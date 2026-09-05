// Package version carries the build information stamped into the binary.
package version

import "fmt"

// Injected at link time by GoReleaser with -ldflags -X. The linker drops a -X
// naming a symbol that does not exist without reporting anything and without
// failing the build, so these names and the paths in .goreleaser.yaml have to
// be kept in step. CI builds a release binary and asserts the values arrived,
// because reading the two files side by side is not enough to catch a typo.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// String renders the build information for the -version flag. The commit and
// build date are omitted when they hold their unstamped defaults, so a locally
// built binary reports what it knows rather than three placeholders.
func String() string {
	s := fmt.Sprintf("apkingo version %s", Version)
	if Commit != "none" {
		s += fmt.Sprintf("\n  commit: %s", Commit)
	}
	if Date != "unknown" {
		s += fmt.Sprintf("\n  built: %s", Date)
	}
	return s
}
