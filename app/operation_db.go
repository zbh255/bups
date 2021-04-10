package app

import (
	"archive/zip"
	"fmt"
	"github.com/mengzushan/bups/common/logger"
	"github.com/mengzushan/bups/utils"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func BackUpForDb() {
	// 通过拼接mysqldump命令实现备份.sql文件
	// mysqldump -u root -p 123456 godb > godb.sql
	//pathHead,_ := os.Getwd()
	pathRune := filepath.FromSlash("/")
	conf := utils.GetConfig()
	//CmdCd := fmt.Sprintf("cd %s%scache%sbackup",pathHead,pathRune,pathRune)
	//command := exec.Command("cd", CmdCd)
	args := fmt.Sprintf("-u%s -p%s %s > .%scache%sbackup%s%s.sql", conf.Database.UserName, conf.Database.UserPasswd, conf.Database.DbName, pathRune, pathRune, pathRune,conf.Database.DbName)
	// 创建.sh文件
	cmd := "mysqldump " + args
	bkF,_ := os.Create("bk.sh")
	_, _ = bkF.Write([]byte(cmd))
	_ = bkF.Close()
	command2 := exec.Command("bash","bk.sh")
	//if err := command.Run() ;err != nil {
	//	log := logger.Std()
	//	log.StdErrorLog("Cd到不存在的目录: " + CmdCd)
	//	defer log.Close()
	//	panic(err)
	//}
	// 初始化日志
	log := logger.Std()
	defer log.Close()
	if err := command2.Run(); err != nil {
		log.StdErrorLog("mysqldump导出错误,导出参数: " + args)
		panic(err)
	}
	err := CreateSqlFileZip()
	if err != nil {
		panic(err)
	}
	log.StdWarnLog("mysqldump导出数据成功,导出参数: " + args)
}

func CreateSqlFileZip() error {
	/*
		重写创建压缩包的逻辑
	*/
	// 读取配置文件
	conf := utils.GetConfig()
	// 创建日志器
	log := logger.Std()
	defer log.Close()
	// 创建带缓存的字节流
	path := "./cache/backup/"
	buf, _ := os.Create(path + conf.Database.DbName + ".zip")
	// 创建一个写入器
	write := zip.NewWriter(buf)
	defer write.Close()
	// 通过元数据访问文件
	fs, _ := os.Stat(path + conf.Database.DbName + ".sql")
	f,err := write.Create(fs.Name())
	if err != nil {
		log.StdErrorLog(err.Error())
		panic(err)
	}
	// 打开文件
	file, err := os.Open(path + fs.Name())
	if err != nil {
		log.StdErrorLog(err.Error())
		panic(err)
	}
	_, _ = io.Copy(f, file)
	return nil
}