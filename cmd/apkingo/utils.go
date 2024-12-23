package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Android SDK version names
var androidName = map[int]string{
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
	34: "Android 14",
	35: "Android 15",
	36: "Android 16",
}

// getFileName retrieves a file name from a file path
func getFileName(path string) string {
	return filepath.Base(path)
}

// printer helps print some values
func printer(s string) {
	if s != "" {
		cyan.Printf("%s\n", s)
	} else {
		italic.Println()
	}
}

// getAPIKey retrieves an API key from the environment variable
func getAPIKey(flagValue, envVar, msg string) string {
	if flagValue == "" {
		flagValue = os.Getenv(envVar)
		if flagValue == "" {
			fmt.Println(msg)
		}
	}
	return flagValue
}

// logError logs an error message if the error is not nil
func logError(msg string, err error) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
	}
}
