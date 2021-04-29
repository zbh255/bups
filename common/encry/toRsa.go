package encry

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	this "github.com/mengzushan/bups/common/error"
	"github.com/mengzushan/bups/common/logger"
	"github.com/mengzushan/bups/utils"
	"io"
	"os"
)

type PemOptions interface {
	CreateRsaPubKeyAndPriKey(pubf,prif io.Writer) this.Error
	MatchPubKeyAndPriKey(pub, pri string) bool
}

type CryptToRsa interface {
	EncryptToRsa(src []byte,pubKey []byte) ([]byte, this.Error)
	DecryptToRsa(cipherText []byte,priKey []byte) ([]byte, this.Error)
}

const RsaBit   int      = 1024 // rsa密钥长度

type Pem struct{}

func (c *CryptBlocks) EncryptToRsa(src []byte,pubKey []byte) ([]byte, this.Error) {
	// pem解码
	block, _ := pem.Decode(pubKey)

	// 使用x509标准转换成可以使用的公钥
	pk, _ := x509.ParsePKIXPublicKey(block.Bytes)

	// 强制转换
	publicKey := pk.(*rsa.PublicKey)

	// 使用公钥加密数据
	cipherText,err := rsa.EncryptPKCS1v15(rand.Reader,publicKey, src)
	if err != nil {
		return nil, this.SetError(err)
	}
	return cipherText, this.Nil
}

func (c *CryptBlocks) DecryptToRsa(cipherText []byte,priKey []byte) ([]byte, this.Error) {
	// 私钥解密
	// pem解密
	block, _ := pem.Decode(priKey)
	privatekey,err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, this.SetError(err)
	}

	// 使用私钥解密密文
	plainText,err := rsa.DecryptPKCS1v15(nil,privatekey,cipherText)
	if err != nil {
		return nil,this.SetError(err)
	}
	return plainText,this.Nil
}

func (p *Pem) CreateRsaPubKeyAndPriKey(pubf, prif io.Writer) this.Error {
	// 生成并在文件中写入可用的RSA公私钥
	// 创建私钥
	private, _ := rsa.GenerateKey(rand.Reader, RsaBit)
	// 获得公钥
	public := private.PublicKey
	// 使用x509标准转换为pem格式
	derText := x509.MarshalPKCS1PrivateKey(private)
	// 创建私钥结构体
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derText,
	}
	// 写入文件
	err1 := pem.Encode(prif, &block)
	if err1 != nil {
		return this.SetError(err1)
	}
	// 将公钥转换为pem格式
	derpText, _ := x509.MarshalPKIXPublicKey(&public)
	// 创建公钥结构体
	block = pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derpText,
	}
	// 写入文件
	err2 := pem.Encode(pubf, &block)
	if err2 != nil {
		return this.SetError(err2)
	}
	// nil Error
	return this.Nil
}

func (p *Pem) ReadPemKey(model int, file io.Reader) (string, this.Error) {
	panic("implement me")
}

// 匹配程序文件下有无生成好的公私钥
// pub,pri 传递path字符串
// 该层有panic()函数
func (p *Pem) MatchPubKeyAndPriKey(pub, pri string) bool {
	// 初始化日志
	log, er := logger.Std(nil)
	defer utils.ReCoverErrorAndPrint()
	defer log.Close()
	if er != this.Nil {
		panic(er)
	}
	// 读取时有错误则打印日志
	pubFile,err := os.Open(pub)
	if err != nil {
		log.StdErrorLog(err.Error())
		return false
	}
	defer pubFile.Close()
	priFile,err := os.Open(pri)
	if err != nil {
		log.StdErrorLog(err.Error())
		return false
	}
	defer priFile.Close()
	// 检查文件内容
	pubFileInfo,_ := pubFile.Stat()
	priFileInfo,_ := priFile.Stat()
	// 1024位密钥经过Base64格式化之后的标准长度,公钥=280,私钥>=887
	if pubFileInfo.Size() != 280 || priFileInfo.Size() < 887 {
		return false
	}
	return true
}
