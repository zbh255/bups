package main

import (
	"github.com/mengzushan/bups/app"
	this "github.com/mengzushan/bups/common/error"
	"os"
)

func main() {
	// 检查是否有参数
	if len(os.Args) > 1 {
		argErr := app.DeCommandArgs()
		if argErr != this.Nil {
			os.Exit(2)
		} else {
			os.Exit(0)
		}
	}
	// 调用函数调用
	app.MainDisPatch()
}

