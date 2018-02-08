package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type testSuite struct {
	Tests    int `xml:"tests,attr"`
	Failures int `xml:"failures,attr"`
}

func main() {

	xmlFile, err := ioutil.ReadFile("./plugins/e2e/results/junit_01.xml")
	if err != nil {
		panic(err)
	}

	t := testSuite{}

	err = xml.Unmarshal([]byte(xmlFile), &t)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	fmt.Printf("%d %d\n", t.Tests, t.Failures)

}
