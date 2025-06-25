package version

import "runtime"

var (
	version 	string
	commit 		string
	buildDate 	string
)

type Info struct {
	Version		string
	Commit		string
	Date		string
	GoVersion	string
	OS			string
	Arch		string
}

func GetInfo() Info {
	return Info{
		Version: version,
		Commit: commit,
		Date: buildDate,
		GoVersion: runtime.Version(),
		OS: runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

func SetBuildInfo(v, c, d string) {
	version = v
	commit = c
	buildDate = d
}