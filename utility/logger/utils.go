package logger

import (
	"bytes"
	"encoding/json"

	"github.com/gogf/gf/v2/text/gstr"
)

// String 结构体或切片字符串输出
//
//	arg 必须是struct或者slice,原始数据类型禁止使用此接口
func String(arg interface{}) string {
	if arg == nil {
		return " "
	}
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(arg)
	if err != nil {
		return err.Error()
	}
	return gstr.TrimRight(buffer.String(), "\n")
}
