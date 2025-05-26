package version

import (
	"fmt"
	"runtime"
)

var (
	Version = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func Info() string {
	return fmt.Sprintf(
		"Version: %s\nGit Commit: %s\nBuild Date: %s\nGo Version: %s\nOS/Arch: %s/%s",
		Version,
		GitCommit,
		BuildDate,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}