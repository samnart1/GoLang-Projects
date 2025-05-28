package version

import (
	"fmt"
	"runtime"
)

var (
	Version = "dev"
	BuildDate = "unknown"
)

func GetVersion() string {
	return fmt.Sprintf("Version: %s (%s/%s)\nBuild Date: %s\n", Version, runtime.GOOS, runtime.GOARCH, BuildDate)
}