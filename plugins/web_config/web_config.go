package web_config

import (
	"encoding/json"
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


var (
	// 插件需要的支持
	support = []uint32{
		plugin.SUPPORT_STDLOG,
		plugin.SUPPORT_ARGS,
		plugin.SUPPORT_RAW_FILE,
	}
	// 处理参数
	sw = flag.String("switch", "off", "web_config的开关")
	bind = flag.String("bind", "127.0.0.1:8080", "web_config绑定的ip&port")
)

func New() plugin.Plugin {
	return &WebConfig{}
}

type WebConfig struct {
	stdLog     bilog.Logger
	file *os.File
	server *http.Server
}

func (w *WebConfig) SetSource(source *plugin.Source) {
	w.stdLog = source.StdLog
	w.file = source.RawFile
}

func (w *WebConfig) GetName() string {
	return Name
}

func (w *WebConfig) GetType() plugin.Type {
	return Type
}

func (w *WebConfig) GetSupport() []uint32 {
	return support
}

func (w *WebConfig) Caller(s plugin.Single) {
	switch s {
	case plugin.Exit:
		err := w.server.Close()
		if err != nil {
			w.stdLog.ErrorFromErr(err)
		}
	}
}

// Start 启动函数
func (w *WebConfig) Start(args []string) {
	// args不为nil时代表参数启动
	if args == nil {
		return
	}
	//os.Args = args
	err := flag.CommandLine.Parse(args)
	if err != nil {
		panic(err)
	}
	if *sw == "off" {
		w.stdLog.Info("off")
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/config", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			writer.WriteHeader(http.StatusOK)
			// from offset 0 read
			buf := make([]byte,0,256)
			var fileRN int64
			for {
				var tmpBuf [64]byte
				readN,err := w.file.ReadAt(tmpBuf[:],fileRN)
				if err != nil && err != io.EOF {
					w.stdLog.ErrorFromErr(err)
					return
				}
				fileRN += int64(readN)
				buf = append(buf,tmpBuf[:readN]...)
				// 读取到末尾代表读取完毕
				if readN != len(tmpBuf) {
					break
				}
			}

			n, err := writer.Write(buf)
			if err != nil {
				w.stdLog.ErrorFromString(err.Error())
			}
			if n != len(buf) {
				w.stdLog.ErrorFromString("write bytes is not equal")
			}
		case http.MethodPost:
			configBytes,err := ioutil.ReadAll(request.Body)
			if err != nil {
				w.stdLog.ErrorFromString(err.Error())
				return
			}
			// 从偏移量0写起
			_, err = w.file.WriteAt(configBytes,0)
			if err != nil {
				w.stdLog.ErrorFromString(err.Error())
				return
			}
			writer.WriteHeader(http.StatusOK)
			rep,err := json.Marshal(&struct {
				status int
				date interface{}
			}{
				status: http.StatusOK,
				date: nil,
			})
			if err != nil {
				w.stdLog.ErrorFromString(err.Error())
				return
			}
			_, err = writer.Write(rep)
			if err != nil {
				w.stdLog.ErrorFromString(err.Error())
			}
		}
	})
	server := &http.Server{
		Addr:              *bind,
		Handler:           mux,
	}
	w.server = server
	w.stdLog.ErrorFromString(server.ListenAndServe().Error())
}
