package daemon

/*
	封装了一些应用通过守护进程启动的方法
*/

import (
	this "github.com/mengzushan/bups/common/error"
	"github.com/mengzushan/bups/common/logger"
	"github.com/mengzushan/bups/utils"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)


var unixPidPath = "/var/run"
var appPidFileName = "bups.pid"
var appDeaMonFileName = "bupsd.pid"

// 该函数返回自定义错误
func Start() this.Error{
	log, logErr := logger.Std(nil)
	defer log.Close()
	defer utils.ReCoverErrorAndPrint()
	if logErr != this.Nil {
		panic(logErr)
	}
	// 启动之前防止有两个守护进程和任务进程在运行
	// TODO:没有考虑强行杀死的情况
	deaBool,taskBoll := false,false
	_,err := os.Open(unixPidPath + "/" + appDeaMonFileName)
	if err != nil {
		deaBool = !deaBool
	}
	_, err = os.Open(unixPidPath + "/" + appPidFileName)
	if err != nil {
		taskBoll = !taskBoll
	}

	if deaBool && taskBoll {
		// 生成守护进程和子进程的pid文件
		deaMonPid := os.Getpid()
		// 创建pid文件
		err = ioutil.WriteFile(unixPidPath+"/"+appDeaMonFileName, []byte(strconv.Itoa(deaMonPid)), 0066)
		if err != nil {
			log.StdErrorLog(err.Error())
		}
		// exec出一个新的子进程随后守护进程退出
		taskPid,_ := syscall.ForkExec(os.Args[0],nil,&syscall.ProcAttr{
			Env:   append(os.Environ(),[]string{"DAEMON=true"}...),
			Files: []uintptr{0,1,2},
			Sys:   &syscall.SysProcAttr{
				Setsid: true,
			},
		})
		err = ioutil.WriteFile(unixPidPath+ "/" + appPidFileName,[]byte(strconv.Itoa(taskPid)),0066)
		if err != nil {
			log.StdErrorLog(err.Error())
		}
		return this.Nil
	} else {
		return this.SetError("不能开启多个守护进程或者任务进程")
	}
}

func ReStart() this.Error {
	err := Stop()
	if err != this.Nil {
		return err
	}
	err = Start()
	if err != this.Nil {
		return err
	}
	return this.Nil
}

func Stop() this.Error {
	// 删除守护进程和任务进程的pid文件并正常退出任务进程
	taskPidFile,err := os.Open(unixPidPath + "/" + appPidFileName)
	if err != nil {
		return this.SetError(err)
	}
	_, err = os.Open(unixPidPath + "/" + appDeaMonFileName)
	if err != nil {
		return this.SetError(err)
	}
	// 读取任务进程的pid并发送正常退出的posix信号
	taskPidNative,_ := ioutil.ReadAll(taskPidFile)
	taskPid,_ := strconv.Atoi(string(taskPidNative))
	err = syscall.Kill(taskPid, syscall.SIGQUIT)
	if err != nil {
		return this.SetError(err)
	}
	// 删除
	err = os.Remove(unixPidPath + "/" + appPidFileName)
	if err != nil {
		return this.SetError(err)
	}
	err = os.Remove(unixPidPath + "/" + appDeaMonFileName)
	if err != nil {
		return this.SetError(err)
	}
	return this.Nil
}
