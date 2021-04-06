package basal

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func IsUTF8(buf []byte) bool {
	nBytes := 0
	for i := 0; i < len(buf); i++ {
		if nBytes == 0 {
			if (buf[i] & 0x80) != 0 { //与操作之后不为0，说明首位为1
				b := buf[i]
				for (b & 0x80) != 0 {
					b <<= 1  //左移一位
					nBytes++ //记录字符共占几个字节
				}
				if nBytes < 2 || nBytes > 6 { //因为UTF8编码单字符最多不超过6个字节
					return false
				}
				nBytes-- //减掉首字节的一个计数
			}
		} else {                     //处理多字节字符
			if buf[i]&0xc0 != 0x80 { //判断多字节后面的字节是否是10开头
				return false
			}
			nBytes--
		}
	}
	return nBytes == 0
}

func Type(value interface{}) reflect.Type {
	return reflect.TypeOf(value)
}

func ToString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int, int8, int16, int32, int64:
		n, err := ToInt64(v)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(n, 10), nil
	case uint, uint8, uint16, uint32, uint64:
		n, err := ToInt64(v)
		if err != nil {
			return "", err
		}
		return strconv.FormatUint(uint64(n), 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case map[string]interface{}, []interface{}:
		b, err := json.Marshal(v)
		return string(b), err
	case []byte:
		return string(v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
	return "", NewError("ToString value type error: %v", Type(value))
}

func ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	case []byte:
		return strconv.ParseFloat(string(v), 64)
	default:
		switch value := reflect.ValueOf(v); value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n := value.Int()
			return float64(n), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n := value.Uint()
			return float64(n), nil
		case reflect.Float64, reflect.Float32:
			return value.Float(), nil
		case reflect.String:
			return strconv.ParseFloat(value.String(), 64)
		case reflect.Slice:
			return strconv.ParseFloat(string(value.Bytes()), 64)
		}
	}
	return 0, NewError("ToFloat64 value type error: %v", Type(value))
}

func ToFloat32(value interface{}) (float32, error) {
	v, err := ToFloat64(value)
	return float32(v), err
}

func ToInt64(value interface{}) (int64, error) {
	switch n := value.(type) {
	case int:
		return int64(n), nil
	case int8:
		return int64(n), nil
	case int16:
		return int64(n), nil
	case int32:
		return int64(n), nil
	case int64:
		return int64(n), nil
	case uint:
		return int64(n), nil
	case uint8:
		return int64(n), nil
	case uint16:
		return int64(n), nil
	case uint32:
		return int64(n), nil
	case uint64:
		return int64(n), nil
	case float64:
		return int64(n), nil
	case float32:
		return int64(n), nil
	case string:
		f, err := strconv.ParseFloat(n, 64)
		return int64(f), err
		//return strconv.ParseInt(n, 10, 64)
	case []byte:
		f, err := strconv.ParseFloat(string(n), 64)
		return int64(f), err
		//return strconv.ParseInt(string(n), 10, 64)
	default:
		switch value := reflect.ValueOf(n); value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return value.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return int64(value.Uint()), nil
		case reflect.Float64, reflect.Float32:
			return int64(value.Float()), nil
		case reflect.String:
			f, err := strconv.ParseFloat(value.String(), 64)
			return int64(f), err
			//return strconv.ParseInt(value.String(), 10, 64)
		case reflect.Slice:
			f, err := strconv.ParseFloat(string(value.Bytes()), 64)
			return int64(f), err
			//return strconv.ParseInt(string(value.Bytes()), 10, 64)
		}
	}
	return 0, NewError("ToInt64 value type error: %v", Type(value))
}

func ToInt32(value interface{}) (int32, error) {
	v, err := ToInt64(value)
	return int32(v), err
}

func ToInt16(value interface{}) (int16, error) {
	v, err := ToInt64(value)
	return int16(v), err
}

func ToInt8(value interface{}) (int8, error) {
	v, err := ToInt64(value)
	return int8(v), err
}

func ToInt(value interface{}) (int, error) {
	v, err := ToInt64(value)
	return int(v), err
}

func ToUint64(value interface{}) (uint64, error) {
	v, err := ToInt64(value)
	return uint64(v), err
}

func ToUint32(value interface{}) (uint32, error) {
	v, err := ToInt64(value)
	return uint32(v), err
}

