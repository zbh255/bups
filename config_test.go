package main

import (
	"github.com/abingzo/bups/common/config"
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	file, err := os.Open("./test/test_read.toml")
	if err != nil {
		panic(err)
	}
	cfg := config.Read(file)
	t.Log(cfg.Plugin["backup"]["file_path"]["log"].(string))
}

func TestWriteConfig(t *testing.T) {
	file, err := os.Open("./test/test_read.toml")
	if err != nil {
		panic(err)
	}
	cfg := config.Read(file)
	writeFile, err := os.OpenFile("./test/test_write.toml", os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	if err := config.Write(writeFile, cfg); err != nil {
		panic(err)
	}
}
