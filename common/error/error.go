package error

import (
	"reflect"
)

type Error string

func (e Error) Error() string { return string(e) }

/*
	常量为固定的重要错误
*/

const (
	Nil         Error = "nil"
	FileOpenErr Error = "File open failed"
	FileNil     Error = "File is empty"
	PathEmpty   Error = "Path is empty"
	KeyNumError Error = "Incorrect number of key bits"
)

/*
	全局变量为可附加信息的重要错误
*/

var LogError Error = "logout error "

/*
	允许自定义错误，参数类型为string或者error接口类型
	其他类型的参数将使用反射获取其值，会影响性能
	不推荐使用除error接口和string之外的类型
*/

func SetError(e interface{}) Error {
	switch e.(type) {
	case error:
		return Error(e.(error).Error())
	case string:
		return Error(e.(string))
	default:
		var value Error
		tp := reflect.TypeOf(e).Kind()
		switch tp {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = Error(tp.String())
			break
		case reflect.Bool:
			if reflect.ValueOf(e).Bool() {
				value = "1"
			} else {
				value = "0"
			}
			break
		case reflect.Float32, reflect.Float64:
			value = Error(tp.String())
			break
		case reflect.Struct, reflect.Map, reflect.Chan:
			value = Error(reflect.ValueOf(e).String())
		}
		return value
	}
}
