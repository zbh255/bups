// go generate build
package app

import (
	"github.com/abingzo/bups/iocc"
	"github.com/abingzo/bups/plugins/backup"
	"github.com/abingzo/bups/plugins/daemon"
	"github.com/abingzo/bups/plugins/encrypt"
	"github.com/abingzo/bups/plugins/recovery"
	"github.com/abingzo/bups/plugins/upload"
	"github.com/abingzo/bups/plugins/web_config"
	
)

func PluginRegister() {
	iocc.RegisterPlugin(backup.New)
	iocc.RegisterPlugin(daemon.New)
	iocc.RegisterPlugin(encrypt.New)
	iocc.RegisterPlugin(recovery.New)
	iocc.RegisterPlugin(upload.New)
	iocc.RegisterPlugin(web_config.New)
	
}
