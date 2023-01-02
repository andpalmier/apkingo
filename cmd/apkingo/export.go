package main

import (
	"encoding/json"
	"os"
)

// ExportJson (jsonpath) - export androidapp struct to json file
func (androidapp *AndroidApp) ExportJson(jsonpath string) error {
	jsonfile, err := json.MarshalIndent(androidapp, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(jsonpath, jsonfile, 0644)
	if err != nil {
		return err
	}
	return nil
}
