package main

import (
	"github.com/abingzo/bups/common/plugin"
	"io"
	"io/ioutil"
	"os"
)

const (
	Name = "cloud-upload"
	Type = plugin.BCallBack
)

var Support = []int{plugin.SupportLogger,plugin.SupportConfigRead}

var StdOut io.Writer = os.Stdout

var conf []byte

func New() plugin.Plugin {
	return &Upload{
		Name: Name,
		Type: Type,
		Support: Support,
		stdout: StdOut,
	}
}

/*
	上传文件的插件
*/

type Upload struct {
	plugin.Plugin
	Name string
	Type plugin.Type
	Support []int
	stdout io.Writer
}

func (u *Upload) SetStdout(out io.Writer) {
	u.stdout = out
}

func (u *Upload) GetName() string {
	return u.Name
}

func (u *Upload) GetType() plugin.Type {
	return u.Type
}

func (u *Upload) GetSupport() []int {
	return u.Support
}

func (u *Upload) ConfRead(reader io.Reader) {
	conf,_ = ioutil.ReadAll(reader)
}

func (u *Upload) ConfWrite(writer io.Writer) {
	_, _ = writer.Write(conf)
}

// Start 启动函数
func (u *Upload) Start(args []string)  {
	if args == nil {
		_, _ = StdOut.Write([]byte("插件测试启动开始\n"))
		_, _ = StdOut.Write([]byte(string(conf) + "\n"))
		_, _ = StdOut.Write([]byte("插件测试启动完成\n"))
	}
}
