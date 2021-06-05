package test

import (
	"bytes"
	"fmt"
	"github.com/mengzushan/bups/app"
	"github.com/mengzushan/bups/common/conf"
	"github.com/mengzushan/bups/common/encry"
	"github.com/mengzushan/bups/common/error"
	"github.com/mengzushan/bups/common/logger"
	"github.com/mengzushan/bups/common/upload"
	"github.com/mengzushan/bups/utils"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func Test_encrypt(t *testing.T) {
	var er encry.Crypt = &encry.CryptBlocks{}
	x1, err := er.EncryptToAes([]byte("hello worldordershowmecode------------->"), []byte("1234567890123456"))
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log("Aes加密测试成功!")
	}
	_, err = er.DecryptToAes(x1, []byte("1234567890123456"))
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log("Aes解密测试成功!")
	}
}

// 测试Rsa公钥私钥加解密
func Test_Encrypt_Rsa(t *testing.T) {
	// 在此之前先创建密钥对
	var po encry.PemOptions = &encry.Pem{}
	pathHead,_ := os.Getwd()
	pubf,err := os.Create(pathHead + "/cache/rsa/public.pem")
	prif,err := os.Create(pathHead + "/cache/rsa/private.pem")
	err = po.CreateRsaPubKeyAndPriKey(pubf, prif)
	if err != error.Nil {
		t.Error("创建密钥对失败: ",err.Error())
	} else {
		t.Log("创建密钥对成功: ",err.Error())
	}
	bl := po.MatchPubKeyAndPriKey(pathHead + "/cache/rsa/public.pem",pathHead + "/cache/rsa/private.pem")
	if bl {
		t.Log("密钥匹配成功")
	} else {
		t.Error("密钥匹配失败")
	}
	// 加解密测试
}

func Test_backUpFileEncrypt(t *testing.T) {
	// 测试备份的zip能否被正确加密
	// 创建一些测试文件
	pwd, _ := os.Getwd()
	_ = os.Mkdir(pwd+"/test_files", 0777)
	_, _ = os.Create(pwd + "/test_files.zip")
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	for i := 0; i < 100; i++ {
		file, err := os.Create(pwd + "/test_files" + "/" + strconv.Itoa(i) + ".txt")
		if err != nil {
			panic(err)
		}
		var data = make([]byte, 27)
		var node int
		for j := 65; j <= 91; j++ {
			data[node] = byte(j)
			node++
		}
		_, _ = file.Write(data)
		_ = file.Close()
	}
	err := app.Zip(pwd+"/test_files", "test_files.zip")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("解压缩测试成功!")
	}

}

// 请执行完zip压缩测试在来测试此aes加密模块
func Test_Encrypt_Aes(t *testing.T) {
	var er encry.Crypt = &encry.CryptBlocks{}
	pwd, err := os.Getwd()
	key := "1234567890123456"
	defer func() {
		err := recover()
		if err != nil {
			t.Error("Aes加密与解密压缩文件测试不通过")
			fmt.Println(err)
		}
	}()
	if err != nil {
		panic(err)
	}
	file, err := os.Open(pwd + "/test_files.zip")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	src, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	src, err = er.EncryptToAes(src, []byte(key))
	if err != nil {
		panic(err)
	}
	// 创建文件写入加密后的数据
	file2, _ := os.Create(pwd + "/test_files.aes")
	_, err = file2.Write(src)
	if err != nil {
		panic(err)
	}
	_ = file2.Close()
	// 读取加密后的文件数据
	file2, err = os.Open(pwd + "/test_files.aes")
	if err != nil {
		panic(err)
	}
	defer file2.Close()
	dst, err := ioutil.ReadAll(file2)
	if err != nil {
		panic(err)
	}
	src, err = er.DecryptToAes(dst, []byte(key))
	if err != nil {
		panic(err)
	} else {
		t.Log("Aes加密与解密压缩文件测试通过")
	}
}

// 测试日志系统的错误
func Test_Error_Log(t *testing.T) {
	log, err := logger.Std(nil)
	defer log.Close()
	if err == error.Nil {
		t.Log("正确捕捉错误类型: ", err.Error())
		// 错误为空则打印日志
		log.StdInfoLog("hell")
		log.StdDebugLog("world")
	} else if err == error.FileOpenErr {
		t.Log("正确捕捉错误类型: ", err.Error())
	}
}

