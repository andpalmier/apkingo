package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

// KoodousInfo - struct for storing Koodous API response
type KoodousInfo struct {
	Url            string   `json:"url"`
	Id             string   `json:"id"`
	SubmissionDate string   `json:"created_at"`
	IconLink       string   `json:"icon-link"`
	Size           int64    `json:"size"`
	Tags           []string `json:"tags"`
	Detected       bool     `json:"is_detected"`
	Rating         int64    `json:"rating"`
	Corrupted      bool     `json:"is_corrupted"`
}

// setKoodousInfo saves Koodous information
func (androidapp *AndroidApp) setKoodousInfo(kapi string) error {
	hash := androidapp.Hashes.Sha256
	url := fmt.Sprintf("https://developer.koodous.com/apks/%s", hash)
	resp, body, errs := gorequest.New().Get(url).Set("Authorization", "Token "+kapi).End()
	if len(errs) > 0 {
		return fmt.Errorf("error reaching Koodous: %v", errs[0])
	}
	defer resp.Body.Close()

	if strings.Contains(body, "Not found") {
		return nil
	}

	if strings.Contains(body, "detail") {
		split := strings.Split(body, "\"")
		if len(split) > 4 {
			var KoodousErr = errors.New(strings.Split(body, "\"")[3])
			return KoodousErr
		} else {
			var KoodousErr = errors.New("error interpreting Koodous response")
			return KoodousErr
		}
	}

	var koodousResult KoodousInfo
	if err := json.Unmarshal([]byte(body), &koodousResult); err != nil {
		return fmt.Errorf("error parsing Koodous result: %s", err)
	}

	androidapp.Koodous = &KoodousInfo{
		Url:            strings.Replace(koodousResult.Url, "developer.", "", 1),
		Id:             koodousResult.Id,
		IconLink:       koodousResult.IconLink,
		Size:           koodousResult.Size,
		Tags:           koodousResult.Tags,
		Detected:       koodousResult.Detected,
		Rating:         koodousResult.Rating,
		Corrupted:      koodousResult.Corrupted,
		SubmissionDate: parseSubmissionDate(koodousResult.SubmissionDate),
	}
	return nil
}

func parseSubmissionDate(date string) string {
	t, err := time.Parse(time.RFC3339Nano, date)
	if err != nil {
		return ""
	}
	return t.Format(time.DateTime)
}
