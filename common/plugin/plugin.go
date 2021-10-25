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

// 一些重要的，插件必需要实现的函数类型
type Start func()
type GetName func() string
type GetType func() Type
type GetSupport func() []int
type ConfRead func(io.Reader)
type ConfWrite func(io.Writer)

type Context struct {
	lock      sync.Mutex
	start     []p.Symbol
	init      []p.Symbol
	bCallBack []p.Symbol
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
	sb, err := pg.Lookup("GetType")
	if err != nil {
		panic(err)
	}
	// 启动函数
	sFn, err := pg.Lookup("Start")
	if err != nil {
		panic(err)
	}
	// 获取名字和支持
	pluginName, err := pg.Lookup("GetName")
	if err != nil {
		panic(err)
	}
	pluginSupport, err := pg.Lookup("GetSupport")
	if err != nil {
		panic(err)
	}
	// 设置标准输出的函数
	setStdout, err := pg.Lookup("SetStdout")
	if err != nil {
		panic(err)
	}
	// 设置配置文件的读取，写入函数
	confRead, err := pg.Lookup("ConfRead")
	if err != nil {
		panic(err)
	}
	confWrite, err := pg.Lookup("ConfWrite")
	if err != nil {
		panic(err)
	}
	c.support[pluginName.(func() string)()] = pluginSupport.(func() []int)()
	// 实现对应的支持
	for _, v := range c.support[pluginName.(func() string)()] {
		switch v {
		case SupportArgs:
			continue
		case SupportLogger:
			setStdout.(func(io.Writer))(c.StdOut)
		case SupportConfigRead:
			confRead.(func(io.Reader))(c.Conf)
		case SupportConfigWrite:
			confWrite.(func(io.Writer))(c.Conf)
		default:
			panic("not support type")
		}
	}
	// 获取类型
	switch sb.(func() Type)() {
	case Init:
		c.init = append(c.init, sFn)
	case BStart:
		c.start = append(c.start, sFn)
	case BCallBack:
		c.bCallBack = append(c.bCallBack, sFn)
	default:
		panic("not support plugin type")
	}
}

func (c *Context) SetState(s Type) {
	c.lock.Lock()
	defer c.lock.Unlock()
	dst := make([]p.Symbol,0,10)
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
		v.(func())()
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
