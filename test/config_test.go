package test

import (
	"errors"
	"github.com/abingzo/bups/common/config"
	"github.com/abingzo/bups/plugins/upload"
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	file, err := os.Open("./test_read.toml")
	if err != nil {
		panic(err)
	}
	cfg := config.Read(file)
	_ = cfg.Plugin["backup"]["file_path"]["root"].(string)
}

func TestWriteConfig(t *testing.T) {
	file, err := os.Open("./test_read.toml")
	if err != nil {
		panic(err)
	}
	cfg := config.Read(file)
	writeFile, err := os.OpenFile("./test_write.toml", os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	if err := config.Write(writeFile, cfg); err != nil {
		panic(err)
	}
}

func TestConfig(t *testing.T) {
	// 设置环境测试配置文件的读取
	envMap := map[string]string{
		"COS_SID":"my sid",
		"COS_SKEY":"my sKey",
		"COS_BUCKET_URL":"my bucket URL",
		"COS_SERVICE_URL":"my service url",
	}
	cmpMap := map[string]interface{}{
		"sId":envMap["COS_SID"],
		"sKey":envMap["COS_SKEY"],
		"bucketUrl":envMap["COS_BUCKET_URL"],
		"serviceUrl":envMap["COS_SERVICE_URL"],
	}
	rawEnvMap := make(map[string]string,4)
	for k := range envMap {
		rawEnvMap[k] = os.Getenv(rawEnvMap[k])
	}
	// 设置回原来的环境值，要不然在GitHub Action中提前设置的环境变量会被覆盖
	defer func() {
		for k,v := range rawEnvMap {
			if v != "" {
				err := os.Setenv(k, v)
				if err != nil {
					t.Fatal(err)
				}
			}
		}
	}()
	for k,v := range envMap{
		err := os.Setenv(k, v)
		if err != nil {
			t.Fatal(err)
		}
	}
	file, err := os.Open("./config.toml")
	if err != nil {
		t.Fatal(err)
	}
	cfg := config.Read(file)
	cfg.SetPluginName(upload.Name)
	cfg.SetPluginScope("cos")
	cfg.RangePluginData(func(k string, v interface{}) {
		vStr := v.(string)
		if mapV,ok := cmpMap[k] ; ok {
			if mapV != vStr {
				t.Fatal(errors.New("set env read failed"))
			}
		} else {
			t.Fatal(errors.New("set env read failed"))
		}
	})
}