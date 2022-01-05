package web_config

import (
	"github.com/abingzo/bups/common/logger"
	"io"
	"os"
	"strings"
	"testing"
)

func TestConfWebReadWrite(t *testing.T) {
	webConfig := New()
	webConfig.SetStdout(os.Stdout)
	webConfig.SetLogOut(logger.New(os.Stdout, "Plugin.web_config"))
	t.Log(os.Args)
	file, err := os.OpenFile("../../config/dev/config.toml", os.O_RDWR|os.O_SYNC, 0777)
	if err != nil {
		panic(err)
	}
	webConfig.ConfRead(io.Reader(file))
	webConfig.ConfWrite(file)
	// TODO: 处理命令行参数
	pluginArgs := os.Args[len(os.Args)-1]
	pluginArgs = pluginArgs[1 : len(pluginArgs)-1]
	args := strings.Split(pluginArgs, " ")
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	args = append([]string{pwd + "/" + webConfig.GetName()}, args...)
	t.Log(args)
	t.Log(len(args))
	webConfig.Start(args)
}
