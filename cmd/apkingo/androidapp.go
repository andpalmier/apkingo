package main

import (
	"encoding/json"
	"os"

	"github.com/shogo82148/androidbinary/apk"
)

// AndroidApp represents information extracted from an APK file
type AndroidApp struct {
	Name         string          `json:"name"`
	PackageName  string          `json:"package-name"`
	Version      string          `json:"version"`
	MainActivity string          `json:"main-activity"`
	MinimumSDK   int32           `json:"minimum-sdk"`
	TargetSDK    int32           `json:"target-sdk"`
	Hashes       Hashes          `json:"hashes"`
	Permissions  []string        `json:"permissions"`
	Metadata     []Metadata      `json:"metadata"`
	Certificate  CertificateInfo `json:"certificate"`
	PlayStore    *PlayStoreInfo  `json:"playstore,omitempty"`
	Koodous      *KoodousInfo    `json:"koodous,omitempty"`
	VirusTotal   *VirusTotalInfo `json:"virustotal,omitempty"`
}

// Hashes represents hash values
type Hashes struct {
	Md5    string `json:"md5"`
	Sha1   string `json:"sha1"`
	Sha256 string `json:"sha256"`
}

// Metadata represents metadata
type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

// ExportJSON exports AndroidApp struct to a JSON file
func (app *AndroidApp) ExportJSON(jsonpath string) error {
	jsonfile, err := json.MarshalIndent(app, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(jsonpath, jsonfile, 0644)
}

// setGeneralInfo sets general information about the APK
func (app *AndroidApp) setGeneralInfo(apk *apk.Apk) error {

	var err error

	app.Name, err = apk.Label(nil)
	if err != nil {
		return err
	}

	app.PackageName = apk.PackageName()

	app.Version, err = apk.Manifest().VersionName.String()
	logError("error getting version information", err)

	app.MainActivity, err = apk.MainActivity()
	logError("error getting main activity information", err)

	app.MinimumSDK, err = apk.Manifest().SDK.Min.Int32()
	logError("error getting minimum SDK information", err)

	app.TargetSDK, err = apk.Manifest().SDK.Target.Int32()
	logError("error getting target SDK information", err)

	for _, n := range apk.Manifest().UsesPermissions {
		permission, _ := n.Name.String()
		if permission != "" {
			app.Permissions = append(app.Permissions, permission)
		}
	}

	for _, n := range apk.Manifest().App.MetaData {
		metadataName, _ := n.Name.String()
		metadataValue, _ := n.Value.String()
		if metadataName != "" {
			app.Metadata = append(app.Metadata, Metadata{
				Name:  metadataName,
				Value: metadataValue,
			})
		}
	}

	return nil
}
