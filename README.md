# apkingo

<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/apkingo.png?raw=true" />
  <p align="center">
    <a href="https://github.com/andpalmier/apkingo/blob/master/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/andpalmier/apkingo"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/andpalmier/apkingo?style=flat-square"></a>
    <a href="https://twitter.com/intent/follow?screen_name=andpalmier"><img src="https://img.shields.io/twitter/follow/andpalmier?style=social&logo=twitter" alt="follow on Twitter"></a>
  </p>
</p>

apkingo is a tool written in Go to get detailed information about an apk file. apkingo will explore the given file to get details on the apk, such as package name, target SDK, permissions and metadata. Use the `-playstore` flag to search the app in the Play Store and retrieve additional information. Use the flag `-cert` to display certificate information contained in the apk.  

## Usage

After downloading the repository, navigate into the directory and build the project with `make apkingo`. This will create a folder `build`, containing an executable called `apkingo`. You can then run the executable with the following flags:

- `-apk` to specify the path to the apk file (**required**)
- `-cert` for printing the certificate information retrieved in the apk file (**sometimes printing a conversion error, but it's still working!**)
- `-playstore` for searching the app in the Play Store by its package name

## Example

<p align="center">
  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/screen_snapseed.png?raw=true" />

  <img alt="apkingo" src="https://github.com/andpalmier/apkingo/blob/main/screen_f-droid.png?raw=true" />
</p>

## Libraries used

- [shogo81148/androidbinary](github.com/shogo82148/androidbinary)
- [avast/apkverifier](github.com/avast/apkverifier)
- [n0madic/google-play-scraper](github.com/n0madic/google-play-scraper)