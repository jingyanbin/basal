package base

import (
	"errors"
	"fmt"
)

func NewError(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}