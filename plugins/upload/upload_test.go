package main

import (
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
	uploadPg.SetLogOut(os.Stdout)
	// no args
	uploadPg.Start(nil)
}
