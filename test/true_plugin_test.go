package test

import (
	"bytes"
	"errors"
	"github.com/abingzo/bups/common/plugin"
	"github.com/abingzo/bups/plugins/backup"
	"github.com/abingzo/bups/plugins/encrypt"
	"github.com/abingzo/bups/plugins/upload"
	"github.com/abingzo/bups/plugins/web_config"
	"io/ioutil"
	"net/http"
	"os"
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
	// 重新生成文件，给测试Upload的代码使用
	ep.Start(nil)
}

func TestUploadPlugin(t *testing.T) {
	ctx := plugin.NewContext()
	ud := upload.New()
	if ud.GetName() != upload.Name {
		t.Fatal(errors.New("backup plugin name not equal"))
	}
	if ud.GetType() != upload.Type {
		t.Fatal(errors.New("backup plugin type not equal"))
	}
	source := LoadPluginSource()
	ctx.RawSource = source
	ctx.RegisterRaw(ud)
	ctx.SetState(upload.Type)
	// Single
	ud.Caller(plugin.Exit)
	//// Args Start
	//ud.Start([]string{
	//	"--download",
	//	"2006-01-02 15:04:05.zip",
	//})
}

func TestWebConfigPlugin(t *testing.T) {
	ctx := plugin.NewContext()
	webConfig := web_config.New()
	if webConfig.GetName() != web_config.Name {
		t.Fatal(errors.New("backup plugin name not equal"))
	}
	if webConfig.GetType() != web_config.Type {
		t.Fatal(errors.New("backup plugin type not equal"))
	}
	source := LoadPluginSource()
	file, err := os.OpenFile("./config.toml", os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	source.RawFile = file
	ctx.RawSource = source
	ctx.RegisterRaw(webConfig)
	args := [][]string{
		{"--switch","off","--bind","127.0.0.1:8080"},
		{"--switch","on","--bind","127.0.0.1:8080"},
	}
	// recover 使阻塞goroutine中断
	defer func() {
		err := recover()
		if err != nil {
			t.Fatal(err)
		}
	}()
	// 从web_config中读取数据
	go func() {
		rep,err := http.Get("http://127.0.0.1:8080/config")
		if err != nil {
			t.Fatal(err)
		}
		buf, err := ioutil.ReadAll(rep.Body)
		if err != nil {
			t.Fatal(err)
		}
		rep.Body.Close()
		rep,err = http.Post("http://127.0.0.1:8080/config","",bytes.NewReader(buf))
		if err != nil {
			t.Fatal(err)
		}
		if rep.StatusCode != http.StatusOK {
			t.Fatal(errors.New("http status is error"))
		}
		webConfig.Caller(plugin.Exit)
	}()
	for _,v := range args {
		webConfig.Start(v)
	}
}