// 测试日志系统的写入结果
func Test_log(t *testing.T) {
	pathCwd, _ := os.Getwd()
	fileName := "/app.log"
	t.Logf("测试程序将在目录:%s,创建文件: %s", pathCwd, fileName)
	file, err := os.Create(pathCwd + fileName)
	if err != nil {
		t.Error(err)
	}
	_ = file.Close()
	log, err := logger.Std(pathCwd + fileName)
	if err != error.Nil {
		t.Error(err)
	}
	defer log.Close()
	data := []string{"h", "e", "l", "l", "o", ",", "w", "o", "r", "l", "d"}
	for i := 0; i < len(data); i++ {
		log.StdInfoLog(data[i])
	}
	file, err = os.Open(pathCwd + fileName)
	if err != nil {
		t.Error(err)
	}
	filedata, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error(err)
	}
	if len(filedata) > 11*2 {
		t.Log("追加日志测试成功")
	} else {
		t.Error("追加日志测试失败")
	}
}

// 测试多协程读写日志
func Test_Multiple_Logs(t *testing.T) {
	pathCwd, _ := os.Getwd()
	fileName := "/app.log"
	t.Logf("测试程序{多协程写日志}将在目录:%s,创建文件: %s", pathCwd, fileName)
	file, err := os.Create(pathCwd + fileName)
	if err != nil {
		t.Error(err)
	}
	_ = file.Close()
	log, err := logger.Std(pathCwd + fileName)
	if err != error.Nil {
		t.Error(err)
	}
	defer log.Close()
	data := []string{"h", "e", "l", "l", "o", ",", "w", "o", "r", "l", "d"}
	for i := 0; i < 10000; i++ {
		go func() {
			log.StdInfoLog(data[1])
		}()
	}
}

func Benchmark_Utils_Equal(t *testing.B) {
	for i := 0; i < t.N; i++ {
		utils.Equal("awfeawfawfwaappaewfawfawefawfawefeawawefawfawfawfwafaiwuhgiuawhgiuwhaighawighaiwhgiawhgiawhgiauwhgawhfiuawhfiawhfiawhfiawuhfeiuwahifhwaihfiwahfiwahfiwhaiefuhwaiefhiawuhfiawhfiawhifhwaifhwiahfiwahfiwuhefiwhaifuhwaifbewaifbiwabf", "hgiuawhgiuwhaighawighaiwhgiawhgiawhgiauwhgawhfiuawhfiawhfia")
	}
	//if bl == false {
	//	t.Error("Equal结果测试不正确")
	//} else {
	//	t.Log("Equal结果测试正确")
	//}
}

func Test_Utils_Equal(t *testing.T) {
	bl := utils.Equal("/User/Harder/Mac/Mini/web.zip", "web.zip")
	if bl == false {
		t.Error("Equal结果测试不正确")
	} else {
		t.Log("Equal结果测试正确")
	}
}

func Benchmark_Utils_Equal_Kmp(b *testing.B) {
	aString := make([]byte,1000000)
	for i := 0 ; i < 1000000; i++ {
		rand.Seed(time.Now().UnixNano())
		rn := rand.Intn(91 - 65) + 65
		aString[i] = byte(rn)
	}
	t1 := time.Now().Unix()
	utils.Equal(string(aString),string(aString[900000:933000]))
	t2 := time.Now().Unix()
	println(t2 - t1)
}

func Benchmark_Utils_Equals(t *testing.B) {
	aString := make([]byte, 1000000)
	for i := 0; i < 1000000; i++ {
		aString[i] = 65
	}
	bString := aString[544444:644444]

	for i := 0; i < t.N; i++ {
		utils.Equals(string(aString), string(bString))
	}
}

func Benchmark_Utils_EqualPro(t *testing.B) {
	aString := make([]byte, 1000000)
	for i := 0; i < 1000000; i++ {
		aString[i] = 65
	}
	bString := aString[544444:644444]

	for i := 0; i < t.N; i++ {
		utils.Equal(string(aString), string(bString))
	}
}

// 对文件匹配列表功能函数的测试
func Test_App_Match(t *testing.T) {
	conf.ConfModel = conf.ConfModelDev
	config := conf.InitConfig()
	fl, err := app.MatchPathFile(config)
	if err != error.Nil {
		t.Log("文件匹配列表功能测试结果: " + err.Error())
		t.Log(fl)
	} else {
		t.Log("文件匹配列表功能测试成功")
	}
}

