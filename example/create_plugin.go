package example

import (
	"github.com/abingzo/bups/common/plugin"
	"github.com/zbh255/bilog"
	"time"
)

// Simple 实现plugin.Plugin接口
type Simple struct {
	name string
	sup []uint32
	_type plugin.Type
	stdLog bilog.Logger
}

// New plugin包定义的必须实现的函数，用于注册组件使用
func New() plugin.Plugin {
	return &Simple{
		name:   "simple",
		sup:    []uint32{plugin.SUPPORT_STDLOG},
		_type:  plugin.Init,
		stdLog: nil,
	}
}

// Start 启动时调用的方法
func (s *Simple) Start(args []string) {
	go func() {
		time.Sleep(time.Second)
		s.stdLog.Info("my is " + s.name)
	}()
}

// Caller 接收可接受的Posix信号时调用的方法
func (s *Simple) Caller(single plugin.Single) {
	return
}

// GetName 返回组件的名字
func (s *Simple) GetName() string {
	return s.name
}

// GetType 返回组件的类型
func (s *Simple) GetType() plugin.Type {
	return s._type
}

// GetSupport 返回组件的需要的资源类型
func (s *Simple) GetSupport() []uint32 {
	return s.sup
}

// SetSource 设置从Context准备好的资源
// context会根据Support Slice声明的类型设置，没有声明的为nil
func (s *Simple) SetSource(source *plugin.Source) {
	s.stdLog = source.StdLog
}
