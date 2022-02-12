# apkingo

<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/apkingo.png?raw=true" />
  <p align="center">
    <a href="https://github.com/andpalmier/apkingo/blob/main/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/andpalmier/apkingo"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/andpalmier/apkingo?style=flat-square"></a>
    <a href="https://twitter.com/intent/follow?screen_name=andpalmier"><img src="https://img.shields.io/twitter/follow/andpalmier?style=social&logo=twitter" alt="follow on Twitter"></a>
  </p>
</p>

apkingo is a tool written in Go to get detailed information about an apk file. apkingo will explore the given file to get details on the apk, such as package name, target SDK, permissions, metadata, certificate serial and issuer. The tool will also retrieve information about the 
specified apk from the Play Store and detect if it is malicious using [Koodous](https://koodous.com/).

## Usage

After downloading the repository, navigate into the directory and build the project with `make apkingo`. This will create a folder `build`, containing an executable called `apkingo`. You can then run the executable with the following flags:

- `-apk` to specify the path to the apk file (**required**)
- `-json`	to specify the path of the json file where the results will be exported


## Screenshots

apkingo analyzing snapseed:
<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/screen_snapseed.png?raw=true" />
</p><br><br>

apkingo analyzing F-Droid:
<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/screen_f-droid.png?raw=true" />
</p><br><br>

apkingo analyzing an android malware:
<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/screen_malware.png?raw=true" />
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/screen_malware1.png?raw=true" />
</p>

## Non-standard libraries used

- [shogo81148/androidbinary](https://github.com/shogo82148/androidbinary)
- [avast/apkverifier](https://github.com/avast/apkverifier)
- [n0madic/google-play-scraper](https://github.com/n0madic/google-play-scraper)
-	[parnurzeal/gorequest](https://github.com/parnurzeal/gorequest)
-	[fatih/color](https://github.com/fatih/color)
