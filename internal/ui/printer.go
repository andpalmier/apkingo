// Package ui provides terminal user interface components for apkingo.
package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

type Printer struct {
	out     io.Writer
	tw      *tabwriter.Writer
	noColor bool
	cyan    *color.Color
	yellow  *color.Color
	red     *color.Color
	italic  *color.Color
}

func NewPrinter(noColor bool) *Printer {
	color.NoColor = noColor
	p := &Printer{
		out:     os.Stdout,
		noColor: noColor,
		cyan:    color.New(color.FgCyan),
		yellow:  color.New(color.FgYellow),
		red:     color.New(color.FgRed),
		italic:  color.New(color.FgWhite, color.Italic),
	}
	p.tw = tabwriter.NewWriter(p.out, 0, 0, 2, ' ', 0)
	return p
}

func (p *Printer) Flush() {
	p.tw.Flush()
}

func (p *Printer) GetRed() *color.Color {
	return p.red
}

func (p *Printer) GetCyan() *color.Color {
	return p.cyan
}

func (p *Printer) GetTabWriter() *tabwriter.Writer {
	return p.tw
}

func (p *Printer) GetOut() io.Writer {
	return p.out
}

func (p *Printer) PrintSectionHeader(title string) {
	p.Flush()
	fmt.Fprintln(p.out)
	p.yellow.Fprintf(p.out, "─── %s ───\n", strings.ToUpper(title))
}

func (p *Printer) PrintKV(key, value string) {
	fmt.Fprintf(p.tw, "%s:\t", key)
	if value != "" && value != "<nil>" && value != "[]" {
		p.cyan.Fprintln(p.tw, value)
	} else {
		p.italic.Fprintln(p.tw, "not found")
	}
}

func (p *Printer) PrintKVRed(key, value string) {
	fmt.Fprintf(p.tw, "%s:\t", key)
	if value != "" && value != "<nil>" && value != "[]" {
		p.red.Fprintln(p.tw, value)
	} else {
		p.italic.Fprintln(p.tw, "not found")
	}
}

func (p *Printer) PrintKVRedBold(key, value string) {
	fmt.Fprintf(p.tw, "%s:\t", key)
	if value != "" && value != "<nil>" && value != "[]" {
		p.red.Add(color.Bold).Fprintln(p.tw, value)
	} else {
		p.italic.Fprintln(p.tw, "not found")
	}
}

func (p *Printer) PrintList(items []string) {
	if len(items) == 0 {
		p.italic.Fprintln(p.out, "none")
		return
	}
	for _, item := range items {
		fmt.Fprintf(p.out, " - %s\n", item)
	}
}

func (p *Printer) PrintBanner() {
	banner := `
 ┌─┐┌─┐┬┌─┬┌┐┌┌─┐┌─┐
 ├─┤├─┘├┴┐│││││ ┬│ │
 ┴ ┴┴  ┴ ┴┴┘└┘└─┘└─┘
by @andpalmier
`
	p.cyan.Fprintln(p.out, banner)
}

func (p *Printer) PrintTitle(title string) {
	p.Flush() // Flush previous section
	p.yellow.Fprintf(p.out, "\n* %s\n", title)
}

func (p *Printer) PrintLabelValue(label, value string) {
	fmt.Fprintf(p.tw, "%s:\t", label)
	if value != "" && value != "<nil>" && value != "[]" {
		p.cyan.Fprintln(p.tw, value)
	} else {
		p.italic.Fprintln(p.tw, "not found")
	}
}

func (p *Printer) PrintLabelValueIndent(label, value string) {
	fmt.Fprintf(p.tw, "\t%s:\t", label)
	if value != "" && value != "<nil>" && value != "[]" {
		p.cyan.Fprintln(p.tw, value)
	} else {
		p.italic.Fprintln(p.tw, "not found")
	}
}

func (p *Printer) PrintError(msg string) {
	p.Flush()
	p.red.Fprintln(p.out, msg)
}

func (p *Printer) PrintSuccess(msg string) {
	p.Flush()
	p.cyan.Fprintln(p.out, msg)
}

func (p *Printer) PrintText(msg string) {
	p.Flush()
	fmt.Fprintln(p.out, msg)
}

func (p *Printer) PrintItalic(msg string) {
	p.Flush()
	p.italic.Fprintln(p.out, msg)
}

func (p *Printer) Printf(format string, a ...interface{}) {
	fmt.Fprintf(p.tw, format, a...)
}

func (p *Printer) Println(a ...interface{}) {
	fmt.Fprintln(p.tw, a...)
}
