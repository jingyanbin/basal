package basal

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFuncName(i interface{}, seps ...rune) string {
	if i == nil {
		return "nil"
	}
	// 获取函数名称
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})
	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}
