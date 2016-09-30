// tlc project tlc.go
package tlc

import (
	"fmt"
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

func choiceuseridfromlist(api anaconda.TwitterApi, listname string, owner_screen_name string, owner_id int64, v url.Values) ([]int64, error) {
	var ret []int64

	users, err := api.GetListMembersBySlug(listname, owner_screen_name, owner_id, nil)
	if err != nil {
		return nil, err
	}

	for _, user := range users.Users {
		ret = append(ret, user.Id)
	}
	return ret, nil
}

func calclist(api anaconda.TwitterApi, operator byte,
	listname1 string, owner_screen_name1 string, owner_id1 int64,
	listname2 string, owner_screen_name2 string, owner_id2 int64,
	v url.Values) ([]int64, error) {
	one, err := choiceuseridfromlist(api, listname1, owner_screen_name1, owner_id1, v)
	if err != nil {
		return nil, err
	}
	another, err := choiceuseridfromlist(api, listname2, owner_screen_name2, owner_id2, v)
	if err != nil {
		return nil, err
	}
	ret := calc(operator, one, another)
	return ret, nil
}

func mergelist(api anaconda.TwitterApi, operator byte,
	listname1 string, owner_screen_name1 string, owner_id1 int64,
	listname2 string, owner_screen_name2 string, owner_id2 int64,
	resultlistname string,
	v url.Values) error {
	resultlist, err := calclist(api, operator,
		listname1, owner_screen_name1, owner_id1,
		listname2, owner_screen_name2, owner_id2,
		v)
	owner, _ := api.GetSelf(v)
	prevresultlist, err := choiceuseridfromlist(api, resultlistname, "", owner.Id, v)
	if err != nil {
		_, err = api.CreateList(resultlistname, listname1+string(operator)+listname2, v)
		if err != nil {
			return err
		}
	}
	removelist := calc('-', prevresultlist, resultlist)
	_, err = api.RemoveUserToListIds(removelist, resultlistname, owner.Id, v)
	if err != nil {
		return err
	}
	addlist := calc('-', resultlist, prevresultlist)
	_, err = api.AddUserToListIds(addlist, resultlistname, owner.Id, v)
	if err != nil {
		return err
	}
	return err
}

func Tlc(key MyTwitterKey) {

	anaconda.SetConsumerKey(key.ConsumerKey)
	anaconda.SetConsumerSecret(key.ConsumerSecret)
	api := anaconda.NewTwitterApi(key.AccessToken, key.AccessTokenSecret)

	aaa, err := choiceuseridfromlist(*api, "aaa", "Goryudyuma", 0, nil)
	if err != nil {
		fmt.Print("Not found list")
		return
	}
	bbb, err := choiceuseridfromlist(*api, "bbb", "Goryudyuma", 0, nil)
	if err != nil {
		fmt.Print("Not found list")
		return
	}
	spew.Dump(aaa)
	spew.Dump(bbb)

	spew.Dump(calc('+', aaa, bbb))
	spew.Dump(calc('*', aaa, bbb))
	spew.Dump(calc('-', aaa, bbb))
	spew.Dump(calc('-', bbb, aaa))

	err = mergelist(*api, '+', "aaa", "Goryudyuma", 0, "bbb", "Goryudyuma", 0, "ccc", nil)
	if err != nil {
		panic(err)
	}
}

/*
		lists, _ := api.CreateList("test", "aaa", nil)
		os.Sleep(10000)
		fmt.Print("!")
		spew.Dump(lists)
		_, err := api.AddUserToList("Goryudyuma", lists.Id, nil)
		if err != nil {
			panic(err)
		}

	_, err := api.RemoveMemberFromList("test4", "Goryudyuma", 0, "Goryudyuma", 0, nil)
	if err != nil {
		panic(err)
	}
	return
*/
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
