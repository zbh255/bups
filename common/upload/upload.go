package upload

import (
	"context"
	cf "github.com/mengzushan/bups/common/conf"
	this "github.com/mengzushan/bups/common/error"
	"github.com/mengzushan/bups/utils"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

/*
	基础的上传至Cos的接口，提供上传，下载，检索
*/

var TestOnFilePath string // 测试文件路径
var ConfModel cf.ConfNum  // 测试配置模块

type Upload interface {
	SetNewLink() this.Error
	Push(fileName string, file io.Reader) this.Error
	Download(fileName string) ([]byte, this.Error)
	Delete(fileName string) this.Error
	Search() (map[string]int64, this.Error)
}

type UploadInterface interface {
	TestSet(num cf.ConfNum, testFilePath string)
	Start() Upload
}

type Func func()

func (h Func) TestSet(num cf.ConfNum, testFilePath string) {
	TestOnFilePath = testFilePath
	ConfModel = num
}

func (h Func) Start() Upload {
	cf.TestOnFilePath = TestOnFilePath
	cf.ConfModel = ConfModel
	conf := cf.InitConfig()
	cosRegion := utils.SplitCosBucketUrl(conf.Bucket.BucketURL)
	var Upload Upload = &Conf{
		sId:        conf.Bucket.Secretid,
		sKey:       conf.Bucket.Secretkey,
		buckUrl:    conf.Bucket.BucketURL,
		serviceUrl: cosRegion,
	}
	return Upload
}

type Conf struct {
	client     *cos.Client
	sId        string
	sKey       string
	buckUrl    string
	serviceUrl string
}

// 初始化一个新的连接
func (c *Conf) SetNewLink() this.Error {
	u, _ := url.Parse(c.buckUrl)
	su, _ := url.Parse(c.serviceUrl)
	bucket := cos.BaseURL{
		BucketURL:  u,
		ServiceURL: su,
	}
	client := cos.NewClient(&bucket, &http.Client{Transport: &cos.AuthorizationTransport{
		SecretID:  c.sId,
		SecretKey: c.sKey,
	}})
	if client == nil {
		return this.SetError("Cos client not connected correctly")
	} else {
		c.client = client
		return this.Nil
	}
}

// 提交文件到腾讯云
// 注意: 该方法会关闭io.Reader接口!
func (c *Conf) Push(fileName string, file io.Reader) this.Error {
	_, err := c.client.Object.Put(context.Background(), fileName, file, nil)
	if err != nil {
		return this.SetError(err)
	}
	return this.Nil
}

func (c *Conf) Download(fileName string) ([]byte, this.Error) {
	res, err := c.client.Object.Get(context.Background(), fileName, nil)
	if err != nil {
		return nil, this.SetError(err)
	}
	file, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, this.SetError(err)
	}
	return file, this.Nil
}

func (c *Conf) Delete(fileName string) this.Error {
	_, err := c.client.Object.Delete(context.Background(), fileName)
	if err != nil {
		return this.SetError(err)
	}
	return this.Nil
}

func (c *Conf) Search() (map[string]int64, this.Error) {
	opt := &cos.BucketGetOptions{Prefix: "test"}
	result, _, err := c.client.Bucket.Get(context.Background(), opt)
	if err != nil {
		return nil, this.SetError(err)
	}
	val := make(map[string]int64)
	for _, v := range result.Contents {
		val[v.Key] = v.Size
	}
	return val, this.Nil
}
