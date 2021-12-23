package main

import "github.com/n0madic/google-play-scraper/pkg/app"

// searchPlayStore() - search app in Play Store using the package name
func (androidapp *AndroidApp) searchPlayStore() error {
	playstoreinfo := app.New(androidapp.Apk.PackageName(), app.Options{
		Country:  "us",
		Language: "us",
	})

	if err = playstoreinfo.LoadDetails(); err != nil {
		return err
	}

	androidapp.PlayStoreInfo = *playstoreinfo
	androidapp.PlayStoreFound = true
	return nil
}
