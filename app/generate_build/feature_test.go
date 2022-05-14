package main

import (
	"os"
	"testing"
)

func TestFeatures(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal(err)
		}
		_ = os.Remove("./generate_file")
	}()
	Generate("../../")
	t.Log(plugin_register)
	GenerateFile("./generate_file")
}
