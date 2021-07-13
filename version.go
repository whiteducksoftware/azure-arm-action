package main

import "fmt"

type Version struct {
	GitSha    string
	GoVersion string
	BuildTime string
}

var (
	__Version__ Version
	gitSha      string
	goVersion   string
	buildTime   string
)

func init() {
	__Version__ = Version{
		GitSha:    gitSha,
		GoVersion: goVersion,
		BuildTime: buildTime,
	}
}

func (ver Version) String() string {
	return fmt.Sprintf("Git SHA: %s, Go Version: %s, Build Time: %s", ver.GitSha, ver.GoVersion, ver.BuildTime)
}
