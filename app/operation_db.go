package app

import (
	"fmt"
	"github.com/mengzushan/bups/common/logger"
	"github.com/mengzushan/bups/utils"
	"os/exec"
)

func BackUpForDb()  {
	// 通过拼接mysqldump命令实现备份.sql文件
	// mysqldump -u root -p 123456 godb > godb.sql
	//pathHead,_ := os.Getwd()
	pathRune := utils.PathRune()
	conf := utils.GetConfig()
	//CmdCd := fmt.Sprintf("cd %s%scache%sbackup",pathHead,pathRune,pathRune)
	//command := exec.Command("cd", CmdCd)
	args := fmt.Sprintf("-u %s -p %s %s > %scache%sbackup%s.sql",conf.Database.UserName,conf.Database.UserPasswd,conf.Database.DbName,pathRune,pathRune,conf.Database.DbName)
	command2 := exec.Command("mysqldump",args)
	//if err := command.Run() ;err != nil {
	//	log := logger.Std()
	//	log.StdErrorLog("Cd到不存在的目录: " + CmdCd)
	//	defer log.Close()
	//	panic(err)
	//}
	// 初始化日志
	log := logger.Std()
	defer log.Close()
	if err := command2.Run() ; err != nil {
		log.StdErrorLog("mysqldump导出错误,导出参数: " + args)
		panic(err)
	}
	log.StdWarnLog("mysqldump导出数据成功,导出参数: " + args)
}