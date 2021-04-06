package base

import (
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
)

const exceptionSkip = 5

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

func CallerLineStack(stack string) (name string, file string, line int, success bool) {
	stackLines := strings.Split(stack, "\n")
	max := len(stackLines)
	for i, v := range stackLines {
		if reLine.MatchString(v) {
			if i+3 < max {
				fls := strings.Trim(stackLines[i+3], "\t")
				fileLines := strings.Split(fls, " ")[0]
				index := strings.LastIndex(fileLines, ":")
				if index == -1 {
					return
				}
				file = fileLines[:index]
				var err error
				line, err = AtoInt(fileLines[index+1:])
				if err != nil {
					return
				}
				name = stackLines[i+2] // strings.Split(stackLines[i + 2], "(")[0]
				success = true
				return
			} else {
				return
			}
		}
	}
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

func formatStack(name, file string, line int, err string, stack []byte) *Buffer {
	buf := NewBuffer(160 + len(stack) + len(name))
	buf.AppendStrings("exception panic: ", err, "\nfile: ", file, ":")
	buf.AppendInt(line, 0)
	buf.AppendStrings("\nfunc: ", name, "\n")
	buf.AppendBytes(stack...)
	return buf
}

func exception(catch func(e error)) {
	if err := recover(); err != nil {
		catch(toError(err))
	}
}

func Exception(catch func(stack string, e error)) {
	if err := recover(); err != nil {
		if catch == nil{
			return
		}
		info := debug.Stack()
		name, file, line, success := CallerLineStack(string(info))
		if !success {
			name, file, line = CallerInFunc(exceptionSkip)
		}
		myErr := toError(err)
		myBuf := formatStack(name, file, line, myErr.Error(), info)
		defer myBuf.Free()
		catch(myBuf.ToString(), myErr)
	}
}

func Try(f func(), Catch func(stack string, e error)) {
	defer Exception(Catch)
	f()
}
