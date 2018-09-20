package errors

import (
	"fmt"
	"io"
	"runtime"
)

// immutable type. none of functions changes it
type _error struct {
	error
	annotatedStack
	errcode interface{}
}

func New(message string, errcode ...interface{}) error {
	return WrapWithMessage_skipstack(1, fmt.Errorf(message), "", errcode)
}

func Errorf(format string, args ...interface{}) error {
	return WrapWithMessage_skipstack(1, fmt.Errorf(format, args...), "")
}

func ErrorfWithCode(errcode interface{}, format string, args ...interface{}) error {
	return WrapWithMessage_skipstack(1, fmt.Errorf(format, args...), "", errcode)
}

func Wrap(err error, errcode ...interface{}) error {
	return WrapWithMessage_skipstack(1, err, "", errcode)
}

func WithMessage(err error, msg string) error {
	return WrapWithMessage_skipstack(1, err, msg)
}

func WrapWithMessage(err error, msg string, errcode ...interface{}) error {
	return WrapWithMessage_skipstack(1, err, msg, errcode)
}

/* for wrappers of this package */

func New_skipstack(message string, skip int, errcode ...interface{}) error {
	return WrapWithMessage_skipstack(skip + 1, fmt.Errorf(message), "", errcode)
}

func Errorf_skipstack(skip int, format string, args ...interface{}) error {
	return WrapWithMessage_skipstack(skip + 1, fmt.Errorf(format, args...), "")
}

func ErrorfWithCode_skipstack(skip int, errcode interface{}, format string, args ...interface{}) error {
	return WrapWithMessage_skipstack(skip + 1, fmt.Errorf(format, args...), "", errcode)
}

func Wrap_skipstack(skip int, err error, errcode ...interface{}) error {
	return WrapWithMessage_skipstack(skip + 1, err, "", errcode)
}

func WithMessage_skipstack(skip int, err error, msg string) error {
	return WrapWithMessage_skipstack(skip + 1, err, msg)
}

func WrapWithMessage_skipstack(skip int, err error, msg string, errcode ...interface{}) error {
	if err == nil { return nil }

	reterr := new(_error)
	switch err1, ok := err.(*_error); {
		default:
			panic("assertion failed - it should be unreachable!!!")
		case ok && len(errcode) == 0 && msg == "":
			return err
		case ok:
			reterr.error = err1.error
			reterr.errcode = err1.errcode
			reterr.stack = err1.stack
			copy_map_msgs(err1.messages, &reterr.messages)
		case !ok:
			reterr.error = err
			reterr.stack = callers(skip + 1)
	}

	if len(errcode) > 0 {
		reterr.errcode = errcode[0]
	}

	if msg != "" {
		if reterr.messages == nil { reterr.messages = make(map[string]string) }
		pc, _, _, _ := runtime.Caller(skip + 1)
		reterr.addMsg(pc, msg)
	}
	return reterr
}

/* methods */

func (f *_error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(s, "ERROR: %s\nERR CODE: %v\n", f.Error(), f.errcode)
		f.annotatedStack.Format(s, verb)
	case 's':
		io.WriteString(s, f.Error())
	case 'q':
		fmt.Fprintf(s, "%q", f.Error())
	}
}

func (w *_error) Cause() error { return w.error }

func (w *_error) ErrCode() interface{} { return w.errcode }

/* getters */

// if it is possible, gets errcode of error.
// otherwise returns empty string("")
func ErrCode(err error) string {
	type errcoder interface {
		ErrCode() string
	}

	if err, ok := err.(errcoder); ok {
		return err.ErrCode()
	}
	return ""
}

// if it is possible, gets underlying error wrapped with this package's error type
// otherwise returns err itself
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	if err, ok := err.(causer); ok {
		return err.Cause()
	}
	return err
}

var _ Handleable = new(_error)
