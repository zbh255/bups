package web_config

import (
	"github.com/abingzo/bups/common/logger"
	"github.com/abingzo/bups/common/plugin"
	"os"
	"strings"
	"testing"
)

func TestConfWebReadWrite(t *testing.T) {
	webConfig := New()
	t.Log(os.Args)
	file, err := os.OpenFile("../../config/dev/config.toml", os.O_RDWR|os.O_SYNC, 0777)
	if err != nil {
		panic(err)
	}
	rawSource := new(plugin.Source)
	rawSource.RawConfig = file
	rawSource.StdLog = logger.New(os.Stdout,logger.PANIC)
	rawSource.AccessLog = logger.New(os.Stdout,logger.DEBUG)
	rawSource.ErrorLog = logger.New(os.Stderr,logger.PANIC)

	webConfig.SetSource(rawSource)
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
