package upload

import (
	"github.com/abingzo/bups/common/config"
	"github.com/abingzo/bups/common/plugin"
	"github.com/abingzo/bups/iocc"
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
	rawSource.AccessLog = iocc.GetAccessLog()
	rawSource.ErrorLog = iocc.GetErrorLog()
	rawSource.StdLog = iocc.GetStdLog()
	rawSource.RawConfig = configFile
	rawSource.Config = config.Read(configFile)
	uploadPg.SetSource(rawSource)
	// no args
	uploadPg.Start(nil)
}
