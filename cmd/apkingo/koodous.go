package main

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

// koodousDetection() - get apk info from Koodous using sha256 hash and store them in the androidapp struct
func (androidapp *AndroidApp) getKoodousDetection(kapi string) error {

	koodousresult := KoodousInfo{}
	_, body, err := gorequest.New().
		Get("https://developer.koodous.com/apks/"+androidapp.Hashes.Sha256).
		Set("Authorization", "Token "+kapi).
		End()

	if err != nil {
		return errors.New("error performing Koodous request, please check your API key")
	}

	koodousRes := []byte(body)
	jsonerr := json.Unmarshal(koodousRes, &koodousresult)

	if jsonerr != nil {
		return errors.New("apk not found in Koodous")
	}

	if strings.HasPrefix(string(koodousRes), "{\"detail") {
		return errors.New("error performing Koodous request, please check your API key")
	}

	androidapp.Koodous = &koodousresult
	t, timerr := time.Parse(time.RFC3339Nano, koodousresult.SubmissionDate)
	if timerr != nil {
		androidapp.Koodous.SubmissionDate = ""
	} else {
		androidapp.Koodous.SubmissionDate = t.Format(time.RFC822)
	}
	return nil
}
