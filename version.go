//+build linux darwin

package main

/*
	用于记录程序的版本信息
*/

import (
	"fmt"
	"runtime"
)

var (
	version      string
	gitBranch    string
	gitTag       string
	gitCommit    string
	gitTreeState string
	buildDate    string
)

type Info struct {
	Version      string `json:"version"`
	GitBranch    string `json:"gitBranch"`
	GitTag       string `json:"gitTag"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// GetInfo 注意：该函数接受的是编译时-ldflags 传入的数据
func GetInfo() *Info {
	return &Info{
		Version:      version,
		GitBranch:    gitBranch,
		GitTag:       gitTag,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}


