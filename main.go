package main

import (
	"github.com/mengzushan/bups/utils"
	"github.com/mengzushan/bups/web"
	"time"
)

func main() {
	conf := utils.GetConfig()
	// webui配置项为开启则进入
	go func() {
		time.Sleep(time.Second * 1)
		if conf.WebConfig.Switch == "on" {
			web.Run()
		}
	}()
	for {
		time.Sleep(time.Second * 1)
	}
}
