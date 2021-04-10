package test

import (
	"github.com/mengzushan/bups/app"
	"github.com/mengzushan/bups/common/encry"
	"github.com/mengzushan/bups/utils"
	"testing"
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

func Test_backUpFileEncrypt(t *testing.T) {
	conf := utils.GetConfig()
	err := app.EncryptFile(app.ENCRYPTON, conf)
	if err != nil {
		t.Error("备份文件加密测试失败: ",err.Error())
	} else {
		t.Log("备份文件加密测试失败!")
	}
}