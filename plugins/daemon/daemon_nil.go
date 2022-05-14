// +build windows

/*
	为了兼容daemon本身，而syscall.Kill Api不支持windows
	所以如果GOOS=windows则构建该文件，下面的代码流程正确。
	但是，并不做任何事情
*/

package daemon

import (
	"github.com/zbh255/bilog"
	"github.com/abingzo/bups/common/path"
	"github.com/abingzo/bups/common/plugin"
)

const (
	Name    = "daemon_not_support"
	PidFile = path.DEFAULT_PATH_BACK_UPCACHE + "/" + Name + "/" + "bups.pid"
	Type    = plugin.Init
)

var support = []uint32{plugin.SUPPORT_ARGS, plugin.SUPPORT_STDLOG}

func New() plugin.Plugin {
	return &Daemon{
		sup: support,
	}
}

type Daemon struct {
	stdLog bilog.Logger
	sup    []uint32
	plugin.Plugin
}

func (d *Daemon) Caller(s plugin.Single) {
	d.stdLog.Info(Name + ":" + "Caller")
}

func (d *Daemon) Start(args []string) {
	d.stdLog.Info("daemon plugin is not support windows")
}

func (d *Daemon) GetName() string {
	return Name
}

func (d *Daemon) GetType() plugin.Type {
	return Type
}

func (d *Daemon) GetSupport() []uint32 {
	return d.sup
}

func (d Daemon) SetSource(source *plugin.Source) {
	d.stdLog = source.StdLog
}
