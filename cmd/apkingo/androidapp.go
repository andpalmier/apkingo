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
	err = os.WriteFile(jsonpath, jsonfile, 0644)
	return err
}

// setGeneralInfo sets general information about the APK
func (app *AndroidApp) setGeneralInfo(apk *apk.Apk) {
	name, _ := apk.Label(nil)
	app.Name = name
	app.PackageName = apk.PackageName()
	version, _ := apk.Manifest().VersionName.String()
	app.Version = version
	main, _ := apk.MainActivity()
	app.MainActivity = main
	sdkMin, _ := apk.Manifest().SDK.Min.Int32()
	app.MinimumSDK = sdkMin
	sdkTarget, _ := apk.Manifest().SDK.Target.Int32()
	app.TargetSDK = sdkTarget
	for _, n := range apk.Manifest().UsesPermissions {
		permission, _ := n.Name.String()
		if permission != "" {
			app.Permissions = append(app.Permissions, permission)
		}
	}
	var m Metadata
	for _, n := range apk.Manifest().App.MetaData {
		metadataName, _ := n.Name.String()
		metadataValue, _ := n.Value.String()
		if metadataName != "" {
			m.Name = metadataName
			m.Value = ""
			if metadataValue != "" {
				m.Value = metadataValue
			}
			app.Metadata = append(app.Metadata, m)
		}
	}
}
