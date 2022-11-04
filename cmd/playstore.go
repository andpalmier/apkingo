package main

import (
	"html"

	"github.com/n0madic/google-play-scraper/pkg/app"
)

// SetPlayStoreInfo(app) - store Play Store info in the androidapp struct
func (androidapp *AndroidApp) setPlayStoreInfo(playstore app.App) {
	playstoreitem := PlayStoreInfo{}
	playstoreitem.Url = playstore.URL
	playstoreitem.Version = playstore.Version
	playstoreitem.Summary = html.UnescapeString(playstore.Summary)
	playstoreitem.Developer = playstore.Developer
	playstoreitem.DeveloperId = playstore.DeveloperID
	playstoreitem.DeveloperMail = playstore.DeveloperEmail
	playstoreitem.DeveloperURL = playstore.DeveloperWebsite
	playstoreitem.Release = playstore.Released
	playstoreitem.Installs = playstore.Installs
	playstoreitem.Score = playstore.Score
	androidapp.PlayStore = &playstoreitem
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
