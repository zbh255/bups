package encry

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

type Crypt interface {
	EncryptToAes(src []byte, key []byte) ([]byte, error)
	DecryptToAes(dst []byte, key []byte) ([]byte, error)
	CryptToRsa
}

type CryptBlocks struct{}

func (c *CryptBlocks) EncryptToAes(src []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	src = c.pkcs7Padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(src, src)
	return src, nil
}

func (c *CryptBlocks) DecryptToAes(dst []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(dst, dst)
	dst = c.pkcs7UnPadding(dst, block.BlockSize())
	return dst, nil
}

// 填充明文
func (c *CryptBlocks) pkcs7Padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

// 去掉填充
func (c *CryptBlocks) pkcs7UnPadding(src []byte, blockSize int) []byte {
	return src[:len(src)-(blockSize-len(src)%blockSize)]
}
