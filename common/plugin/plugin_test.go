package plugin

import (
	"fmt"
	"github.com/abingzo/bups/common/logger"
	"io"
	"os"
	"testing"
)

type TestPlugin struct {
	typ  Type
	log  logger.Logger
	name string
	sup  []int
}

func (t *TestPlugin) Start(args []string) {
	return
}

func (t *TestPlugin) Caller(single Single) {
	t.log.Info(t.name + ".Caller")
}

func (t *TestPlugin) GetName() string {
	return t.name
}

func (t *TestPlugin) GetType() Type {
	return t.typ
}

func (t *TestPlugin) GetSupport() []int {
	return t.sup
}

func (t *TestPlugin) SetStdout(writer io.Writer) {
	panic("implement me")
}

func (t *TestPlugin) SetLogOut(writer logger.Logger) {
	t.log = writer
}

func (t *TestPlugin) ConfRead(reader io.Reader) {
	panic("implement me")
}

func (t *TestPlugin) ConfWrite(writer io.Writer) {
	panic("implement me")
}

func TestPluginContext(t *testing.T) {
	ctx := NewContext()
	file, err := os.OpenFile("./bups.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	ctx.StdOut = os.Stdout
	ctx.LogOut = func() *os.File {
		return file
	}
	pg1 := Plugin(&TestPlugin{})
	log := logger.New(ctx.LogOut(), logger.ERROR)
	pg1.SetLogOut(log)
	ctx.init = append(ctx.init, pg1)
	// 多个插件
	pg2 := Plugin(&TestPlugin{})
	log2 := logger.New(ctx.LogOut(), logger.ERROR)
	pg2.SetLogOut(log2)
	ctx.handle = append(ctx.handle, pg2)
	// 3
	pg3 := Plugin(&TestPlugin{})
	log3 := logger.New(ctx.LogOut(), logger.ERROR)
	pg3.SetLogOut(log3)
	ctx.collect = append(ctx.collect, pg3)
	// 4
	pg4 := Plugin(&TestPlugin{})
	log4 := logger.New(ctx.LogOut(), logger.ERROR)
	pg4.SetLogOut(log4)
	ctx.bCallBack = append(ctx.bCallBack, pg4)
	ctx.RangeAllPlugin(func(k int, v Plugin) {
		v.Caller(Exit)
	})
}

func TestPluginContextMultiLogger(t *testing.T) {
	// 插件管理
	ctx := NewContext()
	file, err := os.OpenFile("./bups.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	ctx.LogOut = func() *os.File {
		return file
	}
	backup := &TestPlugin{
		typ:  BCollect,
		log:  nil,
		name: "backup",
		sup:  []int{SupportLogger},
	}
	daemon := &TestPlugin{
		typ:  Init,
		log:  nil,
		name: "daemon",
		sup:  []int{SupportLogger, SupportArgs},
	}
	ctx.RegisterRaw(backup)
	ctx.RegisterRaw(daemon)
	// Caller
	ctx.RangeAllPlugin(func(k int, v Plugin) {
		v.Caller(Exit)
	})
}
