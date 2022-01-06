package upload

import (
	"github.com/abingzo/bups/common/config"
	"github.com/abingzo/bups/common/logger"
	"github.com/abingzo/bups/common/plugin"
	"log"
	"os"
	"testing"
)

func TestPluginLoadAndStart(t *testing.T) {
	uploadPg := New()
	configFile, err := os.Open("../../config/dev/config.toml")
	if err != nil {
		panic(err)
	}
	log.Writer()
	rawSource := new(plugin.Source)
	rawSource.AccessLog = logger.New(os.Stderr, logger.DEBUG)
	rawSource.ErrorLog = logger.New(os.Stdout, logger.PANIC)
	rawSource.StdLog = logger.New(os.Stdout, logger.PANIC)
	rawSource.RawConfig = configFile
	rawSource.Config = config.Read(configFile)
	uploadPg.SetSource(rawSource)
	// no args
	uploadPg.Start(nil)
}
