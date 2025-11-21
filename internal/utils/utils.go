// Package utils provides utility functions for the apkingo application.
package utils

import (
	"log"
	"path/filepath"
)

// Android SDK version names
var AndroidName = map[int]string{
	1:  "Android 1.0",
	2:  "Android 1.1",
	3:  "Android 1.5 Cupcake",
	4:  "Android 1.6 Donut",
	5:  "Android 2.0 Eclair",
	6:  "Android 2.0.1 Eclair",
	7:  "Android 2.1 Eclair",
	8:  "Android 2.2 Froyo",
	9:  "Android 2.3 Gingerbread",
	10: "Android 2.3.3 Gingerbread",
	11: "Android 3.0 Honeycomb",
	12: "Android 3.1 Honeycomb",
	13: "Android 3.2 Honeycomb",
	14: "Android 4.0 Ice Cream Sandwich",
	15: "Android 4.0.3 Ice Cream Sandwich",
	16: "Android 4.1 Jelly Bean",
	17: "Android 4.2 Jelly Bean",
	18: "Android 4.3 Jelly Bean",
	19: "Android 4.4 KitKat",
	20: "Android 4.4W KitKat",
	21: "Android 5.0 Lollipop",
	22: "Android 5.1 Lollipop",
	23: "Android 6.0 Marshmallow",
	24: "Android 7.0 Nougat",
	25: "Android 7.1 Nougat",
	26: "Android 8.0 Oreo",
	27: "Android 8.1 Oreo",
	28: "Android 9 Pie",
	29: "Android 10",
	30: "Android 11",
	31: "Android 12",
	32: "Android 12L",
	33: "Android 13",
	34: "Android 14",
	35: "Android 15",
	36: "Android 16",
}

// GetFileName retrieves a file name from a file path
func GetFileName(path string) string {
	return filepath.Base(path)
}

// LogError logs an error message if the error is not nil
func LogError(msg string, err error) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}
