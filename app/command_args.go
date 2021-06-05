package app

/*
	处理命令行参数
*/

import (
	"flag"
	"fmt"
	"github.com/mengzushan/bups/common/daemon"
	"github.com/mengzushan/bups/common/encry"
	this "github.com/mengzushan/bups/common/error"
	"github.com/mengzushan/bups/common/info"
	"github.com/mengzushan/bups/common/logger"
	"github.com/mengzushan/bups/utils"
	"github.com/mengzushan/bups/web"
	"os"
	"time"
)

// 注册函数调用
var FuncCall func()

func RegisterCall(call func()) {
	FuncCall = call
}

func DeCommandArgs() this.Error {
	// ce是创建参数
	ce := flag.String("c", "","创建选项,value为rsa_key或者info.json")
	// st是启动参数
	st := flag.String("s","","start(启动) restart(重启) stop(停止)")
	// op是操作参数
	op := flag.String("op","","web_config(开启网页配置模式)")
	flag.Parse()
	fmt.Println(os.Args)
	if *ce != "" {
		return create(*ce)
	}
	if *st != "" {
		return start(*st)
	}
	if *op != "" {
		return option(*op)
	}
	return this.Nil
}

func create(options string) this.Error {
	// 初始化日志
	log, err := logger.Std(nil)
	defer log.Close()
	defer utils.ReCoverErrorAndPrint()
	if err != this.Nil {
		panic(err)
	}
	switch options {
	// 创建公钥私钥文件
	case "rsa_key":
		var po encry.PemOptions = &encry.Pem{}
		pathHead,_ := os.Getwd()
		pubf,err := os.Create(pathHead + "/cache/rsa/public.pem")
		prif,err := os.Create(pathHead + "/cache/rsa/private.pem")
		err = po.CreateRsaPubKeyAndPriKey(pubf, prif)
		if err != this.Nil {
			log.StdErrorLog(err.Error())
			return this.SetError(err)
		}
		bl := po.MatchPubKeyAndPriKey(pathHead + "/cache/rsa/public.pem",pathHead + "/cache/rsa/private.pem")
		if !bl {
			err := this.SetError("创建的密钥不正确")
			log.StdErrorLog(err.Error())
			return err
		}
	// 创建应用信息文件
	case "info.json":
		jsons := info.AppInfo{
			Timer:      0,
			BuildTime:  time.Now().UnixNano(),
			AppVersion: info.Version,
		}
		err = info.SetAppInfo(&jsons)
		if err != this.Nil {
			log.StdErrorLog(err.Error())
		} else {
			log.StdInfoLog("app_info.json创建成功")
		}
		return err
	}
	return this.Nil
}

// 守护进程启动方式
func start(value string) this.Error {
	switch value {
	case "start":
		err := daemon.Start()
		if err != this.Nil {
			return err
		}
		break
	case "restart":
		err := daemon.ReStart()
		if err != this.Nil {
			return err
		}
		break
	case "stop":
		err := daemon.Stop()
		if err != this.Nil {
			return err
		}
		break
	}
	return this.Nil
}

// 操作
func option(op string) this.Error {
	switch op {
	case "web_config":
		web.Run()
		break
	}
	return this.Nil
}