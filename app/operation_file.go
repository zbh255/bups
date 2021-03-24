package app

import (
	"encoding/json"
	"github.com/alexmullins/zip"
	"github.com/mengzushan/bups/common/logger"
	"github.com/mengzushan/bups/utils"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type InfoJson struct {
	CreateTime  int64  `json:"createTime"`
	CreateSize  int64  `json:"createSize"`
	UseTime     int64  `json:"useTime"`
	UseSize     int64  `json:"useSize"`
	WebPassword string `json:"webPassword"`
	DbPassword  string `json:"dbPassword"`
}

type ConfigJson struct {
	Rsa        string `json:"rsa"`
	Aes        string `json:"aes"`
	Format     string `json:"format"`
	WebPath    string `json:"webPath"`
	StaticPath string `json:"staticPath"`
	LogPath    string `json:"logPath"`
	DbName     string `json:"dbName"`
}

func BackUpForFile() {
	// 读取配置文件
	conf := utils.GetConfig()
	// 创建压缩包内的Json配置文件
	var enON = "on"
	if conf.Encryption.Switch == "off" {
		enON = "off"
	}
	jsons := ConfigJson{
		Rsa:        enON,
		Aes:        enON,
		Format:     "zip",
		WebPath:    conf.Local.Web,
		StaticPath: conf.Local.Static,
		LogPath:    conf.Local.Log,
		DbName:     conf.Database.DbName,
	}
	jsonf, _ := json.Marshal(&jsons)
	// 创建文件写入json
	pathHead, _ := os.Getwd()
	file, _ := os.Create(filepath.FromSlash(pathHead + "/cache/backup/config.json"))
	_, err := file.Write(jsonf)
	defer file.Close()
	if err != nil {
		log := logger.Std()
		defer log.Close()
		log.StdInfoLog("Json配置文件写入失败")
	}
	if conf.Local.Web != "" {
		CreateZip(conf.Local.Web, "web.zip")
	}
	if conf.Local.Static != "" {
		CreateZip(conf.Local.Static, "static.zip")
	}
	if conf.Local.Log != "" {
		CreateZip(conf.Local.Log, "log.zip")
	}
}

func CreateZip(srcPath string, createName string) {
	// 创建待写入的压缩文件
	file, err := os.Create(filepath.FromSlash("./cache/backup/") + createName)
	if err != nil {
		log := logger.Std()
		log.StdErrorLog("文件创建失败" + filepath.FromSlash("./cache/backup/"+createName))
		panic(err)
	}
	// 创建压缩包流
	zipFile := zip.NewWriter(file)
	// 遍历目录
	err = filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 检查遍历的路径是否是源路径，避免出现连个相同的文件夹
		// 如果是则进行下一次遍历
		if path == srcPath {
			return nil
		}
		// 获取文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, srcPath+filepath.FromSlash("/"))
		// 判断文件是不是文件夹
		if info.IsDir() {
			header.Name += filepath.FromSlash("/")
		} else {
			// 设置zip文件的压缩算法
			header.Method = zip.Deflate
		}
		// 创建压缩包头部信息
		w, _ := zipFile.CreateHeader(header)
		// 不是文件夹是将文件copy到流中
		if !info.IsDir() {
			newFile, _ := os.Open(path)
			defer newFile.Close()
			_, _ = io.Copy(w, newFile)
		}
		return nil
	})
	if err != nil {
		log := logger.Std()
		log.StdInfoLog("备份的目录不存在: " + srcPath)
		panic(err)
	}
}
