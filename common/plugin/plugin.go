// Package plugin 本包统一插件的机制
package plugin

import (
	"fmt"
	"github.com/abingzo/bups/common/logger"
	"io"
	p "plugin"
	"sync"
)

type Type int

// 描述插件生命周期的常量
const (
	Init      Type = 0 // 初始化时调用的插件
	BCollect  Type = 1 // 需要搜集备份数据时调用的插件
	BHandle   Type = 2 // 需要处理备份的数据时调用的插件
	BCallBack Type = 3 // 处理完备份数据时调用的插件
)

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
	// SupportNativeStdout 共享输出缓冲区的支持
	SupportNativeStdout
)

type plugins []Plugin

// Plugin 插件的插入要实现的接口
type Plugin interface {
	Start(args []string)
	GetName() string
	GetType() Type
	GetSupport() []int
	SetStdout(writer io.Writer)
	SetLogOut(writer logger.Logger)
	ConfRead(reader io.Reader)
	ConfWrite(writer io.Writer)
}

// New 一些重要的，插件必需要实现的函数类型
type New func() Plugin

type Context struct {
	lock sync.Mutex
	// 可以处理参数的插件
	// 该哨兵属性免去了if
	argsPlugin plugins
	collect    plugins
	handle     plugins
	init       plugins
	bCallBack  plugins
	support    map[string][]int
	Conf       io.ReadWriteCloser
	// 共享的日志输出
	LogOut io.Writer
	// 共享的缓冲输出
	StdOut io.Writer
	// 状态的流转，每流入一个状态时则调用对应的插件启动函数
	state Type
}

func (c *Context) Register(s string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	pg, err := p.Open(s)
	if err != nil {
		panic(err)
	}
	// 调用
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
			c.argsPlugin = append(c.argsPlugin, regPlugin)
		case SupportLogger:
			log := logger.New(c.LogOut, fmt.Sprintf("Plugin.%s", regPlugin.GetName()))
			regPlugin.SetLogOut(log)
		case SupportConfigRead:
			regPlugin.ConfRead(c.Conf)
		case SupportConfigWrite:
			regPlugin.ConfWrite(c.Conf)
		case SupportNativeStdout:
			regPlugin.SetStdout(c.StdOut)
		default:
			panic("not support type")
		}
	}
	// 获取类型
	switch regPlugin.GetType() {
	case Init:
		c.init = append(c.init, regPlugin)
	case BCollect:
		c.collect = append(c.collect, regPlugin)
	case BCallBack:
		c.bCallBack = append(c.bCallBack, regPlugin)
	case BHandle:
		c.handle = append(c.handle, regPlugin)
	default:
		panic("not support plugin type")
	}
}

func (c *Context) SetState(s Type) {
	c.lock.Lock()
	defer c.lock.Unlock()
	dst := make(plugins, 0, 10)
	switch s {
	case BCollect:
		dst = append(c.collect)
	case Init:
		dst = append(c.init)
	case BCallBack:
		dst = append(c.bCallBack)
	case BHandle:
		dst = append(c.handle)
	default:
		panic("not support state type")
	}
	// call
	for _, v := range dst {
		v.Start(nil)
	}
	c.state = s
}

// RangeArgsPlugin 遍历支持参数的插件列表
func (c *Context) RangeArgsPlugin(fn func(k int, v Plugin)) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for k, v := range c.argsPlugin {
		fn(k, v)
	}
}

// RangeAllPlugin 遍历所有注册的插件
func (c *Context) RangeAllPlugin(fn func(k int, v Plugin)) {
	c.lock.Lock()
	defer c.lock.Unlock()
	k := 0
	for _, v := range c.init {
		fn(k, v)
		k++
	}
	for _, v := range c.collect {
		fn(k, v)
		k++
	}
	for _, v := range c.handle {
		fn(k, v)
		k++
	}
	for _, v := range c.bCallBack {
		fn(k, v)
		k++
	}
}

func (c *Context) GetState() Type {
	return c.state
}

func NewContext() *Context {
	return &Context{
		support: make(map[string][]int),
	}
}
