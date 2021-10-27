// Package plugin 本包统一插件的机制
package plugin

import (
	"io"
	p "plugin"
	"sync"
)

type Type int

// BStart 单次备份开始调用的插件
const BStart Type = 1

// BCallBack 备份完毕调用的插件
const BCallBack Type = 2

// Init 初始化时调用的插件
const Init Type = 3

// 插件需要的支持
const (
	// SupportArgs 命令行参数支持
	SupportArgs int = iota
	// SupportLogger 输出到内置日志的支持
	SupportLogger
	// SupportConfigRead 配置文件读取的支持
	SupportConfigRead
	// SupportConfigWrite 配置文件写入的支持
	SupportConfigWrite
)

type plugins []Plugin

// Plugin 插件的插入要实现的接口
type Plugin interface {
	Start(args []string)
	GetName() string
	GetType() Type
	GetSupport() []int
	SetStdout(writer io.Writer)
	ConfRead(reader io.Reader)
	ConfWrite(writer io.Writer)
}

// New 一些重要的，插件必需要实现的函数类型
type New func() Plugin

type Context struct {
	lock      sync.Mutex
	start     plugins
	init      plugins
	bCallBack plugins
	support   map[string][]int
	Conf      io.ReadWriteCloser
	StdOut    io.Writer
	// 状态的流转，每流入一个状态时则调用对应的插件启动函数
	state 	  Type
}

func (c *Context) Register(s string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	pg, err := p.Open(s)
	if err != nil {
		panic(err)
	}
	// 调佣
	interFace, err := pg.Lookup("New")
	if err != nil {
		panic(err)
	}
	regPlugin := interFace.(func() Plugin)()
	c.support[regPlugin.GetName()] = regPlugin.GetSupport()
	// 实现对应的支持
	for _, v := range c.support[regPlugin.GetName()] {
		switch v {
		case SupportArgs:
			continue
		case SupportLogger:
			regPlugin.SetStdout(c.StdOut)
		case SupportConfigRead:
			regPlugin.ConfRead(c.Conf)
		case SupportConfigWrite:
			regPlugin.ConfWrite(c.Conf)
		default:
			panic("not support type")
		}
	}
	// 获取类型
	switch regPlugin.GetType() {
	case Init:
		c.init = append(c.init, regPlugin)
	case BStart:
		c.start = append(c.start, regPlugin)
	case BCallBack:
		c.bCallBack = append(c.bCallBack, regPlugin)
	default:
		panic("not support plugin type")
	}
}

func (c *Context) SetState(s Type) {
	c.lock.Lock()
	defer c.lock.Unlock()
	dst := make(plugins,0,10)
	switch s {
	case BStart:
		dst = append(c.start)
	case Init:
		dst = append(c.init)
	case BCallBack:
		dst = append(c.bCallBack)
	default:
		panic("not support state type")
	}
	// call
	for _,v := range dst {
		v.Start(nil)
	}
	c.state = s
}

func (c *Context) GetState() Type {
	return c.state
}

func NewContext() *Context {
	return &Context{
		support: make(map[string][]int),
	}
}
