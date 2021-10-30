package test

import (
	"github.com/abingzo/bups/common/logger"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	log := logger.New(os.Stdout,"Main")
	log.Info("lll")
	log.Debug("ddd")
	defer func() {
		err := recover()
		if err != nil {
			return
		}
	}()
	log.Error("my is error")
	log.Trace("my is msg")
	log.Panic(error(nil))
}