// 对操作Cos资源的接口进行测试
func Test_Cos_Upload(t *testing.T) {
	var u upload.UploadInterface = new(upload.Func)
	pathHead, _ := os.Getwd()
	path := pathHead + "/conf/dev/app.conf.toml"
	u.TestSet(conf.ConfModelDev, path)
	ul := u.Start()
	err := ul.SetNewLink()
	if err != error.Nil {
		t.Error("腾讯Cos连接测试失败" + err.Error())
	} else {
		t.Log("腾讯Cos连接测试成功")
	}
	file, err2 := os.Open("./test_files.zip")
	if err2 != nil {
		t.Error("文件打开失败: ", err.Error())
	}
	cosFileName := "test_files.zip"
	err = ul.Push(cosFileName, file)
	if err != error.Nil {
		t.Error("腾讯Cos提交文件接口测试失败: ", err.Error())
	} else {
		t.Log("腾讯Cos提交文件接口测试成功")
	}
	err = ul.Delete(cosFileName)
	if err != error.Nil {
		t.Error("腾讯Cos删除文件接口测试失败: ", err.Error())
	} else {
		t.Log("腾讯Cos删除文件接口测试成功")
	}
	// 重新打开文件句柄，因为流已经被关闭
	file, err2 = os.Open("./test_files.zip")
	if err2 != nil {
		t.Error(err2)
	}
	// 先提交文件在查找
	t.Log(fmt.Sprintf("Push文件:%s -> Cloud", cosFileName))
	err = ul.Push(cosFileName, file)
	if err != error.Nil {
		t.Error(err.Error())
	}
	mp, err := ul.Search()
	if err != error.Nil {
		t.Error("腾讯Cos查找文件接口测试失败: " + err.Error())
	} else {
		t.Log("腾讯Cos查找文件接口测试成功")
	}
	t.Log(fmt.Sprintf("文件搜索结果: %s", reflect.ValueOf(mp)))
	// 文件下载
	fileData, err := ul.Download(cosFileName)
	if err != error.Nil {
		t.Error("腾讯Cos文件下载接口测试失败: " + err.Error())
	}
	// 子测试，匹配本地文件与下载的文件是否相同
	t.Run("匹配下载文件与本地文件", func(t *testing.T) {
		// 重新打开文件句柄，因为流已经被关闭
		file, _ = os.Open("./test_files.zip")
		defer file.Close()
		fd, _ := ioutil.ReadAll(file)
		if bytes.Equal(fileData, fd) {
			t.Log("本地文件与下载文件匹配，腾讯Cos下载接口测试成功")
		} else {
			t.Error("本地文件与下载文件不匹配，腾讯Cos下载接口测试失败")
		}
	})
}

// 测试配置文件的读取
func Test_Config_Read(t *testing.T) {
	pathHead, _ := os.Getwd()
	path := pathHead + "/conf/dev/app.conf.toml"
	conf.TestOnFilePath = path
	conf.ConfModel = conf.ConfModelDev
	_ = conf.InitConfig()
	defer func() {
		r := recover()
		if r != nil {
			t.Error("配置文件读取测试失败")
		} else {
			t.Log("配置文件读取测试成功")
		}
	}()
}

// utils测试EqualToStrings
func Test_Utils_EqualToStrings(t *testing.T) {
	bl := utils.EqualToStrings([]string{"foo", "bar"}, []string{"bar", "foo"})
	if bl {
		t.Log("utils.EqualToStrings " + "测试成功")
	} else {
		t.Error("utils.EqualToStrings " + "测试失败")
	}
}

// 测试备份文件中的json配置文件
func Test_BackUp_Config_Json(t *testing.T) {
	app.BackUpForFile()
}

// app调用测试
// 无加密文件选项
func Test_App_Timer(t *testing.T) {
	_, err := app.TimerTask(conf.InitConfig())
	if err != error.Nil {
		t.Error(err)
	} else {
		t.Log(err)
	}
}

// app调用测试
// 有加密文件选项
func Test_App_Timer_Encrypt(t *testing.T) {
	config := conf.InitConfig()
	config.Encryption.Switch = "on"
	config.Encryption.Aes = "1234567890123456"
	config.Local.Log = ""
	_, err := app.TimerTask(config)
	if err != error.Nil {
		t.Error(err)
	} else {
		t.Log(err)
	}
}