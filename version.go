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
	return fmt.Sprintf("GitSha: %s, GoVersion: %s, BuildTime: %s", __Version__.GitSha, __Version__.GoVersion, __Version__.BuildTime)
}
