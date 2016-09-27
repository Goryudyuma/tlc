// tlc project tlc.go
package tlc

import (
	"fmt"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

func Tlc(key MyTwitterKey) {

	anaconda.SetConsumerKey(key.ConsumerKey)
	anaconda.SetConsumerSecret(key.ConsumerSecret)
	api := anaconda.NewTwitterApi(key.AccessToken, key.ConsumerSecret)

	v := url.Values{}
	v.Set("count", "30")
	result, _ := api.GetSearch("golang", v)
	fmt.Println(result)
}
