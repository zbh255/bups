package utils

import (
	this "github.com/mengzushan/bups/common/error"
	"github.com/mengzushan/bups/common/logger"
	"os"
	"path/filepath"
	"strings"
)

// 该函数检测b string是否存在于a string中
func Equal(a, b string) bool {
	// 快速判断相同则直接返回
	// 快速判断相同则直接返回
	if a == b {
		return true
	}

	if len(a) / len(b) > 100 {
		return equalForKmp(a,b)
	} else {
		return equalForPtr(a,b)
	}
}

func equalForPtr(a, b string) bool {
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

// kmp的简单实现
func equalForKmp(a, b string) bool {
	aArr := []byte(a)
	bArr := []byte(b)
	// 构造部分匹配列表
	next := createNext(bArr)
	// 记录已经匹配的字符串数
	count := 0
	// 遍历要匹配的字符串
	for i := 0; i < len(aArr) ; i ++ {
		// 先判断匹配个数是否相同，相同直接返回true
		if count == len(b) {
			return true
		}
		// 判断a[i]字符中的值是否与b[count]相同，相同则count计数加1，进入下一轮循环
		if aArr[i] == bArr[count] {
			count++
			continue
		}
		// 判断a[i]不同于b[count]且count计数不为0的情况
		if aArr[i] != bArr[count] && count != 0 {
			// 移动位数 = 已匹配的字符数 - 对应的部分匹配值
			n := count - next[count - 1]
			// 更改i -> count的指针，相当于真实的移位
			count = count - n
			// 因为移位之后可能还会有一些不匹配的值所以还需要在使用原来a[i]的值进行匹配
			i = i - 1
		}
	}
	return false
}

func createNext(b []byte) []int {
	// 部分匹配列表
	pMatches := make([]int,len(b))
	// 原始指针
	raw := 0

	for raw <= len(b) - 1 {
		//// 前缀,后缀指针
		//low,high := 0,0
		// 前缀,后缀匹配列表
		lowList,highList := make([]string,raw),make([]string,raw)
		// 取出原始指针指向的值
		rVal := b[:raw + 1]
		//
		if len(rVal) == 1 {
			raw++
			continue
		}
		// 插入前缀后缀匹配列表
		for i := 0; i < raw ; i ++ {
			lowList[i] = string(rVal[:i + 1])
			highList[i] = string(rVal[i + 1:raw + 1])
		}
		// 部分匹配的字符串长度
		matchLens := 0
		for i := 0; i < raw; i++ {
			for j := 0; j < raw ; j ++ {
				if lowList[i] == highList[j] {
					if matchLens > len(lowList[i]) {
						// 长度小于之前匹配的字符串长度则不添加，因为要找出最长的长度
					} else {
						matchLens = len(lowList[i])
					}
				}
			}
		}
		// 将长度的值加入部分匹配列表对应的地方
		pMatches[raw] = matchLens
		raw++
	}
	return pMatches
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

func MatchUrlFile(url string) bool {
	file,err := os.Open(url)
	if err != nil {
		return false
	} else {
		defer file.Close()
		return true
	}
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
