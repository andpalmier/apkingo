package main

import (
	"fmt"
	"github.com/shogo82148/androidbinary/apk"
)

// mapping SDK to Android version
var androidname = map[int]string{
	1:  "Android 1",
	2:  "Android 1.1",
	3:  "Android 1.5",
	4:  "Android 1.6",
	5:  "Android 2",
	6:  "Android 2",
	7:  "Android 2.1",
	8:  "Android 2.2",
	9:  "Android 2.3",
	10: "Android 2.3.3",
	11: "Android 3",
	12: "Android 3.1",
	13: "Android 3.2",
	14: "Android 4",
	15: "Android 4.0.3",
	16: "Android 4.1",
	17: "Android 4.2",
	18: "Android 4.3",
	19: "Android 4.4",
	20: "Android 4.4W",
	21: "Android 5",
	22: "Android 5.1",
	23: "Android 6",
	24: "Android 7",
	25: "Android 7.1",
	26: "Android 8",
	27: "Android 8.1",
	28: "Android 9",
	29: "Android 10",
	30: "Android 11",
	31: "Android 12",
	32: "Android 12",
	33: "Android 13",
}

// AndroidApp - struct for saving details about apk
type AndroidApp struct {
	Name        string          `json:"name"`
	GeneralInfo GeneralInfo     `json:"generalinfo"`
	Hashes      Hashes          `json:"hashes"`
	Permissions []string        `json:"permissions"`
	Metadata    []Metadata      `json:"metadata"`
	Certificate CertificateInfo `json:"certificate"`
	PlayStore   PlayStoreInfo   `json:"playstore"`
	Koodous     KoodousInfo     `json:"koodous"`
	// using a pointer here to avoid exporting VT info
	// in case VT api was not specified
	VirusTotal *VirusTotalInfo `json:"virustotal,omitempty"`
}

// GeneralInfo - struct for packagename, apk
// version, main activity and SDK values
type GeneralInfo struct {
	PackageName  string `json:"packagename"`
	Version      string `json:"version"`
	MainActivity string `json:"mainactivity"`
	MinimumSdk   string `json:"minimumsdk"`
	TargetSdk    string `json:"targetsdk"`
}

// Hashes - struct for hash values
type Hashes struct {
	Md5    string `json:"md5"`
	Sha1   string `json:"sha1"`
	Sha256 string `json:"sha256"`
}

// Metadata - struct for metadata
type Metadata struct {
	MetadataName  string `json:"name"`
	MetadataValue string `json:"value,omitempty"`
}

// CertificateInfo - struct for certificate info
type CertificateInfo struct {
	Serial    string `json:"serial"`
	Sha1      string `json:"sha1"`
	Subject   string `json:"subject"`
	Issuer    string `json:"issuer"`
	ValidFrom string `json:"validfrom"`
	ValidTo   string `json:"validto"`
}

// PlayStoreInfo - struct for Play Store info
type PlayStoreInfo struct {
	Url       string    `json:"url"`
	Version   string    `json:"version"`
	Summary   string    `json:"summary"`
	Developer Developer `json:"developer"`
	Release   string    `json:"releasedate"`
	Installs  string    `json:"numberinstalls"`
	Score     float64   `json:"score"`
}

// Developer - struct for info about the developer
type Developer struct {
	Name string `json:"name"`
	Id   string `json:"id"`
	Mail string `json:"mail"`
}

// KoodousInfo - struct for info gathered from Koodous
type KoodousInfo struct {
	Url      string `json:"koodousurl"`
	Analyzed bool   `json:"analyzed"`
	Detected bool   `json:"detected"`
}

// VTAnalysStats - struct for analysis details by VirusTotal
type VTAnalysStats struct {
	Harmless         int64 `json:"harmless"`
	TypeUnsupported  int64 `json:"typeunsupported"`
	Suspicious       int64 `json:"suspicious"`
	ConfirmedTimeout int64 `json:"confirmedtimeout"`
	Timeout          int64 `json:"timeout"`
	Failure          int64 `json:"failure"`
	Malicious        int64 `json:"malicious"`
	Undetected       int64 `json:"undetected"`
}

// VTVotes - struct for vote details by VirusTotal
type VTVotes struct {
	Harmless  int64 `json:"harmless"`
	Malicious int64 `json:"malicious"`
}

// VTIcon - struct for icon details by VirusTotal
type VTIcon struct {
	Md5   string `json:"md5"`
	Dhash string `json:"dhash"`
}

// VTVirusTotalInfo - struct for info gathered from VirusTotal
type VirusTotalInfo struct {
	Url          string        `json:"virustotalurl"`
	Names        []string      `json:"names"`
	FirstSubmit  string        `json:"firstsubmitted"`
	TimesSubmit  int64         `json:"timessubmitted"`
	LastAnalysis string        `json:"lastanalysis"`
	AnalysStats  VTAnalysStats `json:"analysisstats"`
	Reput        int64         `json:"reputation"`
	Votes        VTVotes       `json:"votes"`
	Icon         VTIcon        `json:"icon"`
}

// SetPermissions(apk) - get the permission from apk and store
// them in the androidapp struct
func (androidapp *AndroidApp) setPermissions(apk apk.Apk) {
	for _, n := range apk.Manifest().UsesPermissions {
		permission, _ := n.Name.String()
		if permission != "" {
			androidapp.Permissions = append(androidapp.Permissions, permission)
		}
	}
}

// SetMetadata(apk) - get the metadata from apk and store
// them in the androidapp struct
func (androidapp *AndroidApp) setMetadata(apk apk.Apk) {
	var m Metadata
	for _, n := range apk.Manifest().App.MetaData {
		metaname, _ := n.Name.String()
		metavalue, _ := n.Value.String()
		m.MetadataValue = ""
		if metaname != "" {
			m.MetadataName = metaname
			if metavalue != "" {
				m.MetadataValue = metavalue
			}
			androidapp.Metadata = append(androidapp.Metadata, m)
		}
	}
}

// SetGeneralInfo(apk) - get general info from apk and
// store them in the androidapp struct
func (androidapp *AndroidApp) setApkGeneralInfo(apk apk.Apk) {
	androidapp.Name, err = apk.Label(nil)
	if err != nil {
		androidapp.Name = ""
	}
	androidapp.GeneralInfo.PackageName = apk.PackageName()
	androidapp.GeneralInfo.Version, err = apk.Manifest().VersionName.String()
	if err != nil {
		androidapp.GeneralInfo.Version = ""
	}
	androidapp.GeneralInfo.MainActivity, err = apk.MainActivity()
	if err != nil {
		androidapp.GeneralInfo.MainActivity = ""
	}
	sdktarget, err := apk.Manifest().SDK.Target.Int32()
	if err != nil {
		androidapp.GeneralInfo.TargetSdk = ""
	} else {
		androidapp.GeneralInfo.TargetSdk = fmt.Sprintf("%d (%s)", sdktarget, androidname[int(sdktarget)])
	}
	sdkmin, err := apk.Manifest().SDK.Min.Int32()
	if err != nil {
		androidapp.GeneralInfo.MinimumSdk = ""
	} else {
		androidapp.GeneralInfo.MinimumSdk = fmt.Sprintf("%d (%s)", sdkmin, androidname[int(sdkmin)])
	}
}
