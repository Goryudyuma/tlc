package main

import (
	"io/ioutil"

	"github.com/Goryudyuma/anaconda"
	"github.com/Goryudyuma/tlc/tlc"
	"github.com/davecgh/go-spew/spew"
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

func loadyaml() tlc.MyTwitterKey {
	key := tlc.MyTwitterKey{}

	err := yaml.Unmarshal(loadconfig(), &key)
	if err != nil {
		panic(err)
	}
	return key
}

func main() {

	key := loadyaml()
	//	spew.Dump(key)
	//	fmt.Println(key.AccessToken)

	operator := byte('+')
	list1 := tlc.List{Listname: "aaa", Owner_screen_name: "Goryudyuma", Owner_id: 0}
	list2 := tlc.List{Listname: "bbb", Owner_screen_name: "Goryudyuma", Owner_id: 0}
	resultlist := tlc.List{Listname: "ccc", Owner_screen_name: "umaumakey", Owner_id: 0}

	anaconda.SetConsumerKey(key.ConsumerKey)
	anaconda.SetConsumerSecret(key.ConsumerSecret)

	url, test, err := anaconda.AuthorizationURL("http://localhost:8080/callback")
	spew.Dump(url)
	spew.Dump(test)
	spew.Dump(err)
	//spew.Dump(anaconda.GetCredentials(test, test.Secret))

	r := gin.Default()
	r.GET("/login", func(c *gin.Context) {
		c.Redirect(301, url)
	})
	r.GET("/callback/", func(c *gin.Context) {
		a := c.Query("oauth_token")
		b := c.Query("oauth_verifier")
		_, user, _ := anaconda.GetCredentials(test, b)
		//		spew.Dump(user)
		spew.Dump(user.Get("oauth_token"))
		spew.Dump(user.Get("screen_name"))
		spew.Dump(user.Get("oauth_token_secret"))

		api := anaconda.NewTwitterApi(user.Get("oauth_token"), user.Get("oauth_token_secret"))

		apis := make(map[string]anaconda.TwitterApi)
		apis["umaumakey"] = *api

		err = tlc.Tlc(apis, operator, list1, list2, resultlist)
		spew.Dump(err)

		c.String(200, a+b)
	})
	r.Run()

}
