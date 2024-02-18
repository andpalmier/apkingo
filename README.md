# apkingo

<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/img/apkingo.png?raw=true" width="400" />
  <p align="center">
    <a href="https://github.com/andpalmier/apkingo/blob/main/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square"></a>
    <a href="https://godoc.org/github.com/andpalmier/apkingo"><img alt="GoDoc Card" src="https://godoc.org/github.com/andpalmier/apkingo?status.svg"></a>
    <a href="https://goreportcard.com/report/github.com/andpalmier/apkingo"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/andpalmier/apkingo?style=flat-square"></a>
    <a href="https://x.com/intent/follow?screen_name=andpalmier"><img src="https://img.shields.io/twitter/follow/andpalmier?style=social&logo=x" alt="follow on X"></a>
  </p>
</p>

apkingo is a utility designed to extract information from an APK file. By analyzing the provided file, apkingo extracts details like package name, target SDK, permissions, metadata, certificate serial, and issuer. Additionally, the tool fetches information about the specified apk from the Play Store and, if valid API keys are provided, from Koodous and VirusTotal. If the file is not available on VirusTotal, apkingo offers the option to upload it.

## Installation

You can download apkingo from the [releases section](https://github.com/andpalmier/apkingo/releases) or compile it from the source by running:

```
go install github.com/andpalmier/apkingo/cmd/apkingo@latest
```

## Usage

You can run apkingo with the following flags:
```
  -apk string (REQUIRED)
        Path to APK file
  -country string
        Country code of the Play Store (default "us")
  -json string
        Path to export analysis in JSON format
  -kapi string
        Koodous API key (you can export it using the env variable KOODOUS_API_KEY)
  -vtapi string
        VirusTotal API key (you can export it using the env variable VT_API_KEY)
```

## Screenshots

apkingo analyzing snapseed:
<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/img/screen_snapseed.png?raw=true" />
</p>

apkingo analyzing an Android malware (I had to cut the screenshot on the permissions section):
<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/img/screen_malware.png?raw=true" />
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/img/screen_malware2.png?raw=true" />
</p>

## Features

Here is the full list of information which apkingo can retrieve:

- General information: app name, package name, app version, main activity, minimum and target SDK
- Hashes: md5, sha1 and sha256
- Permissions
- Metadata
- Certificate information: serial, thumbprint, validity, date, expiration date, issuer and subject
- Play Store information: Play Store url, version, release date, last update date, genre, summary, number of installations, score, developer name, developer ID, developer mail and developer website
- Koodous info (API key required): Koodous url, Koodous ID, Koodous link to the app icon, size, Koodous tags, trusted (boolean), Koodous rating, corrupted (boolean) and submission date
- VirusTotal info (API key required): VirusTotal url, apk names, submission date, number of submissions, last analysis date and results, community votes (harmless and malicious), md5 and dhash of icon, providers, receivers, services, interesting strings and permissions that are considered dangerous

## 3rd party libraries and API documentation 

- shogo81148/androidbinary: [GitHub repo](https://github.com/shogo82148/androidbinary) and [Go reference](https://pkg.go.dev/github.com/shogo82148/androidbinary)
- avast/apkverifier: [GitHub repo](https://github.com/avast/apkverifier) and [Go reference](https://pkg.go.dev/github.com/avast/apkverifier)
- fatih/color: [GitHub repo](https://github.com/fatih/color) and [Go reference](https://pkg.go.dev/github.com/fatih/color)
- n0madic/google-play-scraper: [GitHub repo](https://github.com/n0madic/google-play-scraper) and [Go reference](https://pkg.go.dev/github.com/n0madic/google-play-scraper)
- parnurzeal/gorequest: [GitHub repo](https://github.com/parnurzeal/gorequest) and [Go reference](https://pkg.go.dev/github.com/parnurzeal/gorequest)
- VirusTotal/vt-go: [GitHub repo](https://github.com/VirusTotal/vt-go) and [Go reference](https://pkg.go.dev/github.com/VirusTotal/vt-go)
- [VirusTotal API documentation](https://docs.virustotal.com/reference/overview)
- [Koodous API documentation](https://docs.koodous.com/api/apks.html)
