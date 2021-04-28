package test

import (
	"github.com/mengzushan/bups/common/error"
	"os"
	"testing"
)

func TestError(t *testing.T) {
	err := rError()
	if err != error.FileNil {
		t.Error("错误不匹配")
	} else if err == error.FileNil {
		t.Log("错误匹配")
	}
}

func rError() error.Error {
	return error.FileNil
}

// 测试错误处理
func Test_Error_Handing(t *testing.T) {
	// 人为制造一个打开文件的错误
	_, err := os.Open(".ahwbfiuawgifawbifwabiufb")
	var thisErr error.Error
	if err != nil {
		thisErr = error.SetError(err)
	}
	// 查看通过error接口类型转化为自定义错误处理的结果
	if thisErr.Error() != "" {
		t.Log("从error接口转换为自定义错误处理成功")
	} else {
		t.Error("从error接口转换为自定义错误处理失败")
	}
	// 通过string设置自定义错误处理
	thisErr = error.SetError("filed is filed")
	if thisErr.Error() != "" {
		t.Log("从string转换为自定义错误处理成功")
	} else {
		t.Error("从string转换为自定义错误处理失败")
	}
	// 从其他数据类型自定义错误处理
	//var data = 123456
	//thisErr = error.SetError(data)
	//if thisErr.Error() == strconv.Itoa(data) {
	//	t.Log("从int转换为自定义错误处理成功")
	//} else {
	//	t.Error("从int转换为自定义错误处理失败")
	//}
}
