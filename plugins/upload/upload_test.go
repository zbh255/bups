package main

import (
	"github.com/abingzo/bups/common/logger"
	"log"
	"os"
	"testing"
)

func TestPluginLoadAndStart(t *testing.T) {
	uploadPg := New()
	configFile, err := os.Open("../../conf/dev/config.toml")
	if err != nil {
		panic(err)
	}
	log.Writer()
	uploadPg.ConfRead(configFile)
	uploadPg.ConfWrite(configFile)
	uploadPg.SetStdout(os.Stdout)
	uploadPg.SetLogOut(logger.New(os.Stdout, "Plugin.upload."))
	// no args
	uploadPg.Start(nil)
}
