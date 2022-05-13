package web_config

import (
	"flag"
	"github.com/abingzo/bups/common/plugin"
	"github.com/zbh255/bilog"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	Name = "web_config"
	Type = plugin.Init
)

// 插件需要的支持
var support = []uint32{
	plugin.SUPPORT_STDLOG,
	plugin.SUPPORT_ARGS,
	plugin.SUPPORT_RAW_CONFIG,
}

func New() plugin.Plugin {
	return &WebConfig{
		name:    Name,
		typ:     Type,
		support: support,
	}
}

type WebConfig struct {
	stdLog     bilog.Logger
	name       string
	typ        plugin.Type
	support    []uint32
	confReader io.Reader
	confWriter io.Writer
	plugin.Plugin
}

func (w *WebConfig) SetSource(source *plugin.Source) {
	w.stdLog = source.StdLog
	w.confReader = source.RawConfig
	w.confWriter = source.RawConfig
}

func (w *WebConfig) GetName() string {
	return w.name
}

func (w *WebConfig) GetType() plugin.Type {
	return w.typ
}

func (w *WebConfig) GetSupport() []uint32 {
	return w.support
}

func (w *WebConfig) Caller(s plugin.Single) {
	w.stdLog.Info(Name + ".Caller")
}

// Start 启动函数
func (w *WebConfig) Start(args []string) {
	// args不为nil时代表参数启动
	if args == nil {
		return
	}
	os.Args = args
	_ = flag.CommandLine.Parse(args)
	// 处理参数
	sw := flag.String("switch", "off", "web_config的开关")
	bind := flag.String("bind", "127.0.0.1:8080", "web_config绑定的ip&port")
	flag.Parse()
	if *sw == "off" {
		w.stdLog.Info("off")
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/config", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			writer.WriteHeader(http.StatusOK)
			configData, err := ioutil.ReadAll(w.confReader)
			if err != nil {
				w.stdLog.ErrorFromString(err.Error())
				return
			}
			n, err := writer.Write(configData)
			if err != nil {
				w.stdLog.ErrorFromString(err.Error())
			}
			if n != len(configData) {
				w.stdLog.ErrorFromString("write bytes is not equal")
			}
		case http.MethodPost:
			configBytes,err := ioutil.ReadAll(request.Body)
			if err != nil {
				w.stdLog.ErrorFromString(err.Error())
				return
			}
			_, err = w.confWriter.Write(configBytes)
			if err != nil {
				w.stdLog.ErrorFromString(err.Error())
				return
			}
		}
	})
	server := &http.Server{
		Addr:              *bind,
		Handler:           mux,
	}
	w.stdLog.ErrorFromString(server.ListenAndServe().Error())
}
