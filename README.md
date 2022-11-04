# apkingo

<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/img/apkingo.png?raw=true" />
  <p align="center">
    <a href="https://github.com/andpalmier/apkingo/blob/main/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/andpalmier/apkingo"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/andpalmier/apkingo?style=flat-square"></a>
    <a href="https://twitter.com/intent/follow?screen_name=andpalmier"><img src="https://img.shields.io/twitter/follow/andpalmier?style=social&logo=twitter" alt="follow on Twitter"></a>
  </p>
</p>

apkingo is a tool written in Go to get detailed information about an apk file. apkingo will explore the given file to get details on the apk, such as package name, target SDK, permissions, metadata, certificate serial and issuer. The tool will also retrieve information about the specified apk from the Play Store, and (if valid API keys are provided) from [Koodous](https://koodous.com/) and [VirusTotal](https://virustotal.com).

## Installation

You can can download apkingo from the [Releases section](https://github.com/andpalmier/apkingo/releases) or compile it from the source by downloading the repository, navigating into the `apkingo` directory and building the project with `make apkingo`. This will create a  `build` folder, containing the resulting executable.

## Usage

You can run apkingo with the following flags:

- `-apk` to specify the path to the apk file (**required**)
- `-json` to specify the path of the json file where the results will be exported
- `-vt` to specify VirusTotal API key (required for VirusTotal analysis)
- `-k` to specify Koodous API key (required for Koodous analysis)

Example:
```apkingo -apk snapseed.apk -json snapseed_analysis.json```

## Screenshots

apkingo analyzing snapseed:
<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/img/screen_snapseed.png?raw=true" />
</p>

apkingo analyzing F-Droid:
<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/img/screen_f-droid.png?raw=true" />
</p>

apkingo analyzing an Android malware (I had to cut the screenshot on the permissions section):
<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/img/screen_malware.png?raw=true" />
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/img/screen_malware2.png?raw=true" />
</p>

## Features

Here is the full list of information which apkingo can retrieve:

- General information: app name, package name, app version, MainActivity, minimum and target SDK
- Hashes: md5, sha1 and sha256
- Permissions
- Metadata
- Certificate information: serial, sha1, subject, issuer, validity date and expiration date
- Play Store information: Play Store url, version, summary, release date, number of installations, score, developer name, developer ID, developer mail and developer website
- Koodous info (API key required): Koodous url, Koodous ID, app name, package name, company, version, Koodous link to the app icon, size, Koodous tags, trusted (boolean), installed on devices (boolean), Koodous rating, detected (boolean), corrupted (boolean), statically analyzed (boolean), dynamically analyzed (boolean) and date when the app was submitted to Koodous for the first time
- VirusTotal info (API key required): VirusTotal url, apk names, first submission date, number of submissions, last analysis date and results, reputation, community votes (harmless and malicious), md5 and dhash of icon, providers, receivers, services, interesting string and permissions that are considered dangerous.

## 3rd party libraries used

- [shogo81148/androidbinary](https://github.com/shogo82148/androidbinary)
- [avast/apkverifier](https://github.com/avast/apkverifier)
- [fatih/color](https://github.com/fatih/color)
- [n0madic/google-play-scraper](https://github.com/n0madic/google-play-scraper)
- [parnurzeal/gorequest](https://github.com/parnurzeal/gorequest)
- [VirusTotal/vt-go](https://github.com/VirusTotal/vt-go)
