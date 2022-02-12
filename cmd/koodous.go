package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

// Koodousresult - struct for storing resuls from Koodous
type Koodousresult struct {
	CreatedOn        int           `json:"created_on"`
	Rating           int           `json:"rating"`
	Image            string        `json:"image"`
	Tags             []interface{} `json:"tags"`
	Md5              string        `json:"md5"`
	Sha1             string        `json:"sha1"`
	Sha256           string        `json:"sha256"`
	App              string        `json:"app"`
	PackageName      string        `json:"package_name"`
	Company          string        `json:"company"`
	DisplayedVersion string        `json:"displayed_version"`
	Size             int           `json:"size"`
	Stored           bool          `json:"stored"`
	Analyzed         bool          `json:"analyzed"`
	IsApk            bool          `json:"is_apk"`
	Trusted          bool          `json:"trusted"`
	Detected         bool          `json:"detected"`
	Corrupted        bool          `json:"corrupted"`
	Repo             string        `json:"repo"`
	OnDevices        bool          `json:"on_devices"`
}

// koodousDetection() - get apk info from Koodous using sha256 hash and store them in the androidapp struct
func (androidapp *AndroidApp) getKoodousDetection() error {
	_, body, err := gorequest.New().Get("https://api.koodous.com/apks/" + androidapp.Hashes.Sha256).End()
	if err != nil {
		return errors.New("error performing koodous request")
	}
	koodousRes := []byte(body)
	koodousresult := Koodousresult{}
	jsonerr := json.Unmarshal(koodousRes, &koodousresult)
	if jsonerr != nil {
		return errors.New("apk not found in koodous")
	}
	androidapp.Koodous.Url = fmt.Sprintf("https://koodous.com/apks/%s", androidapp.Hashes.Sha256)
	androidapp.Koodous.Analyzed = koodousresult.Analyzed
	androidapp.Koodous.Detected = koodousresult.Detected
	return nil
}
