package basal

import (
	"fmt"
	"os"
)

type ILogger interface {
	ErrorF(format string, v ...interface{})
}

type fileLogger struct {
	f *os.File
}

func (m *fileLogger) ErrorF(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	sLen := len(s)
	if (sLen > 0 && s[sLen-1] != '\n') || sLen == 0 {
		m.f.WriteString("[ERROR] " + s + "\n")
	}
}

func newFileLogger(f *os.File) *fileLogger {
	return &fileLogger{f: f}
}

var log = newFileLogger(os.Stdout)

func GetLogger() ILogger {
	return log
}
