package main

import (
	"flag"
	"fmt"
	"github.com/abingzo/bups/common/logger"
	"github.com/abingzo/bups/common/path"
	"github.com/abingzo/bups/common/plugin"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

/*
	将项目中的原daemon模块移植为插件的形式
*/

const (
	Name    = "daemon"
	PidFile = path.PathBackUpCache + "/" + Name + "/" + "bups.pid"
	Type    = plugin.Init
)

var support = []int{plugin.SupportArgs, plugin.SupportNativeStdout}

func New() plugin.Plugin {
	return &Daemon{
		sup: support,
	}
}

type Daemon struct {
	stdOut io.Writer
	sup    []int
	plugin.Plugin
}

func (d *Daemon) Start(args []string) {
	if args == nil {
		return
	}
	os.Args = args
	// st是启动参数
	st := flag.String("s", "", "start(启动) restart(重启) stop(停止)")
	flag.Parse()
	switch *st {
	case "start":
		start(d.stdOut)
	case "stop":
		stop(d.stdOut)
	case "restart":
		restart(d.stdOut)
	default:
		_, _ = os.Stdout.Write([]byte(fmt.Sprintf("不支持的参数:%v\n", args)))
	}
}

func (d *Daemon) GetName() string {
	return Name
}

func (d *Daemon) GetType() plugin.Type {
	return Type
}

func (d *Daemon) GetSupport() []int {
	return d.sup
}

func (d *Daemon) SetStdout(writer io.Writer) {
	d.stdOut = writer
}

func (d *Daemon) SetLogOut(logger logger.Logger) {
	return
}

func (d *Daemon) ConfRead(reader io.Reader) {
	return
}

func (d *Daemon) ConfWrite(writer io.Writer) {
	return
}

func pidFileExist() bool {
	info, err := os.Stat(PidFile)
	if err != nil || !info.IsDir() {
		return false
	}
	return true
}

// 守护进程操作相关函数
// 写入进程号到pidFile,异步启动主进程后退出
func start(stdOut io.Writer) {
	// 同时只能打开一个子进程
	if pidFileExist() {
		_, _ = stdOut.Write([]byte("You cannot open two processes at the same time"))
		return
	}
	pidFile, err := os.OpenFile(PidFile, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		_, _ = stdOut.Write([]byte(err.Error()))
		return
	}
	cmd := exec.Command(os.Args[0])
	if err := cmd.Start(); err != nil {
		_, _ = stdOut.Write([]byte(err.Error()))
		return
	}
	// 写入pid
	_, err = pidFile.Write([]byte(strconv.Itoa(cmd.Process.Pid)))
	if err != nil {
		_, _ = stdOut.Write([]byte(fmt.Sprintf("write to pid file failed: %s", err.Error())))
	}
}
func stop(stdOut io.Writer) {
	if !pidFileExist() {
		_, _ = stdOut.Write([]byte(fmt.Sprintf("%s is not found", PidFile)))
		return
	}
	// 发送信号和清理pidFile
	pidFile, err := os.Open(PidFile)
	if err != nil {
		_, _ = stdOut.Write([]byte(err.Error()))
		return
	}
	bytes, err := ioutil.ReadAll(pidFile)
	if err != nil {
		_, _ = stdOut.Write([]byte(err.Error()))
		return
	}
	pid, err := strconv.Atoi(string(bytes))
	if err != nil {
		panic(err)
	}
	err = syscall.Kill(pid, syscall.SIGQUIT)
	if err != nil {
		_, _ = stdOut.Write([]byte(err.Error()))
		return
	}
	// 信号发送成功则清理pidFile
	err = os.Remove(PidFile)
	if err != nil {
		_, _ = stdOut.Write([]byte(err.Error()))
		return
	}
}
func restart(stdOut io.Writer) {
	stop(stdOut)
	start(stdOut)
}
