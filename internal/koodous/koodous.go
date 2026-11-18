// Package koodous provides integration with the Koodous Android malware analysis API.
package koodous

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/andpalmier/apkingo/v3/internal/constants"
	"github.com/andpalmier/apkingo/v3/internal/utils"
)

const maxKoodousResponseBytes = 4 << 20

var (
	koodousBaseURL = "https://developer.koodous.com"
	httpClient     = &http.Client{Timeout: constants.DefaultHTTPTimeout}
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

// GetInfo retrieves Koodous information for a given hash.
// It makes an authenticated API request to Koodous with a timeout.
func GetInfo(kapi, hash string) (*KoodousInfo, error) {
	url := fmt.Sprintf("%s/apks/%s", koodousBaseURL, hash)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create Koodous request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+kapi)

	resp, err := httpClient.Do(req)
	if err != nil {
		utils.LogError("error reaching Koodous", err)
		return nil, fmt.Errorf("error reaching Koodous: %w", err)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			if err := resp.Body.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to close response body: %v\n", err)
			}
		}
	}()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxKoodousResponseBytes+1))
	if err != nil {
		utils.LogError("error reading Koodous response", err)
		return nil, fmt.Errorf("error reading Koodous response: %w", err)
	}
	if len(body) > maxKoodousResponseBytes {
		return nil, fmt.Errorf("koodous response exceeds %d bytes", maxKoodousResponseBytes)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var apiErr struct {
			Detail string `json:"detail"`
		}
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Detail != "" {
			err := errors.New(apiErr.Detail)
			utils.LogError("error interpreting Koodous response", err)
			return nil, err
		}
		return nil, fmt.Errorf("koodous returned %s", resp.Status)
	}

	var koodousResult KoodousInfo
	if err := json.Unmarshal(body, &koodousResult); err != nil {
		utils.LogError("error parsing Koodous result", err)
		return nil, fmt.Errorf("error parsing Koodous result: %w", err)
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
