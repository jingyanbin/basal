package basal

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
)

const JSON_APPEND = -1          //追加
const JSON_APPEND_IN_FRONT = -2 //追加在前面
const JSON_INDEX_LAST = -3      //最后的位置
const JSON_INDEX_FIRST = -4     //第一的位置

type Json struct {
	data interface{}
}

func (my *Json) String() string {
	s, _ := my.ToString(false)
	return s
}

func (my *Json) IsNil() bool {
	return my.data == nil
}

func (my *Json) ToString(indent bool) (string, error) {
	return ToString(my.data, indent)
}

func (my *Json) ToInt64() (int64, error) {
	return ToInt64(my.data)
}

func (my *Json) ToInt32() (int32, error) {
	return ToInt32(my.data)
}

func (my *Json) ToFloat64() (float64, error) {
	return ToFloat64(my.data)
}

func (my *Json) ToFloat32() (float32, error) {
	return ToFloat32(my.data)
}

func (my *Json) Float64() (float64, bool) {
	if number, ok := my.data.(json.Number); ok {
		v, err := number.Float64()
		if err != nil {
			return 0, false
		}
		return v, true
	}
	return 0, false
}

func (my *Json) Float32() (float32, bool) {
	if v, ok := my.Float64(); ok {
		if v > math.MaxFloat32 || v < -math.MaxFloat32 {
			return 0, false
		}
		return float32(v), true
	}
	return 0, false
}

func (my *Json) Int64() (int64, bool) {
	if number, ok := my.data.(json.Number); ok {
		v, err := number.Int64()
		if err != nil {
			return 0, false
		}
		return v, true
	}
	return 0, false
}

func (my *Json) Int32() (int32, bool) {
	if v, ok := my.Int64(); ok {
		if v > math.MaxInt32 || v < math.MinInt32 {
			return 0, false
		}
		return int32(v), true
	}
	return 0, false
}

func (my *Json) Int() (int, bool) {
	if v, ok := my.Int64(); ok {
		return int(v), true
	}
	return 0, false
}

func (my *Json) Bool() (v bool, ok bool) {
	v, ok = my.data.(bool)
	return
}

func (my *Json) Slice() ([]interface{}, error) {
	v, ok := my.data.([]interface{})
	if ok {
		return v, nil
	} else {
		return nil, NewError("[]interface{} value type error: %v", Type(my.data))
	}
}

func (my *Json) Map() (map[string]interface{}, error) {
	v, ok := my.data.(map[string]interface{})
	if ok {
		return v, nil
	} else {
		return nil, NewError("map[string]interface{} value type error: %v", Type(my.data))
	}
}

func (my *Json) Bytes() ([]byte, error) {
	if my.data == nil {
		return nil, NewError("json is nil")
	}
	return json.Marshal(my.data)
}

func (my *Json) GetJson(keys ...interface{}) *Json {
	return &Json{my.Get(keys...)}
}

func (my *Json) Get(keys ...interface{}) interface{} {
	var v = my.data
	var ok bool
	for _, key := range keys {
		switch k := key.(type) {
		case string:
			v, ok, _ = findMapKey(v, k)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			index, _ := ToInt(k)
			v, ok, _ = findSliceIndex(v, index)
		default:
			ok = false
		}
		if !ok {
			return nil
		}
	}
	return v
}

func (my *Json) Loads(js interface{}) error {
	obj, err := NewJson(js)
	if err != nil {
		return err
	}
	my.data = obj.data
	return nil
}

func (my *Json) create(keys []interface{}) (interface{}, error) {
	length := len(keys)
	if length < 2 {
		return nil, NewError("json create error: keys num less 2, keys=%v", keys)
	}

	var lastRoot interface{}
	pos := length - 1
	lastRoot = keys[pos]
	for i := pos - 1; i >= 0; i-- {
		switch k := keys[i].(type) {
		case string:
			parent := map[string]interface{}{}
			parent[k] = lastRoot
			lastRoot = parent
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			index, _ := ToInt(k)
			if index == JSON_APPEND || index == JSON_APPEND_IN_FRONT {
				index = 0
			}
			if index == 0 {
				lastRoot = []interface{}{lastRoot}
			} else {
				return nil, NewError("json create error: slice out of range, keys=%v, index=%v", keys, i)
			}
		default:
			return nil, NewError("json create error: not found key type, keys=%v, index=%v, type=%v", keys, i, Type(keys[i]))
		}
	}
	return lastRoot, nil
}

func (my *Json) set(root interface{}, args []interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, NewError("json set error: args num less 2")
	}
	switch data := root.(type) {
	case *interface{}:
		switch v := (*data).(type) {
		case []interface{}:
			var index int
			switch idx := args[0].(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				index, _ = ToInt(idx)
			default:
				return nil, NewError("json set error: key not is index %v", args)
			}
			maxLen := len(v)
			if index > maxLen || index < JSON_INDEX_FIRST {
				return nil, NewError("json set error: out of range[%v, %v] error: index=%v", JSON_INDEX_FIRST, maxLen, index)
			}
			if len(args) == 2 {
				value := args[1]
				if index == JSON_APPEND || index == maxLen {
					v = append(v, value)
				} else if index == JSON_APPEND_IN_FRONT {
					v = append([]interface{}{value}, v...)
				} else if index == JSON_INDEX_LAST {
					if maxLen > 0 {
						v[maxLen-1] = value
					}
				} else if index == JSON_INDEX_FIRST {
					v[0] = value
				} else {
					v[index] = value
				}
			} else {
				if index == JSON_APPEND || index == maxLen {
					value, err := my.create(args[1:])
					if err != nil {
						return nil, err
					}
					v = append(v, value)
				} else if index == JSON_APPEND_IN_FRONT {
					value, err := my.create(args[1:])
					if err != nil {
						return nil, err
					}
					v = append([]interface{}{value}, v...)
				} else if index == JSON_INDEX_LAST {
					value, err := my.set(&v[maxLen-1], args[1:])
					if err != nil {
						return nil, err
					}
					v[maxLen-1] = value
				} else if index == JSON_INDEX_FIRST {
					value, err := my.set(&v[0], args[1:])
					if err != nil {
						return nil, err
					}
					v[0] = value
				} else {
					value, err := my.set(&v[index], args[1:])
					if err != nil {
						return nil, err
					}
					v[index] = value
				}
			}
			return v, nil

		case map[string]interface{}:
			key, ok := args[0].(string)
			if !ok {
				return nil, NewError("json set error: key not is string %v", args)
			}
			if len(args) == 2 {
				v[key] = args[1]
			} else {
				next, ok := v[key]
				if ok {
					value, err := my.set(&next, args[1:])
					if err != nil {
						return nil, err
					}
					v[key] = value
				} else {
					value, err := my.create(args[1:])
					if err != nil {
						return nil, err
					}
					v[key] = value
				}
			}
			return v, nil
		}
	}
	return nil, NewError("json set type error: args=%v", args)
}

