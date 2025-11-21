package analyzer

import (
	"errors"
	"html"
	"time"

	"github.com/andpalmier/apkingo/internal/utils"
	playapp "github.com/n0madic/google-play-scraper/pkg/app"
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

// SetPlayStoreInfo searches for an app in Play Store using the package name
func (app *AndroidApp) SetPlayStoreInfo(country string) error {
	packagename := app.PackageName
	if packagename == "" {
		app.PlayStore = nil
		return errors.New("no package name found in the app")
	}

	playStoreInfo := playapp.New(packagename, playapp.Options{
		Country:  country,
		Language: "us",
	})

	if err := playStoreInfo.LoadDetails(); err != nil {
		utils.LogError("Play Store information not found", err)
		return errors.New("play Store information not found")
	}

	app.PlayStore = &PlayStoreInfo{
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
