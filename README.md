# apkingo

<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/img/apkingo.png?raw=true" width="400" />
  <p align="center">
    <a href="https://github.com/andpalmier/apkingo/blob/main/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/License-Apache%202.0-blue.svg"></a>
    <a href="https://godoc.org/github.com/andpalmier/apkingo"><img alt="GoDoc Card" src="https://godoc.org/github.com/andpalmier/apkingo?status.svg"></a>
    <a href="https://goreportcard.com/report/github.com/andpalmier/apkingo"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/andpalmier/apkingo?style=flat-square"></a>
    <a href="https://x.com/intent/follow?screen_name=andpalmier"><img src="https://img.shields.io/twitter/follow/andpalmier?style=social&logo=x" alt="follow on X"></a>
  </p>
</p>

> **apkingo** is an APK analysis tool written in Go. It extracts information from Android applications, such as permissions, metadata, certificate details, and integrates with VirusTotal and Koodous for malware detection.

## Features

### Core Analysis
- **General Info**: Package name, version, main activity, SDK versions
- **Hashes**: MD5, SHA1, SHA256
- **Permissions**: Complete list of requested permissions
- **Metadata**: Application metadata
- **Certificate**: Serial, thumbprint, validity, issuer, subject

### File Format Support
- **XAPK/APKS Support**: Automatically detects and extracts APKs from Android App Bundle archives
- **Directory Analysis**: Batch analyze multiple APKs in a directory with the `-dir` flag

### Offline Analysis
- **No Play Store Required**: Use the `-no-play-store` flag to analyze APKs without internet access (useful in restricted network environments, air-gapped systems, or regions with limited connectivity)

### External Intelligence
- **Play Store Integration**: Scrapes application info from Google Play Store
- **VirusTotal Analysis** (requires VirusTotal API key):
  - Malware detection stats with highlighted red flags
  - Popular threat classification (e.g., "trojan.pegasus/chrysaor")
  - File reputation score
  - Community detection (Sigma, YARA, IDS)
  - File tags and characteristics
- **VirusTotal Androguard** (automatic with VirusTotal API key):
  - Complete APK structure analysis
  - Activities, Services, Providers, Receivers
  - Libraries and SDK versions
  - Dangerous permissions highlighted in red
- **Koodous Integration** (requires Koodous API key):
  - Malware detection status
  - Community rating and trust score
  - Positive/Negative votes
  - Repository information when available

### Output & Export
- **Enhanced Terminal Output**: Colored results with **bold red warnings** for malware indicators
- **JSON Export**: Pretty-printed analysis export including all VirusTotal/Koodous data

## Installation

### From GitHub Releases

Download the pre-compiled binary for your system from the [Releases](https://github.com/andpalmier/apkingo/releases) page.

### From Source

```bash
go install github.com/andpalmier/apkingo/cmd/apkingo@latest
```

### From Homebrew

```bash
brew tap andpalmier/tap
brew install apkingo
```

## Usage

### Using Docker (Recommended)

You can run **apkingo** directly using Docker without installing Go or downloading binaries.

```bash
# Analyze an APK (mount the directory containing the APK)
docker run --rm -v $(pwd):/mnt ghcr.io/andpalmier/apkingo -apk /mnt/target.apk

# Analyze an XAPK file
docker run --rm -v $(pwd):/mnt ghcr.io/andpalmier/apkingo -apk /mnt/app.xapk

# Analyze all APKs in a directory
docker run --rm -v $(pwd):/mnt ghcr.io/andpalmier/apkingo -dir /mnt

# Analyze and export JSON report
docker run --rm -v $(pwd):/mnt ghcr.io/andpalmier/apkingo -apk /mnt/target.apk -json /mnt/report.json
```

### CLI Usage

```bash
# Analyze a single APK
apkingo -apk <path_to_apk>

# Analyze an XAPK/APKS file
apkingo -apk <path_to_xapk>

# Analyze all APKs in a directory
apkingo -dir <path_to_directory>

# Analyze with API keys and export JSON
apkingo -apk <path_to_apk> -vtapi <VT_KEY> -kapi <KOODOUS_KEY> -json report.json
```

### API Keys

For enhanced analysis, you can provide API keys for VirusTotal and Koodous either via command-line flags or environment variables:

**Environment Variables (Recommended):**
```bash
export VT_API_KEY="your_virustotal_api_key"
export KOODOUS_API_KEY="your_koodous_api_key"
apkingo -apk <path_to_apk>
```

**Command-Line Flags:**
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
| `-vtapi` | VirusTotal API key (can also use `VT_API_KEY` env var) |
| `-kapi` | Koodous API key (can also use `KOODOUS_API_KEY` env var) |
| `-no-play-store` | Skip Play Store scraping for offline analysis |
| `-vtupload` | Upload the APK to VirusTotal after analysis (interactive prompt) |

### Example

```bash
apkingo -apk <path_to_apk>
```

## Screenshot

apkingo analyzing an Android malware:
<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/img/screen_malware.png?raw=true" />
</p>

## 3rd party libraries and API documentation 

- shogo82148/androidbinary: [GitHub repo](https://github.com/shogo82148/androidbinary) and [Go reference](https://pkg.go.dev/github.com/shogo82148/androidbinary)
- avast/apkverifier: [GitHub repo](https://github.com/avast/apkverifier) and [Go reference](https://pkg.go.dev/github.com/avast/apkverifier)
- fatih/color: [GitHub repo](https://github.com/fatih/color) and [Go reference](https://pkg.go.dev/github.com/fatih/color)
- n0madic/google-play-scraper: [GitHub repo](https://github.com/n0madic/google-play-scraper) and [Go reference](https://pkg.go.dev/github.com/n0madic/google-play-scraper)
- parnurzeal/gorequest: [GitHub repo](https://github.com/parnurzeal/gorequest) and [Go reference](https://pkg.go.dev/github.com/parnurzeal/gorequest)
- VirusTotal/vt-go: [GitHub repo](https://github.com/VirusTotal/vt-go) and [Go reference](https://pkg.go.dev/github.com/VirusTotal/vt-go)
- [VirusTotal API documentation](https://docs.virustotal.com/reference/overview)
- [Koodous API documentation](https://docs.koodous.com/api/apks.html)
