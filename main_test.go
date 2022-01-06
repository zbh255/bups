package main

import (
	"github.com/abingzo/bups/common/config"
	"github.com/abingzo/bups/common/logger"
	"io/ioutil"
	"os"
	"sync"
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

// 测试多个日志器共享一个底层文件
func TestMultiLogger(t *testing.T) {
	logFile, err := os.OpenFile("./test/bups.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	daemonLog := logger.New(logFile, logger.ERROR)
	backupLog := logger.New(logFile, logger.ERROR)
	wg := sync.WaitGroup{}
	wg.Add(400)
	go func() {
		for i := 0; i < 200; i++ {
			daemonLog.Info("header")
			wg.Done()
		}
	}()
	go func() {
		for i := 0; i < 200; i++ {
			backupLog.Info("hello world again")
			wg.Done()
		}
	}()
	wg.Wait()
}
