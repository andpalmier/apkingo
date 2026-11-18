package utils_test

import (
	"strings"
	"testing"

	"github.com/andpalmier/apkingo/v3/internal/utils"
)

func TestSanitizeText(t *testing.T) {
	got := utils.SanitizeText("safe\x1b]0;pwned\x07\u202eevil")
	if strings.ContainsAny(got, "\x1b\x07") {
		t.Fatalf("Control characters survived sanitization: %q", got)
	}
	if !strings.Contains(got, "pwned") || !strings.Contains(got, "evil") {
		t.Errorf("Sanitization removed printable text: %q", got)
	}
}
