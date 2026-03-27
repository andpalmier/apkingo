// Package constants holds shared constants used across the application.
package constants

import "time"

// Sentinel values used when checking if values are empty/missing.
const (
	NilValue   = "<nil>"
	EmptySlice = "[]"
)

// Default values for configuration.
const (
	// DefaultCountry is the default country code for Play Store lookups.
	DefaultCountry = "us"
	// DefaultHTTPTimeout is the default timeout for HTTP requests.
	DefaultHTTPTimeout = 30 * time.Second
)
