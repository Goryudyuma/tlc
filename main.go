package main

import (
	"io/ioutil"

	"github.com/Goryudyuma/tlc/tlc"
	"github.com/davecgh/go-spew/spew"
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
	spew.Dump(key)
	//	fmt.Println(key.AccessToken)
	tlc.Tlc(key)
}