func ToUint16(value interface{}) (uint16, error) {
	v, err := ToInt64(value)
	return uint16(v), err
}

func ToUint8(value interface{}) (uint8, error) {
	v, err := ToInt64(value)
	return uint8(v), err
}

func ToUint(value interface{}) (uint, error) {
	v, err := ToInt64(value)
	return uint(v), err
}

func Float64(value interface{}) (float64, error) {
	v, ok := value.(float64)
	if !ok {
		return 0, NewError("Float64 value type error: %v", Type(value))
	}
	return v, nil
}

func Float32(value interface{}) (float32, error) {
	v, ok := value.(float32)
	if !ok {
		return 0, NewError("Float32 value type error: %v", Type(value))
	}
	return v, nil
}

func Int64(value interface{}) (int64, error) {
	v, ok := value.(int64)
	if !ok {
		return 0, NewError("Int64 value type error: %v", Type(value))
	}
	return v, nil
}

func Uint64(value interface{}) (uint64, error) {
	v, ok := value.(uint64)
	if !ok {
		return 0, NewError("Uint64 value type error: %v", Type(value))
	}
	return v, nil
}

func Int32(value interface{}) (int32, error) {
	v, ok := value.(int32)
	if !ok {
		return 0, NewError("Int32 value type error: %v", Type(value))
	}
	return v, nil
}

func Uint32(value interface{}) (uint32, error) {
	v, ok := value.(uint32)
	if !ok {
		return 0, NewError("Uint32 value type error: %v", Type(value))
	}
	return v, nil
}

func Int16(value interface{}) (int16, error) {
	v, ok := value.(int16)
	if !ok {
		return 0, NewError("Int16 value type error: %v", Type(value))
	}
	return v, nil
}

func Uint16(value interface{}) (uint16, error) {
	v, ok := value.(uint16)
	if !ok {
		return 0, NewError("Uint16 value type error: %v", Type(value))
	}
	return v, nil
}

func Int8(value interface{}) (int8, error) {
	v, ok := value.(int8)
	if !ok {
		return 0, NewError("Int8 value type error: %v", Type(value))
	}
	return v, nil
}

func Uint8(value interface{}) (uint8, error) {
	v, ok := value.(uint8)
	if !ok {
		return 0, NewError("Uint8 value type error: %v", Type(value))
	}
	return v, nil
}

func Int(value interface{}) (int, error) {
	v, ok := value.(int)
	if !ok {
		return 0, NewError("Int value type error: %v", Type(value))
	}
	return v, nil
}

func Uint(value interface{}) (uint, error) {
	v, ok := value.(uint)
	if !ok {
		return 0, NewError("Uint value type error: %v", Type(value))
	}
	return v, nil
}

//const UINT64_MIN uint64 = 0
//const UINT64_MAX uint64 = ^UINT64_MIN
//const INT64_MIN = ^UINT64_MAX
//const INT64_MAX  = int64(^uint64(0)>>1)

const UINT32_MIN uint32 = 0
const UINT32_MAX uint32 = ^UINT32_MIN
const INT32_MIN = ^UINT32_MAX
const INT32_MAX = int32(^uint32(0) >> 1)

func AtoInt64(s string) (x int64, err error) {
	neg := false
	if s == "" {
		return 0, NewError("param error: %s", s)
	}

	if s[0] == '-' || s[0] == '+' {
		neg = s[0] == '-'
		s = s[1:]
	} else if s[0] < '0' || s[0] > '9' {
		return 0, NewError("param error: %s", s)
	}

	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > (1<<63-1)/10 {
			// overflow
			return 0, NewError("param error: overflow %v", s)
		}
		x = x*10 + int64(c) - '0'
		if x < 0 {
			// overflow
			return 0, NewError("param error: overflow %v", s)
		}
	}
	if neg {
		x = -x
	}
	return x, nil
}

func AtoInt(s string) (int, error) {
	x, err := AtoInt64(s)
	return int(x), err
}

func AbsInt64(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func AbsInt32(n int32) int32 {
	y := n >> 31
	return (n ^ y) - y
}

func AbsFloat32(n float32) float32 {
	if n < 0 {
		return -n
	}
	return n
}

func AbsFloat64(n float64) float64 {
	if n < 0 {
		return -n
	}
	return n
}

func Round(value float64, digit int) float64 {
	p10 := math.Pow10(digit)
	return math.Trunc((value+0.5/p10)*p10) / p10
}
