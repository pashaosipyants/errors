package errors

import (
	"errors"
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

// returns error with error code, if specified
// errcode is optional param. first of variadic parameters is used, else are ignored
func New(message string, errcode ...interface{}) error {
	return WrapAnnotated_skipstack(1, errors.New(message), "", errcode...)
}

// returns error with formatted message and without error code
func Errorf(format string, args ...interface{}) error {
	return WrapAnnotated_skipstack(1, fmt.Errorf(format, args...), "")
}

// returns error with formatted message and specified error code
func Codef(errcode interface{}, format string, args ...interface{}) error {
	return WrapAnnotated_skipstack(1, fmt.Errorf(format, args...), "", errcode)
}

// returns error with underlying err and error code, if specified
// if err already has error code(it means it was created earlier by one of this package's functions) and errcode is specified, it is overridden
// errcode is optional param. first of variadic parameters is used, else are ignored
func Wrap(err error, errcode ...interface{}) error {
	return WrapAnnotated_skipstack(1, err, "", errcode...)
}

// returns annotated error with underlying err and error code, if specified
// if err already has error code(it means it was created earlier by one of this package's functions) and errcode is specified, it is overridden
// errcode is optional param. first of variadic parameters is used, else are ignored
func WrapAnnotated(err error, annotation string, errcode ...interface{}) error {
	return WrapAnnotated_skipstack(1, err, annotation, errcode...)
}

/* for wrappers of this package */
// with these functions one can specify correct first stack frame to print in stack trace to skip stack frames of
// wrapper objects

func New_skipstack(message string, skip int, errcode ...interface{}) error {
	return WrapAnnotated_skipstack(skip+1, errors.New(message), "", errcode...)
}

func Errorf_skipstack(skip int, format string, args ...interface{}) error {
	return WrapAnnotated_skipstack(skip+1, fmt.Errorf(format, args...), "")
}

func Codef_skipstack(skip int, errcode interface{}, format string, args ...interface{}) error {
	return WrapAnnotated_skipstack(skip+1, fmt.Errorf(format, args...), "", errcode)
}

func Wrap_skipstack(skip int, err error, errcode ...interface{}) error {
	return WrapAnnotated_skipstack(skip+1, err, "", errcode...)
}

func WrapAnnotated_skipstack(skip int, err error, annotation string, errcode ...interface{}) error {
	if err == nil {
		return nil
	}

	reterr := new(_error)
	switch err1, ok := err.(*_error); {
	default:
		panic("assertion failed - it should be unreachable!!!\n Please, make an issue for developers of this package")
	case ok && len(errcode) == 0 && annotation == "":
		return err
	case ok:
		reterr.error = err1.error
		reterr.errcode = err1.errcode
		reterr.stack = err1.stack
		copy_map_annots(err1.annotations, &reterr.annotations)
	case !ok:
		reterr.error = err
		reterr.stack = callers(skip + 1)
	}

	if len(errcode) > 0 {
		reterr.errcode = errcode[0]
	}

	if annotation != "" {
		if reterr.annotations == nil {
			reterr.annotations = make(map[string]string)
		}
		pc, _, _, _ := runtime.Caller(skip + 1)
		reterr.addAnnotation(pc, annotation)
	}
	return reterr
}

/* methods */

func (f *_error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(s, "ERROR: %s\nERR CODE: %v\n", f.Error(), f.errcode)
		fmt.Fprint(s, f.annotatedStack)
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
// otherwise returns nil
func ErrCode(err error) interface{} {
	if err, ok := err.(*_error); ok {
		return err.ErrCode()
	}
	return nil
}

// if it is possible, gets underlying error wrapped with this package's error type
// otherwise returns err itself
func Cause(err error) error {
	if err, ok := err.(*_error); ok {
		return err.Cause()
	}
	return err
}

var _ Handleable = (*_error)(nil)
