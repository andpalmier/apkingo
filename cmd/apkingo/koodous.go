package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

// KoodousInfo - struct for storing resuls from Koodous
type KoodousInfo struct {
	Url            string        `json:"url"`
	SubmissionDate string        `json:"created_at"`
	Id             string        `json:"id"`
	IconLink       string        `json:"image"`
	Size           int64         `json:"size"`
	Tags           []interface{} `json:"tags"`
	Rating         int64         `json:"rating"`
	Detected       bool          `json:"is_detected"`
	Corrupted      bool          `json:"is_corrupted"`
}

// koodousInfo(koodousAPIkey, hash) - get apk info from Koodous using Koodous API key and sha256 hash
func koodousInfo(kapi string, hash string) error {

	koodousresult := KoodousInfo{}
	_, body, err := gorequest.New().
		Get("https://developer.koodous.com/apks/"+hash).
		Set("Authorization", "Token "+kapi).
		End()

	if err != nil {
		return errors.New("error reaching Koodous")
	}
	koodousbody := []byte(body)

	if strings.Contains(string(koodousbody), "detail") {
		if strings.Contains(string(koodousbody), "Not found") {
			italic.Println("app not found in Koodous")
			return nil
		}
		return errors.New(strings.Split(body, "\"")[3])
	}

	jsonerr := json.Unmarshal(koodousbody, &koodousresult)
	if jsonerr != nil {
		return errors.New("error parsing Koodous result")
	}

	fmt.Printf("URL:\t\t")
	printer(koodousresult.Url)

	t, terr := time.Parse(time.RFC3339Nano, koodousresult.SubmissionDate)
	submissiondate := ""
	if terr == nil {
		submissiondate = t.Format(time.RFC822)
	}

	fmt.Printf("Submiss. date:\t")
	printer(submissiondate)
	fmt.Printf("Koodous ID:\t")
	printer(koodousresult.Id)
	fmt.Printf("Icon URL:\t")
	printer(koodousresult.IconLink)
	fmt.Printf("Size:\t\t")
	printer(fmt.Sprintf("%d", koodousresult.Size))
	fmt.Printf("Tags:\t\t")
	if len(koodousresult.Tags) == 0 {
		italic.Printf("no tags found\n")
	} else {
		cyan.Printf("%v\n", koodousresult.Tags)
	}
	fmt.Printf("Rating:\t\t")
	if koodousresult.Rating < 0 {
		red.Printf("%d\n", koodousresult.Rating)
	} else {
		cyan.Printf("%d\n", koodousresult.Rating)
	}
	fmt.Printf("Detected:\t")
	if koodousresult.Detected {
		red.Printf("[!] apk detected as malicious\n")
	} else {
		cyan.Printf("%t\n", koodousresult.Detected)
	}
	fmt.Printf("Corrupted:\t")
	if koodousresult.Corrupted {
		red.Printf("%t\n", koodousresult.Corrupted)
	} else {
		cyan.Printf("%t\n", koodousresult.Corrupted)
	}
	return nil
}
