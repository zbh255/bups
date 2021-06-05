package test

import (
	"github.com/mengzushan/bups/app"
	"github.com/mengzushan/bups/common/conf"
	this "github.com/mengzushan/bups/common/error"
	"os"
	"testing"
)


func Test_Check_Config(t *testing.T) {
	pathHead,_ := os.Getwd()
	path := pathHead + "/conf/dev/app.conf.toml"
	conf.TestOnFilePath = path
	config := conf.InitConfig()
	err := app.CheckTomlConfig(config)
	if err != this.Nil {
		t.Error("测试失败: " + err.Error())
	} else {
		t.Log("测试成功")
	}
}