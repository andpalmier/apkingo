// Package koodous provides integration with the Koodous Android malware analysis API.
package koodous

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/andpalmier/apkingo/internal/utils"
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
	Trusted        bool     `json:"trusted"`
}

// GetInfo saves Koodous information
func GetInfo(kapi, hash string) (*KoodousInfo, error) {
	url := fmt.Sprintf("https://developer.koodous.com/apks/%s", hash)
	resp, body, errs := gorequest.New().Get(url).Set("Authorization", "Token "+kapi).End()
	if len(errs) > 0 {
		utils.LogError("error reaching Koodous", errs[0])
		return nil, fmt.Errorf("error reaching Koodous: %v", errs[0])
	}
	defer resp.Body.Close()

	if strings.Contains(body, "Not found") {
		return nil, nil
	}

	if strings.Contains(body, "detail") {
		split := strings.Split(body, "\"")
		if len(split) > 4 {
			err := errors.New(strings.Split(body, "\"")[3])
			utils.LogError("error interpreting Koodous response", err)
			return nil, err
		} else {
			err := errors.New("error interpreting Koodous response")
			utils.LogError("error interpreting Koodous response", err)
			return nil, err
		}
	}

	var koodousResult KoodousInfo
	if err := json.Unmarshal([]byte(body), &koodousResult); err != nil {
		utils.LogError("error parsing Koodous result", err)
		return nil, fmt.Errorf("error parsing Koodous result: %s", err)
	}

	return &KoodousInfo{
		Url:            strings.Replace(koodousResult.Url, "developer.", "", 1),
		Id:             koodousResult.Id,
		IconLink:       koodousResult.IconLink,
		Size:           koodousResult.Size,
		Tags:           koodousResult.Tags,
		Detected:       koodousResult.Detected,
		Rating:         koodousResult.Rating,
		Corrupted:      koodousResult.Corrupted,
		Trusted:        koodousResult.Trusted,
		SubmissionDate: parseSubmissionDate(koodousResult.SubmissionDate),
	}, nil
}

func parseSubmissionDate(date string) string {
	t, err := time.Parse(time.RFC3339Nano, date)
	if err != nil {
		utils.LogError("error parsing submission date", err)
		return ""
	}
	return t.Format(time.DateTime)
}
