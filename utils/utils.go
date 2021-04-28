package utils

import (
	this "github.com/mengzushan/bups/common/error"
	"github.com/mengzushan/bups/common/logger"
	"os"
	"path/filepath"
	"strings"
)

// 该函数检测b string是否存在于a string中
// 采用与Equals不同的实现方式
// 在比较较长的字符串时比较占优势
func Equal(a, b string) bool {
	// 快速判断相同则直接返回
	// 快速判断相同则直接返回
	if a == b {
		return true
	}
	buf1 := []byte(a)
	buf2 := []byte(b)
	buf2Num := len(buf2)
	if buf2Num > len(buf1) {
		return false
	}
	var bNode int
	var equalNum int
	bNode = len(buf2) - 1
	for i := 0; i < len(buf1); i++ {
		if bNode > len(buf1)-1 {
			return false
		}
		if buf1[i] == buf2[0] && buf1[bNode] == buf2[buf2Num-1] {
			for j := i + 1; j < bNode; j++ {
				if buf1[j] == buf2[j-i] {
					equalNum++
				}
				if equalNum+2 == buf2Num {
					return true
				}
			}
		}
		bNode++
	}
	return false
}

func Equals(a, b string) bool {
	buf1 := []byte(a)
	buf2 := []byte(b)

	equalNum := 0
	for i := 0; i < len(buf1); i++ {
		if len(buf1)-i >= len(buf2) {
			if buf1[i] == buf2[0] {
				for j := 0; j < len(buf2); j++ {
					if buf1[i+j] == buf2[j] {
						equalNum++
					}
					if equalNum+1 == len(buf2) {
						return true
					}
				}
			}
		} else {
			return false
		}
	}
	return true
}

// 本身函数中已存在log功能，该恢复函数不往日志上写错误
func ReCoverErrorAndPrint() {
	err := recover()
	switch err.(type) {
	case this.Error:
		// 日志错误不可恢复
		if Equal(err.(this.Error).Error(), this.LogError.Error()) {
			panic(err.(this.Error))
		} else {
			println(err.(this.Error).Error())
		}
	case error:
		println(err.(error).Error())
	}
}

// 该恢复函数用于没有使用打印日志功能的函数
func ReCoverErrorAndLog() {
	err := recover()
	log, customErr := logger.Std(nil)
	defer log.Close()
	defer ReCoverErrorAndPrint()
	if customErr != this.Nil {
		panic(err)
	}
	switch err.(type) {
	case this.Error:
		log.StdErrorLog(err.(this.Error).Error())
		break
	case error:
		log.StdErrorLog(err.(error).Error())
		break
	}
}

// 清理退出程序时的缓存
func CleanUpCache(cachePath string) this.Error {
	// 遍历文件夹并删除
	err := filepath.Walk(cachePath, func(path string, info os.FileInfo, err error) error {
		if path != cachePath {
			err = os.Remove(path)
		}
		return err
	})
	if err != nil {
		return this.SetError(err)
	}
	return this.Nil
}

func MatchAppInfo() {

}

// 按.分隔url -> "https://examplebucket-1250000000.cos.COS_REGION.myqcloud.com"
func SplitCosBucketUrl(url string) string {
	url2 := strings.SplitN(url, ".", 5)
	return url2[2]
}

// 判断字符串切片中的内容是否相同
func EqualToStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	var eqNum int
	for i := 0; i < len(a); i++ {
		for _, v := range b {
			if a[i] == v {
				eqNum++
			}
			if eqNum == len(a) {
				return true
			}
		}
	}
	return false
}
