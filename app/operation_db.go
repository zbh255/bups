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
	// 创建sql文件的压缩包
	conf := utils.GetConfig()
	file, err := os.Open(filepath.FromSlash("./cache/backup/" + conf.Database.DbName + ".sql"))
	defer func() {
		// 关闭并删除源sql文件
		_ = file.Close()
		_ = os.RemoveAll("./cache/backup/" + conf.Database.DbName + ".sql")
	}()
	if err != nil {
		return err
	}
	// 创建zip文件
	zipFile, _ := os.Create(filepath.FromSlash("./cache/backup/" + conf.Database.DbName + ".zip"))
	// 创建zip写入器
	archive := zip.NewWriter(zipFile)
	// 写入文件头信息
	info,err := os.Stat(filepath.FromSlash("./cache/backup/" + conf.Database.DbName + ".sql"))
	header, _ := zip.FileInfoHeader(info)
	//path := filepath.FromSlash("./cache/" + conf.Database.DbName + ".sql")
	header.Name = "/"
	header.Method = zip.Deflate
	writer,err := archive.CreateHeader(header)
	_, _ = io.Copy(writer, zipFile)
	return nil
}