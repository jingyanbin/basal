package basal

import (
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
)

const exceptionSkip = 3

// 调用信息短文件名
func CallerShort(skip int) (file string, line int) {
	var ok bool
	_, file, line, ok = runtime.Caller(skip)
	if ok {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
	} else {
		file = "???"
		line = 0
	}
	return
}

// 调用信息长文件名
func Caller(skip int) (file string, line int) {
	var ok bool
	_, file, line, ok = runtime.Caller(skip)
	if !ok {
		file = "???"
		line = 0
	}
	return
}

func CallerInFunc(skip int) (name string, file string, line int) {
	var pc uintptr
	var ok bool
	pc, file, line, ok = runtime.Caller(skip)
	if ok {
		inFunc := runtime.FuncForPC(pc)
		name = inFunc.Name()
	} else {
		file = "???"
		name = "???"
	}
	return
}

var reLine = regexp.MustCompile(`^panic\([a-z 0-9]+,\s*[a-z 0-9]+\)$`)

func CallerLineStack(stack string) (name string, file string) {
	name = "???"
	file = "???"
	stackArr := strings.SplitN(stack, "panic.go", 2)
	if len(stackArr) != 2 {
		return
	}
	stackLines := strings.SplitN(stackArr[1], "\n", 4)
	if len(stackLines) != 4 {
		return
	}
	//name = strings.Trim(stackLines[1], "\n")
	//file = strings.Trim(stackLines[2], "\n")
	name = strings.TrimSpace(stackLines[1])
	file = strings.TrimSpace(stackLines[2])
	return
}

func toError(r interface{}) (err error) {
	switch x := r.(type) {
	case string:
		err = NewError(x)
	case error:
		err = x
	default:
		err = NewError("unknown error: %v", x)
	}
	return
}

func formatStack(name, file string, err string, stack []byte) *Buffer {
	buf := NewBuffer(160 + len(stack) + len(name))
	buf.AppendStrings("exception: ", err, "\nfile: ", file)
	buf.AppendStrings("\nfunc: ", name, "\n")
	buf.AppendBytes(stack...)
	return buf
}

func exception(catch func(e error)) {
	if err := recover(); err != nil {
		catch(toError(err))
	}
}

func Exception(catchs ...func(stack string, e error)) {
	if err := recover(); err != nil {
		if len(catchs) == 0 {
			return
		}
		info := debug.Stack()
		name, file := CallerLineStack(string(info))
		myErr := toError(err)
		myBuf := formatStack(name, file, myErr.Error(), info)
		defer myBuf.Free()
		for _, catch := range catchs {
			if catch == nil {
				continue
			}
			catch(myBuf.ToString(), myErr)
		}
	}
}

func Try(f func(), catch func(stack string, e error)) {
	defer Exception(catch)
	f()
}
