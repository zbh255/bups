package info

import (
	"encoding/json"
	this "github.com/mengzushan/bups/common/error"
	"github.com/mengzushan/bups/utils"
	"os"
	"time"
)

/*
	Timer : 定时器时间
	BuildTime : 更改时间
	AppVersion : app版本号
*/
type AppInfo struct {
	Timer      int64   `json:"timer"`
	BuildTime  int64   `json:"build_time"`
	AppVersion float64 `json:"app_version"`
}

const (
	Path    string  = "/dir/app_info.json"
	Version float64 = 0.10
)

// 程序运行时成功的任务数量
var RunTimeNum int = 0

func GetAppInfo() *AppInfo {
	pwd, err := os.Getwd()
	defer utils.ReCoverErrorAndLog()
	if err != nil {
		panic(err)
	}
	file, err := os.Open(pwd + Path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	jsond := AppInfo{}
	jsons := json.NewDecoder(file)
	err = jsons.Decode(&jsond)
	return &jsond
}

// 设置json文件参数为AppInfo结构体指针
// 返回的错误为自定义错误
func SetAppInfo(ptr *AppInfo) this.Error {
	ptr.AppVersion = Version
	ptr.BuildTime = time.Now().Unix()
	pwd, err := os.Getwd()
	if err != nil {
		return this.SetError(err)
	}
	file, err := os.OpenFile(pwd+Path, os.O_WRONLY, 0)
	if err != nil {
		return this.SetError(err)
	}
	defer file.Close()
	jsons, err := json.Marshal(ptr)
	_, err = file.Write(jsons)
	if err != nil {
		return this.SetError(err)
	}
	return this.Nil
}
