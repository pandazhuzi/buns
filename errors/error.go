package errors

import (
	"fmt"
	"errors"
	"runtime/debug"
)

type BunsError struct {
	e error
}

func (e BunsError) Error() string{
	return e.e.Error()
}

func MakeError(e interface{}, v ...interface{}) error{

	switch err := e.(type) {
	case BunsError:
		return err
	case string:
		debug.PrintStack()
		return BunsError{
			e : errors.New(fmt.Sprintf(e.(string),v...))}
	case error:
		debug.PrintStack()
		return BunsError{
			e: err}
	}

	panic("error type not support")
}