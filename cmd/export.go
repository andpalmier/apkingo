package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ExportJson (jsonpath) - export androidapp struct to json file
func (androidapp *AndroidApp) ExportJson(jsonpath string) error {
	jsonfile, err := json.MarshalIndent(androidapp, "", "\t")
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = ioutil.WriteFile(jsonpath, jsonfile, 0644)
	if err != nil {
		fmt.Printf("%s", err)
	}
	return nil
}
