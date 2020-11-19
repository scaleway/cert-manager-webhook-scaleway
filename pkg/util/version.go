package util

import (
	"fmt"
	"runtime"
)

// These are set during build time via -ldflags
var (
	version   = "0.0.1+dev"
	gitCommit string
	buildDate string
)

// VersionInfo represents the current running version
type VersionInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
	Compiler  string `json:"compiler"`
	Platform  string `json:"platform"`
}

// GetVersion returns the current running version
func GetVersion() VersionInfo {
	return VersionInfo{
		Version:   version,
		GitCommit: gitCommit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
