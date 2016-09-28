// tlc project tlc.go
package tlc

import (
	"net/url"

	"github.com/Goryudyuma/anaconda"
	"github.com/davecgh/go-spew/spew"
	"github.com/deckarep/golang-set"
)

func changeint64interface(vs []int64) []interface{} {
	ret := make([]interface{}, len(vs))
	for i, v := range vs {
		ret[i] = v
	}
	return ret
}

func calc(operator byte, one, another []int64) []int64 {
	setone := mapset.NewSetFromSlice(changeint64interface(one))
	setanother := mapset.NewSetFromSlice(changeint64interface(another))

	result := mapset.NewSet()

	switch operator {
	case '+':
		result = setone.Union(setanother)
	case '*':
		result = setone.Intersect(setanother)
	case '-':
		result = setone.Difference(setanother)
	default:
		panic("can't find operator")
	}

	ret := make([]int64, 0, result.Cardinality())

	for _, v := range result.ToSlice() {
		ret = append(ret, v.(int64))
	}
	return ret
}

func choiceuseridfromlist(api anaconda.TwitterApi, listname string, owner_screen_name string, owner_id int64, v url.Values) []int64 {
	var ret []int64

	users, err := api.GetListMembersBySlug(listname, owner_screen_name, owner_id, nil)
	if err != nil {
		panic(err)
	}

	for _, user := range users.Users {
		ret = append(ret, user.Id)
	}
	return ret
}

func Tlc(key MyTwitterKey) {

	anaconda.SetConsumerKey(key.ConsumerKey)
	anaconda.SetConsumerSecret(key.ConsumerSecret)
	api := anaconda.NewTwitterApi(key.AccessToken, key.AccessTokenSecret)

	users2016 := choiceuseridfromlist(*api, "2016", "Goryudyuma", 0, nil)
	userstlctest1 := choiceuseridfromlist(*api, "tlctest1", "Goryudyuma", 0, nil)
	spew.Dump(users2016)
	spew.Dump(userstlctest1)

	spew.Dump(calc('+', users2016, userstlctest1))
	spew.Dump(calc('*', users2016, userstlctest1))
	spew.Dump(calc('-', users2016, userstlctest1))

}

/*
	lists, err := api.GetListsOwnedBy(119667108, nil)
	if err != nil {
		panic(err)
	}
*/
//spew.Dump(lists)
/*
	for _, list := range lists {
		spew.Dump(list)
	}
	fmt.Print("!")
*/
/*
	lists, err := api.CreateList("TLCtest", "TLC", nil)
	spew.Dump(lists)

	if err != nil {
		panic(err)
	}

	users, err := api.AddUserToList("Goryudyuma", lists.Id, nil)

	spew.Dump(users)

	users, err = api.AddUserToList("Goryudyuma2", lists.Id, nil)

	spew.Dump(users)
*/
/*
	users, err := api.GetListMembersBySlug("2016", "Goryudyuma", 119667108, nil)
	if err != nil {
		panic(err)
	}
*/
