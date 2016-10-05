package main

import (
	"io/ioutil"

	"github.com/Goryudyuma/anaconda"
	"github.com/Goryudyuma/tlc/tlc"
	"github.com/davecgh/go-spew/spew"
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

func loadyaml() tlc.Config {
	key := tlc.Config{}

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

func main() {

	key := loadyaml()
	//	spew.Dump(key)
	//	fmt.Println(key.AccessToken)

	/*
		operator := byte('+')
		list1 := tlc.List{Listname: "aaa", Owner_screen_name: "Goryudyuma", Owner_id: 0}
		list2 := tlc.List{Listname: "bbb", Owner_screen_name: "Goryudyuma", Owner_id: 0}
		resultlist := tlc.List{Listname: "ccc", Owner_screen_name: "umaumakey", Owner_id: 0}
	*/

	anaconda.SetConsumerKey(key.ConsumerKey)
	anaconda.SetConsumerSecret(key.ConsumerSecret)

	url, test, err := anaconda.AuthorizationURL("http://localhost:8080/callback")
	spew.Dump(url)
	spew.Dump(test)
	spew.Dump(err)
	//spew.Dump(anaconda.GetCredentials(test, test.Secret))

	r := gin.Default()

	store := sessions.NewCookieStore([]byte(key.SeedString))
	//store.Options(sessions.Options{Secure: true})
	r.Use(sessions.Sessions("tlcsession", store))

	r.GET("/", func(c *gin.Context) {
		if !checklogin(c) {
			c.Redirect(301, "/login")
		}
		c.String(200, "logined")
	})
	r.GET("/login", func(c *gin.Context) {
		if checklogin(c) {
			c.Redirect(301, "/")
		}
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
		c.Redirect(301, "/")
	})
	api := r.Group("/api")
	{
		api.POST("/query", func(c *gin.Context) {
			if !checklogin(c) {
				c.String(403, "Not login")
			}
			c.String(200, "queryaaa")
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
