package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/Goryudyuma/anaconda"
	"github.com/Goryudyuma/tlc/tlc"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func loadconfig() []byte {
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	return data
}

func loadyaml() Config {
	key := Config{}

	err := yaml.Unmarshal(loadconfig(), &key)
	if err != nil {
		panic(err)
	}
	return key
}

func checklogin(c *gin.Context) bool {
	session := sessions.Default(c)
	OathToken := session.Get("OathToken")
	OauthTokenSecret := session.Get("OauthTokenSecret")
	if OathToken == nil || OauthTokenSecret == nil {
		return false
	}
	return true
}

type formula struct {
	Operator   string   `json:"Operator"`
	List1      tlc.List `json:"List1"`
	List2      tlc.List `json:"List2"`
	ResultList tlc.List `json:"ResultList"`
}

var mutexControl sync.Mutex

func run(c *gin.Context, d formula) error {
	session := sessions.Default(c)
	OathToken := session.Get("OathToken").(string)
	OauthTokenSecret := session.Get("OauthTokenSecret").(string)

	operator := byte(d.Operator[0])
	list1 := d.List1
	list2 := d.List2
	resultlist := d.ResultList

	mutexControl.Lock()
	defer mutexControl.Unlock()

	api := anaconda.NewTwitterApi(OathToken, OauthTokenSecret)
	defer api.Close()
	err := tlc.Tlc(*api, operator, list1, list2, resultlist)

	return err
}

func searchlists(c *gin.Context, username string) ([]string, error) {
	session := sessions.Default(c)
	OathToken := session.Get("OathToken").(string)
	OauthTokenSecret := session.Get("OauthTokenSecret").(string)

	mutexControl.Lock()

	api := anaconda.NewTwitterApi(OathToken, OauthTokenSecret)

	user, err := api.GetUsersShow(username, nil)
	if err != nil {
		return nil, err
	}
	lists, err := api.GetListsOwnedBy(user.Id, nil)
	if err != nil {
		return nil, err
	}

	api.Close()
	mutexControl.Unlock()

	ret := make([]string, len(lists))
	for i, list := range lists {
		ret[i] = list.Name
	}
	return ret, nil
}

func listusers(c *gin.Context, list tlc.List) (anaconda.UserCursor, error) {
	session := sessions.Default(c)
	OathToken := session.Get("OathToken").(string)
	OauthTokenSecret := session.Get("OauthTokenSecret").(string)

	mutexControl.Lock()
	defer mutexControl.Unlock()
	api := anaconda.NewTwitterApi(OathToken, OauthTokenSecret)
	defer api.Close()

	users, err := api.GetListMembersBySlug(list.Listname, list.Owner_screen_name, list.Owner_id, nil)
	return users, err

}

func main() {

	key := loadyaml()

	anaconda.SetConsumerKey(key.ConsumerKey)
	anaconda.SetConsumerSecret(key.ConsumerSecret)

	url, test, err := anaconda.AuthorizationURL(key.CallbackURL)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.LoadHTMLGlob("content/*")

	store := sessions.NewCookieStore([]byte(key.SeedString))
	//store.Options(sessions.Options{Secure: true})
	r.Use(sessions.Sessions("tlcsession", store))

	r.GET("/", func(c *gin.Context) {
		if !checklogin(c) {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Redirect(301, "/login")
		}
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/login", func(c *gin.Context) {
		if checklogin(c) {
			c.Redirect(301, "/")
		}
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Redirect(301, url)
	})
	r.GET("/logout", func(c *gin.Context) {

		session := sessions.Default(c)
		session.Clear()
		session.Save()

		c.String(200, "logout")
	})
	r.GET("/callback/", func(c *gin.Context) {
		session := sessions.Default(c)
		//		a := c.Query("oauth_token")
		b := c.Query("oauth_verifier")

		_, user, _ := anaconda.GetCredentials(test, b)

		session.Set("OathToken", user.Get("oauth_token"))
		session.Set("OauthTokenSecret", user.Get("oauth_token_secret"))
		session.Save()

		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Redirect(301, "/")
	})
	api := r.Group("/api")
	{
		api.POST("/query", func(c *gin.Context) {
			if !checklogin(c) {
				c.String(403, "Not login")
			} else {
				query := c.PostForm("query")
				//spew.Dump(query)
				var d formula
				err := json.Unmarshal([]byte(query), &d)
				if err != nil {
					c.String(500, err.Error())
				} else {
					//spew.Dump(d)
					err := run(c, d)
					if err != nil {
						c.String(500, err.Error())
					} else {
						c.String(200, "ok")
					}
				}
			}
		})
		api.POST("/list", func(c *gin.Context) {
			if !checklogin(c) {
				c.String(403, "Not login")
			} else {
				username := c.PostForm("username")
				if username == "" {
					c.String(500, "username is empty")
				} else {
					lists, err := searchlists(c, username)
					if err != nil {
						c.String(500, err.Error())
					} else {
						data, err := json.Marshal(lists)
						if err != nil {
							c.String(500, err.Error())
						} else {
							c.String(200, string(data))
						}
					}
				}
			}
		})
		api.POST("/listusers", func(c *gin.Context) {
			if !checklogin(c) {
				c.String(403, "Not login")
			} else {
				username := c.PostForm("username")
				listname := c.PostForm("listname")
				if username == "" {
					c.String(500, "username is empty")
				} else if listname == "" {
					c.String(500, "listname is empty")
				} else {
					users, err := listusers(c, tlc.List{Listname: listname, Owner_screen_name: username})
					if err != nil {
						c.String(500, err.Error())
					} else {
						c.JSON(200, users)
					}
				}
			}

		})
	}
	r.Run()

}

/*

	api := anaconda.NewTwitterApi(user.Get("oauth_token"), user.Get("oauth_token_secret"))

	apis := make(map[string]anaconda.TwitterApi)
	apis["umaumakey"] = *api

	err = tlc.Tlc(apis, operator, list1, list2, resultlist)
	spew.Dump(err)

	c.String(200, a+b)
*/

//{"Operator":"*","List1":{"Listname":"aaa","OwnerScreenName":"Goryudyuma","OwnerId":0},"List2":{"Listname":"bbb","OwnerScreenName":"Goryudyuma","OwnerId":0},"ResultList":{"Listname":"ccc","OwnerScreenName":"Goryudyuma","OwnerId":0}}
