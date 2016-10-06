// tlc project tlc.go
package tlc

import (
	"errors"
	"net/url"

	"github.com/Goryudyuma/anaconda"
	"github.com/davecgh/go-spew/spew"
	//	"github.com/davecgh/go-spew/spew"
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

func choiceuseridfromlist(api anaconda.TwitterApi, list List, v url.Values) ([]int64, error) {
	var ret []int64

	users, err := api.GetListMembersBySlug(list.Listname, list.Owner_screen_name, list.Owner_id, nil)
	if err != nil {
		return nil, err
	}

	for _, user := range users.Users {
		ret = append(ret, user.Id)
	}
	return ret, nil
}

func calclist(api anaconda.TwitterApi, operator byte,
	list1 List, list2 List,
	v url.Values) ([]int64, error) {
	one, err := choiceuseridfromlist(api, list1, v)
	if err != nil {
		return nil, err
	}
	another, err := choiceuseridfromlist(api, list2, v)
	if err != nil {
		return nil, err
	}
	ret := calc(operator, one, another)
	return ret, nil
}

func mergelist(api anaconda.TwitterApi, operator byte,
	list1 List, list2 List,
	resultlistname string,
	v url.Values) error {
	resultlist, err := calclist(api, operator, list1, list2, v)
	if err != nil {
		spew.Dump("aaa")
		return err
	}
	owner, err := api.GetSelf(v)
	if err != nil {
		spew.Dump(err.Error())
		return err
	}
	prevresultlist, err := choiceuseridfromlist(api, List{resultlistname, "", owner.Id}, v)
	if err != nil {
		_, err = api.CreateList(resultlistname, list1.Listname+string(operator)+list2.Listname, v)
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

func checkapiforlist(apis map[string]anaconda.TwitterApi, list List, v url.Values) ([]int64, error) {
	for _, api := range apis {
		result, err := choiceuseridfromlist(api, list, v)
		if err == nil {
			return result, nil
		}
	}
	return nil, errors.New("有効なAPIが見つかりませんでした")
}

func Tlc(api anaconda.TwitterApi, operator byte, list1 List, list2 List, result List) error {
	err := mergelist(api, operator, list1, list2, result.Listname, nil)
	return err
}

/*
func Tlc(apis map[string]anaconda.TwitterApi, operator byte, list1 List, list2 List, result List) error {

		aaa, err := checkapiforlist(apis, list1, nil)
		if err != nil {
			return err
		}
		bbb, err := checkapiforlist(apis, list2, nil)
		if err != nil {
			return err
		}

			spew.Dump(aaa)
			spew.Dump(bbb)

			spew.Dump(calc('+', aaa, bbb))
			spew.Dump(calc('*', aaa, bbb))
			spew.Dump(calc('-', aaa, bbb))
			spew.Dump(calc('-', bbb, aaa))

	err := mergelist(apis[result.Owner_screen_name], operator, list1, list2, result.Listname, nil)
	return err
}

*/
