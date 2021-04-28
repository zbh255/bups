package main

import (
	"github.com/mengzushan/bups/app"
	cf "github.com/mengzushan/bups/common/conf"
	this "github.com/mengzushan/bups/common/error"
	"github.com/mengzushan/bups/common/info"
	"github.com/mengzushan/bups/common/logger"
	"github.com/mengzushan/bups/utils"
	"github.com/mengzushan/bups/web"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conf := cf.InitConfig()
	// webui配置项为开启则进入
	log, err := logger.Std(nil)
	if err != this.Nil {
		panic(err)
	}
	defer log.Close()
	if conf.WebConfig.Switch == "on" {
		go web.Run()
	}
	// 接收退出信号是清理cache文件
	go func() {
		sig := make(chan os.Signal)
		// Ctrl-C Ctrl-/ quit stop
		signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGTSTP)
		<-sig
		pwd, _ := os.Getwd()
		err = utils.CleanUpCache(pwd + "/cache/backup")
		if err != this.Nil {
			log.StdErrorLog(err.Error())
		}
		os.Exit(0)
	}()
	// 管道用于控制函数调用
	callChan := make(chan byte)
	go func() {
		for true {
			select {
			case <-callChan:
				app.TimerTask()
			}
		}
	}()
	/*
		异常退出的情况与第一次启动的情况
		读取app_info.json中的timer字段，并更新
		为避免误差设置提前2000毫秒
	*/
	appInfo := info.GetAppInfo()
	if appInfo.Timer+2000 < time.Now().Unix() && info.RunTimeNum == 0 {
		// 重新设置下一次的时间
		appInfo.Timer = appInfo.Timer + int64(time.Duration(conf.SaveTime)*time.Hour)
		err = info.SetAppInfo(appInfo)
		if err != this.Nil {
			log.StdErrorLog(err.Error())
		} else {
			log.StdInfoLog("Reset the timing task successfully")
			callChan <- 0
			info.RunTimeNum++
		}
	}
	// 根据配置项创建定时任务
	// 配置备份的时间按小时算
	// 卡死循环
	for true {
		oldTime := appInfo.Timer
		newTime := time.Now().Unix()
		SaveTime := time.NewTimer(time.Duration(oldTime - newTime))
		select {
		case <-SaveTime.C:
			callChan <- 0
			info.RunTimeNum++
		}
	}
}
