package main

import (
	"errors"
	"github.com/n0madic/google-play-scraper/pkg/app"
	"html"
	"time"
)

// PlayStoreInfo - struct for Play Store information
type PlayStoreInfo struct {
	Url       string    `json:"url"`
	Version   string    `json:"version"`
	Release   string    `json:"release-date"`
	Updated   time.Time `json:"updated"`
	Genre     string    `json:"genre"`
	Summary   string    `json:"summary"`
	Installs  string    `json:"number-installs"`
	Score     float64   `json:"score"`
	Developer Developer `json:"developer"`
}

// Developer - struct for information about the developer
type Developer struct {
	Name string `json:"name"`
	Id   string `json:"id"`
	Mail string `json:"mail"`
	URL  string `json:"url"`
}

// setPlayStoreInfo searches for an app in Play Store using the package name
func (androidapp *AndroidApp) setPlayStoreInfo(country string) error {
	packagename := androidapp.PackageName
	if packagename == "" {
		androidapp.PlayStore = nil
		return errors.New("no package name found in the app")
	}

	playStoreInfo := app.New(packagename, app.Options{
		Country:  country,
		Language: "us",
	})

	if err := playStoreInfo.LoadDetails(); err != nil {
		return errors.New("error loading Play Store information")
	}

	androidapp.PlayStore = &PlayStoreInfo{
		Url:       playStoreInfo.URL,
		Version:   playStoreInfo.Version,
		Release:   playStoreInfo.Released,
		Updated:   playStoreInfo.Updated,
		Genre:     playStoreInfo.Genre,
		Summary:   html.UnescapeString(playStoreInfo.Summary),
		Installs:  playStoreInfo.Installs,
		Score:     playStoreInfo.Score,
		Developer: Developer{Name: playStoreInfo.Developer, Id: playStoreInfo.DeveloperID, Mail: playStoreInfo.DeveloperEmail, URL: playStoreInfo.DeveloperWebsite},
	}

	return nil
}
