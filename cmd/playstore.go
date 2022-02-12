package main

import (
	"github.com/n0madic/google-play-scraper/pkg/app"
	"html"
)

// SetPlayStoreInfo(app) - store Play Store info in the androidapp struct
func (androidapp *AndroidApp) setPlayStoreInfo(playstore app.App) {
	androidapp.PlayStore.Url = playstore.URL
	androidapp.PlayStore.Version = playstore.Version
	androidapp.PlayStore.Summary = html.UnescapeString(playstore.Summary)
	androidapp.PlayStore.Developer.Name = playstore.Developer
	androidapp.PlayStore.Developer.Mail = playstore.DeveloperEmail
	androidapp.PlayStore.Developer.Id = playstore.DeveloperID
	androidapp.PlayStore.Release = playstore.Released
	androidapp.PlayStore.Installs = playstore.Installs
	androidapp.PlayStore.Score = playstore.Score
}

// searchPlayStore() - search app in Play Store using the package name
func (androidapp *AndroidApp) searchPlayStore() (app.App, error) {
	packagename := androidapp.GeneralInfo.PackageName
	if packagename == "" {
		return app.App{}, err
	}
	playstoreinfo := app.New(packagename, app.Options{
		Country:  "us",
		Language: "us",
	})
	if err = playstoreinfo.LoadDetails(); err != nil {
		return app.App{}, err
	}
	return *playstoreinfo, nil
}
