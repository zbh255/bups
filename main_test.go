package main

import (
	"github.com/abingzo/bups/common/config"
	"github.com/abingzo/bups/common/plugin"
	"io/ioutil"
	"os"
	"testing"
)

// 测试操作配置文件的函数
func TestCFGOptions(t *testing.T) {
	cfg := &CFG{}
	file, err := os.OpenFile("./config/dev/config.toml", os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	cfg.Open(file)
	_, err = ioutil.ReadAll(cfg)
	if err != nil {
		panic(err)
	}
	_, err = ioutil.ReadAll(cfg)
	if err != nil {
		panic(err)
	}
	// 读取配置文件
	conf := config.Read(cfg)
	t.Log(conf.Plugin)
}

func TestNativeCfg(t *testing.T) {
	file, err := os.OpenFile("./config/dev/config.toml", os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	file2, err := os.OpenFile("./config/dev/config.toml", os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	// 读取配置文件
	bytes2, _ := ioutil.ReadAll(file2)
	cfg := &CFG{}
	cfg.Open(file)
	bytes1, _ := ioutil.ReadAll(cfg)
	t.Log(string(bytes2) == string(bytes1))
}

// 测试插件的加载和管理
func TestPluginLoadAndManager(t *testing.T) {
	ctx := LoadPlugin("./build_release/plugins",
		"./build_release/log/bups.log", "./build_release/config/config.toml")
	ctx.SetState(plugin.Init)
}
