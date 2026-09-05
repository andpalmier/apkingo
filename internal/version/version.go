// Package version carries the build information stamped into the binary.
package version

import (
	"fmt"
	"runtime/debug"
)

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

// shortHashLength matches the abbreviation git uses by default, so a commit
// recovered from the build info reads the same as one passed by -ldflags.
const shortHashLength = 7

// init recovers the build information that -ldflags did not supply.
//
// "go install module@version" applies no -ldflags at all, so a binary installed
// the way the README describes used to report itself as "dev" with no way to
// tell which release it came from. The toolchain does record the module version
// in that case, and records the revision and commit time instead when building
// inside a checkout, so between the two the binary can always say where it came
// from. Values that -ldflags did set are left alone.
func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	// "(devel)" is what a build from a working tree reports, which says less
	// than the revision below.
	if Version == "dev" && info.Main.Version != "" && info.Main.Version != "(devel)" {
		Version = info.Main.Version
	}

	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			if Commit == "none" && len(s.Value) >= shortHashLength {
				Commit = s.Value[:shortHashLength]
			}
		case "vcs.time":
			if Date == "unknown" {
				Date = s.Value
			}
		}
	}
}

// String renders the build information for the -version flag. The commit and
// build date are omitted when they hold their unstamped defaults, so a binary
// built without either source of information reports what it knows rather than
// three placeholders.
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
