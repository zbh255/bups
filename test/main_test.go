package test

import (
	"errors"
	"github.com/abingzo/bups/app"
	"github.com/abingzo/bups/common/config"
	"github.com/zbh255/bilog"
	"io/ioutil"
	"os"
	"sync"
	"testing"
)

// 测试操作配置文件的函数
func TestCFGOptions(t *testing.T) {
	cfg := &app.CFG{}
	file, err := os.OpenFile("./config.toml", os.O_RDWR, 0777)
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
	_ = config.Read(cfg)
}

func TestNativeCfg(t *testing.T) {
	file, err := os.OpenFile("./config.toml", os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	file2, err := os.OpenFile("./config.toml", os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	// 读取配置文件
	bytes2, _ := ioutil.ReadAll(file2)
	cfg := &app.CFG{}
	cfg.Open(file)
	bytes1, _ := ioutil.ReadAll(cfg)
	if !(string(bytes2) == string(bytes1)) {
		t.Fatal(errors.New("write config file is not equal read config file"))
	}
}

// 测试多个日志器共享一个底层文件
func TestMultiLogger(t *testing.T) {
	logFile, err := os.OpenFile("./bups.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	daemonLog := bilog.NewLogger(logFile,bilog.ERROR)
	backupLog := bilog.NewLogger(logFile,bilog.ERROR)
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
