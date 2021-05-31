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
	// 检查配置的正确性
	err = app.CheckTomlConfig(conf)
	if err != this.Nil {
		log.StdErrorLog(err.Error())
		panic(err)
	}
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
	// 备份信息
	backupInfo := make(chan byte)
	go func() {
		for true {
			select {
			case <-callChan:
				taskInfo,err := app.TimerTask(conf)
				if err != this.Nil {
					log.StdErrorLog(err.Error())
				} else {
					log.StdInfoLog(taskInfo)
					backupInfo <- 0
				}
			}
		}
	}()
	/*
		异常退出的情况与第一次启动的情况
		读取app_info.json中的timer字段，并更新
		为避免误差设置提前2秒
	*/
	appInfo := info.GetAppInfo()
	liveTime := time.Now().UnixNano()
	// 备份时间为0,则设置从现在开始
	if appInfo.Timer == 0 {
		appInfo.Timer = liveTime
		_ = info.SetAppInfo(appInfo)
	}
	if appInfo.Timer + int64(2 * time.Second) < liveTime && info.RunTimeNum == 0 {
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
		appBackupTime := appInfo.Timer
		liveTime := time.Now().UnixNano()
		SaveTime := time.NewTimer(time.Duration(appBackupTime - liveTime))
		select {
		case <-SaveTime.C:
			callChan <- 0
			// 更新备份时间
			<-backupInfo
			liveTime = time.Now().UnixNano()
			appInfo.Timer = liveTime + int64(time.Duration(conf.SaveTime) * time.Hour)
			_ = info.SetAppInfo(appInfo)
			info.RunTimeNum++
		}
	}
}
