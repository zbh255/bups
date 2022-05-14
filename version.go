// +build linux darwin windows

package main

/*
	用于记录程序的版本信息
*/

import (
	"fmt"
	"github.com/abingzo/bups/app"
	"runtime"
)

var (
	version      string
	gitBranch    string
	gitCommit    string
	gitTreeState string
	buildDate    string
)


// GetInfo 注意：该函数接受的是编译时-ldflags 传入的数据
// 返回指针是为了兼容标准库的json包
func GetInfo() *app.Info {
	return &app.Info{
		Version:      version,
		GitBranch:    gitBranch,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
