package test

import (
	"github.com/abingzo/bups/common/plugin"
	"io"
	"os"
	"testing"
)

type Metadata struct {
	b []byte
	c bool
	io.ReadWriteCloser
}

func (m *Metadata) Read(p []byte) (int, error) {
	ptr := 0
	for ptr < len(p) && ptr < len(m.b) {
		p[ptr] = m.b[ptr]
		ptr++
	}
	if ptr == len(m.b) {
		return ptr,io.EOF
	}
	return ptr,nil
}

func (m *Metadata) Write(p []byte) (int,error) {
	m.b = append(m.b,p...)
	return len(m.b),nil
}

func (m *Metadata) Close() error {
	return nil
}

func TestPluginLoad(t *testing.T) {
	ctx := plugin.NewContext()
	ctx.StdOut = os.Stdout
	ctx.Conf = &Metadata{}
	_, _ = ctx.Conf.Write([]byte("这就是配置文件"))
	path,_ := os.Getwd()
	ctx.Register(path + "/plugins/upload")
	ctx.SetState(plugin.BCallBack)
}
