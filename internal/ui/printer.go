// Package ui provides terminal user interface components for apkingo.
package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/andpalmier/apkingo/internal/constants"
	"github.com/fatih/color"
)

// Printer provides colored terminal output for APK analysis results.
type Printer struct {
	out     io.Writer
	tw      *tabwriter.Writer
	cyan    *color.Color
	yellow  *color.Color
	red     *color.Color
	redBold *color.Color
	italic  *color.Color
}

// NewPrinter creates a new Printer for terminal output.
// It initializes all color objects including bold variants for warnings.
func NewPrinter() *Printer {
	p := &Printer{
		out:     os.Stdout,
		cyan:    color.New(color.FgCyan),
		yellow:  color.New(color.FgYellow),
		red:     color.New(color.FgRed),
		redBold: color.New(color.FgRed, color.Bold),
		italic:  color.New(color.FgWhite, color.Italic),
	}
	p.tw = tabwriter.NewWriter(p.out, 0, 0, 2, ' ', 0)
	return p
}

func (p *Printer) Flush() {
	_ = p.tw.Flush()
}

func (p *Printer) GetRed() *color.Color {
	return p.red
}

// GetRedBold returns a red color with bold styling for warnings.
func (p *Printer) GetRedBold() *color.Color {
	return p.redBold
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
	_, _ = fmt.Fprintln(p.out)
	_, _ = p.yellow.Fprintf(p.out, "─── %s ───\n", strings.ToUpper(title))
}

func (p *Printer) PrintKV(key, value string) {
	_, _ = fmt.Fprintf(p.tw, "%s:\t", key)
	if value != "" && value != constants.NilValue && value != constants.EmptySlice {
		_, _ = p.cyan.Fprintln(p.tw, value)
	} else {
		_, _ = p.italic.Fprintln(p.tw, "not found")
	}
}

func (p *Printer) PrintKVRed(key, value string) {
	_, _ = fmt.Fprintf(p.tw, "%s:\t", key)
	if value != "" && value != constants.NilValue && value != constants.EmptySlice {
		_, _ = p.red.Fprintln(p.tw, value)
	} else {
		_, _ = p.italic.Fprintln(p.tw, "not found")
	}
}

func (p *Printer) PrintKVRedBold(key, value string) {
	_, _ = fmt.Fprintf(p.tw, "%s:\t", key)
	if value != "" && value != constants.NilValue && value != constants.EmptySlice {
		_, _ = p.redBold.Fprintln(p.tw, value)
	} else {
		_, _ = p.italic.Fprintln(p.tw, "not found")
	}
}

func (p *Printer) PrintList(items []string) {
	if len(items) == 0 {
		_, _ = p.italic.Fprintln(p.out, "none")
		return
	}
	for _, item := range items {
		_, _ = fmt.Fprintf(p.out, " - %s\n", item)
	}
}

func (p *Printer) PrintBanner() {
	banner := `
  ┌─┐┌─┐┬┌─┬┌┐┌┌─┐┌─┐
  ├─┤├─┘├┴┐│││││ ┬│ │
  ┴ ┴┴  ┴ ┴┴┘└┘└─┘└─┘
by @andpalmier
`
	_, _ = p.cyan.Fprintln(p.out, banner)
}

func (p *Printer) PrintText(msg string) {
	p.Flush()
	_, _ = fmt.Fprintln(p.out, msg)
}

func (p *Printer) PrintItalic(msg string) {
	p.Flush()
	_, _ = p.italic.Fprintln(p.out, msg)
}

func (p *Printer) Printf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(p.tw, format, a...)
}

func (p *Printer) Println(a ...interface{}) {
	_, _ = fmt.Fprintln(p.tw, a...)
}
