package main

import "github.com/abingzo/bups/common/plugin"

const (
	Name = "encrypt"
	Type = plugin.BHandle
)

var support = []int{plugin.SupportLogger, plugin.SupportConfigRead}