func (my *Json) Set(args ...interface{}) error {
	if my.data == nil {
		value, err := my.create(args)
		if err != nil {
			return NewError("json root is nil create error: %v", err)
		}
		my.data = value
		return nil
	} else {
		value, err := my.set(&my.data, args)
		if err != nil {
			return err
		}
		my.data = value
		return nil
	}
}

func (my *Json) delete(root interface{}, keys []interface{}) (interface{}, bool, bool) {
	if len(keys) == 0 {
		return nil, false, true
	}
	switch data := root.(type) {
	case *interface{}:
		switch v := (*data).(type) {
		case []interface{}:
			switch idx := keys[0].(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				index, err := ToInt(idx)
				if err == nil {
					maxLen := len(v)
					if index == JSON_INDEX_LAST {
						index = maxLen - 1
					} else if index == JSON_INDEX_FIRST {
						index = 0
					}
					if index >= 0 && index < maxLen {
						value, success, found := my.delete(&v[index], keys[1:])
						if value != nil {
							v[index] = value
						}
						if found {
							v = append(v[:index], v[index+1:]...)
							return v, true, false
						} else {
							return nil, success, false
						}
					}
				}
			}
		case map[string]interface{}:
			switch key := keys[0].(type) {
			case string:
				next, ok := v[key]
				if ok {
					value, success, found := my.delete(&next, keys[1:])
					if value != nil {
						v[key] = value
					}
					if found {
						delete(v, key)
						return v, true, false
					} else {
						return nil, success, false
					}
				}
			}
		}
	}
	return nil, false, false
}

func (my *Json) Delete(keys ...interface{}) bool {
	if my.data != nil {
		value, success, _ := my.delete(&my.data, keys)
		if value != nil {
			my.data = value
		}
		return success
	}
	return false
}

func (my *Json) Clear() {
	my.data = nil
}

func findMapKey(data interface{}, key string) (v interface{}, ok bool, err error) {
	switch m := data.(type) {
	case map[string]interface{}:
		v, ok = m[key]
	default:
		ok = false
		err = NewError("not is map[string]interface{}, type=%v", Type(m))
	}
	return
}

func findSliceIndex(data interface{}, index int) (interface{}, bool, error) {
	if data == nil {
		return nil, false, nil
	}
	slice, ok := data.([]interface{})
	if !ok {
		return nil, false, NewError("not is []interface{}, type=%v", Type(data))
	}
	maxLen := len(slice)
	if index >= 0 && index < maxLen {
		return slice[index], true, nil
	} else if index == JSON_INDEX_LAST && maxLen > 0 {
		return slice[maxLen-1], true, nil
	} else if index == JSON_INDEX_FIRST && maxLen > 0 {
		return slice[0], true, nil
	} else {
		return nil, false, nil
	}
}

func NewJson(js interface{}) (*Json, error) {
	switch v := js.(type) {
	case string:
		return loadJson([]byte(v))
	case []byte:
		return loadJson(v)
	case map[string]interface{}, []interface{}:
		return linkJson(v)
	case *os.File:
		bs, err := ioutil.ReadAll(v)
		if err != nil {
			return nil, err
		}
		return loadJson(bs)
	}
	return nil, NewError("new json type error: %v", Type(js))
}

func loadJson(v []byte) (*Json, error) {
	js := &Json{}
	decoder := json.NewDecoder(bytes.NewBuffer(v))
	decoder.UseNumber()
	err := decoder.Decode(&js.data)
	if err != nil {
		return nil, err
	}
	return js, nil
}

func linkJson(js interface{}) (*Json, error) {
	switch v := js.(type) {
	case map[string]interface{}, []interface{}:
		return &Json{v}, nil
	}
	return nil, NewError("link json type error: %v", Type(js))
}

func LoadJsonFileTo(jsFileName string, toPtr interface{}) error {
	//if isNil, typ := IsNilPointer(v); isNil {
	//	return NewError("LoadJsonFileTo IsNilPointer: %s", typ)
	//}
	data, err := ioutil.ReadFile(jsFileName)
	if err != nil {
		return err
	}
	return LoadJsonBytesTo(data, toPtr)
}

func LoadJsonBytesTo(js []byte, toPtr interface{}) error {
	dec := json.NewDecoder(bytes.NewBuffer(js))
	dec.UseNumber()
	err := dec.Decode(toPtr)
	if err != nil {
		return err
	}
	return nil
}

func LoadJsonStringTo(js string, toPtr interface{}) error {
	return LoadJsonBytesTo([]byte(js), toPtr)
}
