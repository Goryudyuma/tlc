package main

import (
	"fmt"
	"io/ioutil"

	"github.com/Goryudyuma/tlc/tlc"
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
	fmt.Print("%v", key)
	return key
}

func main() {
	key := loadyaml()
	tlc.Tlc(key)
}
