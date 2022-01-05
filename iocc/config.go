package iocc

import (
	"github.com/abingzo/bups/common/config"
	"io"
)

var (
	cfg *config.AutoGenerated
)

func RegisterConfig(reader io.Reader) {
	cfg  = config.Read(reader)
}

func GetConfig() *config.AutoGenerated {
	return cfg
}