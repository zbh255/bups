package main

import (
	"github.com/abingzo/bups/common/plugin"
	"io"
	"io/ioutil"
)

const (
	Name = "cloud-upload"
	Type = plugin.BCallBack
)

var Support = []int{plugin.SupportLogger,plugin.SupportConfigRead}

var StdOut io.Writer

var conf []byte

/*
	上传文件的插件
*/

func SetStdout(out io.Writer) {
	StdOut = out
}

func GetName() string {
	return Name
}

func GetType() plugin.Type {
	return Type
}

func GetSupport() []int {
	return Support
}

func ConfRead(reader io.Reader) {
	conf,_ = ioutil.ReadAll(reader)
}

func ConfWrite(writer io.Writer) {}

// Start 启动函数
func Start()  {
	_, _ = StdOut.Write([]byte("插件测试启动开始\n"))
	_, _ = StdOut.Write([]byte(string(conf) + "\n"))
	_, _ = StdOut.Write([]byte("插件测试启动完成\n"))
}