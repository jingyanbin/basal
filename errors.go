package basal

import (
	"errors"
	"fmt"
)

func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func NewError(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}
