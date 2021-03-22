package app

import (
	"github.com/mengzushan/bups/utils"
	"os"
	"go.uber.org/zap"
)

func BackUpForDb()  {
	pathHead,_ := os.Getwd()
	path := pathHead + "/cache/backup"
	conf := utils.GetConfig()
}