package main

import "fmt"

type Version struct {
	GitHash   string
	GoVersion string
	BuildTime string
}

var (
	__Version__ Version
	gitHash     string
	goVersion   string
	buildTime   string
)

func init() {
	__Version__ = Version{
		GitHash:   gitHash,
		GoVersion: goVersion,
		BuildTime: buildTime,
	}
}

func (ver Version) String() string {
	return fmt.Sprintf("GitHash: %s, GoVersion: %s, BuildTime: %s", __Version__.GitHash, __Version__.GoVersion, __Version__.BuildTime)
}
