package main

import (
	"errors"
	"fmt"
	"html"

	"github.com/n0madic/google-play-scraper/pkg/app"
)

// searchPlayStore() - search app in Play Store using the package name
func searchPlayStore(packagename string) error {
	if packagename == "" {
		return errors.New("no package name found in the app")
	}
	playstoreinfo := app.New(packagename, app.Options{
		Country:  "us",
		Language: "us",
	})
	if err = playstoreinfo.LoadDetails(); err != nil {
		italic.Println("package name not found in Play Store")
	} else {
		fmt.Printf("URL:\t\t")
		printer(playstoreinfo.URL)
		fmt.Printf("Version:\t")
		printer(playstoreinfo.Version)
		fmt.Printf("Summary:\t")
		printer(html.UnescapeString((playstoreinfo.Summary)))
		fmt.Printf("Release date:\t")
		printer(playstoreinfo.Released)
		fmt.Printf("# installs:\t")
		printer(playstoreinfo.Installs)
		fmt.Printf("Score:\t\t")
		printer(playstoreinfo.ScoreText)
		fmt.Printf("Developer:\t")
		printer(playstoreinfo.Developer)
		fmt.Printf("Developer ID:\t")
		printer(playstoreinfo.DeveloperID)
		fmt.Printf("Developer mail:\t")
		printer(playstoreinfo.DeveloperEmail)
		fmt.Printf("Developer URL:\t")
		printer(playstoreinfo.DeveloperWebsite)
	}
	return nil
}
