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
	"os/exec"
	"strconv"
	"syscall"
	"time"
)


var unixPidPathEND = "/dir"
var appPidFileName = "bups.pid"
var appDeaMonFileName = "bupsd.pid"

// 启动
func Start() this.Error {
	pathHead, _ := os.Getwd()
	log, logErr := logger.Std(nil)
	defer log.Close()
	defer utils.ReCoverErrorAndPrint()
	if logErr != this.Nil {
		panic(logErr)
	}
	// 格式化守护进程日志
	stdout,err := os.OpenFile(pathHead + "/log/service.log",os.O_WRONLY|os.O_APPEND,0666)
	if err != nil {
		log.StdErrorLog(err.Error())
		return this.SetError(err)
	}
	// 判断子进程是否在还在运行
	file, err := os.Open(pathHead + unixPidPathEND + "/" + appPidFileName)
	// 打不开视为不存在反之视为已存在
	if err != nil {

	} else {
		pid,_ := ioutil.ReadAll(file)
		log.StdErrorLog("子进程已经存在: " + string(pid))
		return this.SetError("The child process already exists")
	}
	// task process
	cmd := exec.Command(os.Args[0])
	cmd.Stderr = stdout
	cmd.Stdout = stdout
	// 异步启动
	err = cmd.Start()
	if err != nil {
		log.StdErrorLog(err.Error())
		return this.SetError(err)
	}
	return this.Nil
}

// write app pid
// print log
func WritePid(pid int,log *logger.Logger) this.Error {
	pathHead, _ := os.Getwd()
	pidFile, err := os.OpenFile(pathHead + unixPidPathEND + "/" + appPidFileName,os.O_WRONLY|os.O_APPEND|os.O_CREATE,0666)
	if err != nil {
		log.StdErrorLog(err.Error())
		return this.SetError(err)
	}
	_, err = pidFile.Write([]byte(strconv.Itoa(pid)))
	if err != nil {
		log.StdErrorLog("写入应用pid失败" + err.Error())
	}
	log.StdInfoLog("写入应用pid成功:pid-"+ strconv.Itoa(pid))
	return this.Nil
}

// 正常退出删除pid文件
func DelPid() this.Error {
	pathHead, _ := os.Getwd()
	err := os.Remove(pathHead + unixPidPathEND + "/" + appPidFileName)
	return this.SetError(err)
}

//// 该函数返回自定义错误
//func Start() this.Error{
//	log, logErr := logger.Std(nil)
//	defer log.Close()
//	defer utils.ReCoverErrorAndPrint()
//	if logErr != this.Nil {
//		panic(logErr)
//	}
//	// 启动之前防止有两个守护进程和任务进程在运行
//	// TODO:没有考虑强行杀死的情况
//	deaBool,taskBoll := false,false
//	_,err := os.Open(unixPidPath + "/" + appDeaMonFileName)
//	if err != nil {
//		deaBool = !deaBool
//	}
//	_, err = os.Open(unixPidPath + "/" + appPidFileName)
//	if err != nil {
//		taskBoll = !taskBoll
//	}
//
//	if deaBool && taskBoll {
//		// 生成守护进程和子进程的pid文件
//		deaMonPid := os.Getpid()
//		// 创建pid文件
//		err = ioutil.WriteFile(unixPidPath+"/"+appDeaMonFileName, []byte(strconv.Itoa(deaMonPid)), 0644)
//		if err != nil {
//			log.StdErrorLog(err.Error())
//		}
//		// exec出一个新的子进程随后守护进程退出
//		taskPid,_ := syscall.ForkExec(os.Args[0],nil,&syscall.ProcAttr{
//			Env:   append(os.Environ(),[]string{"DAEMON=true"}...),
//			Files: []uintptr{0,1,2},
//			Sys:   &syscall.SysProcAttr{
//				Setsid: true,
//			},
//		})
//		err = ioutil.WriteFile(unixPidPath+ "/" + appPidFileName,[]byte(strconv.Itoa(taskPid)),0644)
//		if err != nil {
//			log.StdErrorLog(err.Error())
//		}
//		return this.Nil
//	} else {
//		return this.SetError("不能开启多个守护进程或者任务进程")
//	}
//}

func ReStart() this.Error {
	err := Stop()
	if err != this.Nil {
		return err
	}
	// 等待几秒执行完成
	time.Sleep(time.Second * 3)
	err = Start()
	if err != this.Nil {
		return err
	}
	return this.Nil
}

//func Stop() this.Error {
//	// 删除守护进程和任务进程的pid文件并正常退出任务进程
//	taskPidFile,err := os.Open(unixPidPath + "/" + appPidFileName)
//	if err != nil {
//		return this.SetError(err)
//	}
//	_, err = os.Open(unixPidPath + "/" + appDeaMonFileName)
//	if err != nil {
//		return this.SetError(err)
//	}
//	// 读取任务进程的pid并发送正常退出的posix信号
//	taskPidNative,_ := ioutil.ReadAll(taskPidFile)
//	taskPid,_ := strconv.Atoi(string(taskPidNative))
//	err = syscall.Kill(taskPid, syscall.SIGQUIT)
//	if err != nil {
//		return this.SetError(err)
//	}
//	// 删除
//	err = os.Remove(unixPidPath + "/" + appPidFileName)
//	if err != nil {
//		return this.SetError(err)
//	}
//	err = os.Remove(unixPidPath + "/" + appDeaMonFileName)
//	if err != nil {
//		return this.SetError(err)
//	}
//	return this.Nil
//}

func Stop() this.Error {
	pathHead, _ := os.Getwd()
	pidFile, err := os.Open(pathHead + unixPidPathEND + "/" + appPidFileName)
	if err != nil {
		return this.SetError(err)
	}
	pidB, _ := ioutil.ReadAll(pidFile)
	pid,_ := strconv.Atoi(string(pidB))
	// 向运行中的子进程发送信号
	err = syscall.Kill(pid,syscall.SIGQUIT)
	if err != nil {
		return this.SetError(err)
	}
	return this.Nil
}