package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func main() {

	dat, err := ioutil.ReadFile("people-test.yml")
	if err != nil {
		panic(err)
	}

	givers := Group{}
	err = yaml.Unmarshal(dat, &givers)
	if err != nil {
		panic(err)
	}

	recievers := Group{}
	err = yaml.Unmarshal(dat, &recievers)
	if err != nil {
		panic(err)
	}

	matches, err := CreateMatches(givers, recievers)

	fmt.Printf("%+v", matches)

}
