package main

/*
	配合go generate自动生成插件的注册代码
*/

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var plugin_register = `// go generate build
package app

import (
	"github.com/abingzo/bups/iocc"
	%s
)

func PluginRegister() {
	%s
}
`

func main() {
	Generate("../")
	GenerateFile("./plugin_register.go")
}

func Generate(initPath string) {
	startF := "plugins"
	importAppend := ""
	funcAppend := ""
	// os.FileInfo 是为了兼容go 1.16以下的版本
	err := filepath.Walk(initPath+startF, func(path string, info os.FileInfo, err error) error {
		// 遍历到自己
		if info.Name() == startF {
			return nil
		}
		// 非文件
		if !info.IsDir() {
			return nil
		}
		importAppend += fmt.Sprintf("%sgithub.com/abingzo/bups/%s/%s%s\n\t",`"`,startF,info.Name(),`"`)
		funcAppend += fmt.Sprintf("iocc.RegisterPlugin(%s.New)\n\t", info.Name())
		return nil
	})
	if err != nil {
		panic(err)
	}
	plugin_register = fmt.Sprintf(plugin_register,importAppend,funcAppend)
}

func GenerateFile(path string) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writeN, err := file.Write([]byte(plugin_register))
	if err != nil {
		panic(err)
	}
	if writeN != len(plugin_register) {
		panic(errors.New("write bytes not equal"))
	}
}
