package plugin

import (
	"github.com/abingzo/bups/common/config"
	"github.com/abingzo/bups/common/logger"
	"io"
)

// Source 插件使用到的资源定义
// 根据需求注册
type Source struct {
	Config *config.AutoGenerated
	// 原生配置文件接口
	RawConfig io.ReadWriteCloser
	AccessLog logger.Logger
	ErrorLog logger.Logger
	StdLog logger.Logger
}

func (s *Source) GetConfigReader() io.Reader {
	return s.RawConfig
}

func (s *Source) GetConfigWriter() io.Writer {
	return s.RawConfig
}

func (s *Source) GetConfigReadWriter() io.ReadWriter {
	return s.RawConfig
}