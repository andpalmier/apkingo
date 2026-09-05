# apkingo

<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/img/apkingo.png?raw=true" width="400" />
  <p align="center">
    <a href="https://github.com/andpalmier/apkingo/blob/main/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/License-Apache%202.0-blue.svg"></a>
    <a href="https://godoc.org/github.com/andpalmier/apkingo"><img alt="GoDoc Card" src="https://godoc.org/github.com/andpalmier/apkingo?status.svg"></a>
    <a href="https://github.com/andpalmier/apkingo/actions/workflows/linter.yaml"><img alt="golangci-lint" src="https://github.com/andpalmier/apkingo/actions/workflows/linter.yaml/badge.svg"></a>
    <a href="https://x.com/intent/follow?screen_name=andpalmier"><img src="https://img.shields.io/twitter/follow/andpalmier?style=social&logo=x" alt="follow on X"></a>
  </p>
</p>

> **apkingo** is an APK analysis tool written in Go. It extracts information from Android applications, such as permissions, metadata, certificate details, version code, supported CPU architectures, and integrates with VirusTotal and Koodous for malware detection.

## Features

From the APK itself, apkingo reads the package name, version name and version
code, main activity, SDK versions and supported CPU architectures. It lists
every requested permission and all application metadata, computes MD5, SHA1 and
SHA256, and reports the signing certificate's serial, thumbprint, validity,
issuer and subject. Pass `-locale` to pull the app name in a specific language,
such as `-locale zh-CN`.

It reads XAPK and APKS archives too, detecting and extracting the APKs inside,
and `-dir` analyzes every APK in a directory in one run.

By default it also scrapes the Google Play Store listing. `-no-play-store`
skips that, which is what you want on an air-gapped machine, behind a
restrictive network, or anywhere Play is unreachable.

With a VirusTotal API key it adds detection statistics with red flags
highlighted, the popular threat classification such as
`trojan.pegasus/chrysaor`, the file reputation score, community detections from
Sigma, YARA and IDS rules, and the file's tags. The same key pulls VirusTotal's
Androguard analysis: the full APK structure, its activities, services,
providers and receivers, its libraries and SDK versions, with dangerous
permissions highlighted. A Koodous API key adds that service's detection
status, community rating and trust score, vote counts, and the linked
repository when one exists.

Results print to the terminal in colour, with malware indicators in bold red.
`-json` writes the whole analysis, VirusTotal and Koodous data included, as
pretty-printed JSON.

## Installation

### From GitHub releases

Download the pre-compiled binary for your system from the [Releases](https://github.com/andpalmier/apkingo/releases) page.

### From source

```bash
go install github.com/andpalmier/apkingo/v2/cmd/apkingo@latest
```

### From Homebrew

```bash
brew install --cask andpalmier/tap/apkingo
```

Homebrew casks are macOS only. On Linux, use `go install` or a pre-built binary.

## Usage

### With Docker

Running the published image needs neither Go nor a downloaded binary. Mount the
directory holding the APK:

```bash
# Analyze an APK
docker run --rm -v $(pwd):/mnt ghcr.io/andpalmier/apkingo -apk /mnt/target.apk

# Analyze an XAPK file
docker run --rm -v $(pwd):/mnt ghcr.io/andpalmier/apkingo -apk /mnt/app.xapk

# Analyze every APK in a directory
docker run --rm -v $(pwd):/mnt ghcr.io/andpalmier/apkingo -dir /mnt

# Analyze and export a JSON report
docker run --rm -v $(pwd):/mnt ghcr.io/andpalmier/apkingo -apk /mnt/target.apk -json /mnt/report.json
```

### From the command line

```bash
apkingo -apk <path_to_apk>            # a single APK
apkingo -apk <path_to_xapk>           # an XAPK or APKS archive
apkingo -dir <path_to_directory>      # every APK in a directory

# With API keys, exporting JSON
apkingo -apk <path_to_apk> -vtapi <VT_KEY> -kapi <KOODOUS_KEY> -json report.json

# The app name in a specific language, with a matching Play Store country
apkingo -apk target.apk -locale zh-CN -country cn
```

### API keys

VirusTotal and Koodous keys can come from the environment, which keeps them out
of your shell history:

```bash
export VT_API_KEY="your_virustotal_api_key"
export KOODOUS_API_KEY="your_koodous_api_key"
apkingo -apk <path_to_apk>
```

They can also be passed as flags:

```bash
apkingo -apk <path_to_apk> -vtapi <YOUR_VT_KEY> -kapi <YOUR_KOODOUS_KEY>
```

### Options

| Flag | Description |
|------|-------------|
| `-apk` | Path to APK or XAPK file to analyze (required) |
| `-dir` | Analyze all APKs in a directory |
| `-json` | Path to export analysis in JSON format |
| `-country` | Country code of the Play Store (default: "us") |
| `-locale` | Locale for localized app name extraction (e.g., `en`, `zh-CN`) |
| `-vtapi` | VirusTotal API key (can also use `VT_API_KEY` env var) |
| `-kapi` | Koodous API key (can also use `KOODOUS_API_KEY` env var) |
| `-no-play-store` | Skip Play Store scraping for offline analysis |
| `-vtupload` | Upload the APK to VirusTotal after analysis (interactive prompt) |
| `-version` | Print the version, commit and build date, then exit |

## Screenshot

apkingo analyzing an Android malware:
<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/img/screen_malware.png?raw=true" />
</p>

## Third party libraries and API documentation

- shogo82148/androidbinary: [GitHub repo](https://github.com/shogo82148/androidbinary) and [Go reference](https://pkg.go.dev/github.com/shogo82148/androidbinary)
- avast/apkverifier: [GitHub repo](https://github.com/avast/apkverifier) and [Go reference](https://pkg.go.dev/github.com/avast/apkverifier)
- fatih/color: [GitHub repo](https://github.com/fatih/color) and [Go reference](https://pkg.go.dev/github.com/fatih/color)
- n0madic/google-play-scraper: [GitHub repo](https://github.com/n0madic/google-play-scraper) and [Go reference](https://pkg.go.dev/github.com/n0madic/google-play-scraper)
- parnurzeal/gorequest: [GitHub repo](https://github.com/parnurzeal/gorequest) and [Go reference](https://pkg.go.dev/github.com/parnurzeal/gorequest)
- VirusTotal/vt-go: [GitHub repo](https://github.com/VirusTotal/vt-go) and [Go reference](https://pkg.go.dev/github.com/VirusTotal/vt-go)
- [VirusTotal API documentation](https://docs.virustotal.com/reference/overview)
- [Koodous API documentation](https://docs.koodous.com/api/apks.html)
