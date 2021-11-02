package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/abingzo/bups/common/config"
	"github.com/abingzo/bups/common/logger"
	"github.com/abingzo/bups/common/path"
	"github.com/abingzo/bups/common/plugin"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	Name           = "upload"
	DownloadCached = path.PathBackUpCache + "/download"
	Type           = plugin.BCallBack
)

var Support = []int{
	plugin.SupportLogger,
	plugin.SupportConfigRead,
	plugin.SupportArgs,
	plugin.SupportArgs,
}

var StdOut io.Writer = os.Stdout

func New() plugin.Plugin {
	return &Upload{
		Name:       Name,
		Type:       Type,
		Support:    Support,
		stdout:     StdOut,
		cosElement: nil,
	}
}

func InitCosElement(u *Upload) {
	// 初始化实例
	u.cosElement = &CosElement{}
	cfg := config.Read(u.confReader)
	cfg.SetPluginName(u.Name)
	cfg.SetPluginScope("cos")
	// 设置属性
	u.cosElement.sId = cfg.PluginGetData("sId").(string)
	u.cosElement.sKey = cfg.PluginGetData("sKey").(string)
	u.cosElement.bucketUrl = cfg.PluginGetData("bucketUrl").(string)
	u.cosElement.serviceUrl = cfg.PluginGetData("serviceUrl").(string)
	// 连接服务端
	bu, _ := url.Parse(u.cosElement.bucketUrl)
	bsu, _ := url.Parse(u.cosElement.serviceUrl)
	bucket := cos.BaseURL{
		BucketURL:  bu,
		ServiceURL: bsu,
	}
	client := cos.NewClient(&bucket, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  u.cosElement.sId,
			SecretKey: u.cosElement.sKey,
		},
	})
	u.cosElement.client = client
}

/*
	配置文件选项:plugin.upload.cos
	基础的上传至Cos的接口，提供上传，下载，检索
*/

type CosElement struct {
	client     *cos.Client
	sId        string
	sKey       string
	bucketUrl  string
	serviceUrl string
}

func (c *CosElement) Push(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	tmpContainer := strings.Split(path, "/")
	fileName := tmpContainer[len(tmpContainer)]
	_, err = c.client.Object.Put(context.Background(), fileName, file, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *CosElement) Download(fileName string) ([]byte, error) {
	res, err := c.client.Object.Get(context.Background(), fileName, nil)
	if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (c *CosElement) Delete(fileName string) error {
	_, err := c.client.Object.Delete(context.Background(), fileName)
	if err != nil {
		return err
	}
	return nil
}

func (c *CosElement) Search() {}

type Upload struct {
	plugin.Plugin
	Name       string
	Type       plugin.Type
	confReader io.Reader
	confWriter io.Writer
	Support    []int
	stdout     io.Writer
	loggerOut  logger.Logger
	cosElement *CosElement
}

func (u *Upload) SetStdout(out io.Writer) {
	u.stdout = out
}

func (u *Upload) SetLogOut(out logger.Logger) {
	u.loggerOut = out
}

func (u *Upload) GetName() string {
	return u.Name
}

func (u *Upload) GetType() plugin.Type {
	return u.Type
}

func (u *Upload) GetSupport() []int {
	return u.Support
}

func (u *Upload) ConfRead(reader io.Reader) {
	u.confReader = reader
}

func (u *Upload) ConfWrite(writer io.Writer) {
	u.confWriter = writer
}

// Start 启动函数
func (u *Upload) Start(args []string) {
	// 初始化实例
	if u.cosElement == nil {
		InitCosElement(u)
	}
	if args == nil || len(args) == 0 {
		err := u.cosElement.Push(path.PathBackUpCache + "/backup/backup.zip")
		if err != nil {
			panic(err)
		}
		return
	} else {
		os.Args = args
	}
	downloadFileName := flag.String("download", "", "需要下载的文件名")
	searchFileName := flag.String("search", "", "需要搜索的文件名")
	if *downloadFileName != "" {
		bytes, err := u.cosElement.Download(*downloadFileName)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(DownloadCached+"/"+*downloadFileName, bytes, 0755)
		if err != nil {
			panic(err)
		}
		// 打印消息
		_, _ = u.stdout.Write([]byte(fmt.Sprintf("%s 下载成功\n", *downloadFileName)))
		u.loggerOut.Info(fmt.Sprintf("%s 下载成功\n", *downloadFileName))
	} else if *searchFileName != "" {

	}
}
