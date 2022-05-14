package test

import (
	"errors"
	"github.com/abingzo/bups/common/plugin"
	"github.com/abingzo/bups/plugins/backup"
	"github.com/abingzo/bups/plugins/encrypt"
	"testing"
)

// 测试真实插件调用
func TestBackupPlugin(t *testing.T) {
	ctx := plugin.NewContext()
	bpp := backup.New()
	if bpp.GetName() != backup.Name {
		t.Fatal(errors.New("backup plugin name not equal"))
	}
	if bpp.GetType() != backup.Type {
		t.Fatal(errors.New("backup plugin type not equal"))
	}
	source := LoadPluginSource()
	ctx.RawSource = source
	ctx.RegisterRaw(bpp)
	ctx.SetState(backup.Type)
	// Single
	bpp.Caller(plugin.Exit)
}

// 测试Backup插件的代码必须发生在Encrypt之前
// 因为Encrypt依赖Backup的一些结果
func TestEncryptPlugin(t *testing.T) {
	ctx := plugin.NewContext()
	ep := encrypt.New()
	if ep.GetName() != encrypt.Name {
		t.Fatal(errors.New("backup plugin name not equal"))
	}
	if ep.GetType() != encrypt.Type {
		t.Fatal(errors.New("backup plugin type not equal"))
	}
	source := LoadPluginSource()
	ctx.RawSource = source
	ctx.RegisterRaw(ep)
	ctx.SetState(encrypt.Type)
	// Single
	ep.Caller(plugin.Exit)
}