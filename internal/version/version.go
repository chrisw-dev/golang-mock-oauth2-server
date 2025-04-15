// Package version provides version information for the application
package version

import (
	"fmt"
	"runtime"
)

// Build information set by ldflags during build
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = ""
	GoVersion = runtime.Version()
)

// Info contains version information
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
}

// GetVersion returns version information
func GetVersion() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
		GoVersion: GoVersion,
	}
}

// String returns a formatted string with version information
func (i Info) String() string {
	return fmt.Sprintf("Version: %s\nCommit: %s\nBuildDate: %s\nGo: %s",
		i.Version, i.Commit, i.BuildDate, i.GoVersion)
}
