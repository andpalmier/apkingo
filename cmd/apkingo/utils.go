package main

import "github.com/fatih/color"

// colors to improve readability
var italic = color.New(color.FgWhite, color.Italic)
var yellow = color.New(color.FgYellow)
var cyan = color.New(color.FgCyan)
var red = color.New(color.FgRed)

// printBanner() - like the cool kids
func printBanner() {
	var Banner string = `
	┌─┐┌─┐┬┌─┬┌┐┌┌─┐┌─┐
	├─┤├─┘├┴┐│││││ ┬│ │
	┴ ┴┴  ┴ ┴┴┘└┘└─┘└─┘
by @andpalmier
`
	cyan.Println(Banner)
}

// printer(string) - help print some values
func printer(s string) {
	if s != "" {
		cyan.Printf("%s\n", s)
	} else {
		italic.Println("not found")
	}
}

// mapping SDK to Android version
var androidname = map[int]string{
	0:  "Not found",
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
