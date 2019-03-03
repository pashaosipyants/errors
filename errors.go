package errors

import (
	"errors"
	"fmt"
)

// immutable type. none of functions changes it
type _error struct {
	error
	annotatedStack
	errcode interface{}
	suppressed error
}

// New returns error with error code, if specified.
//
// errcode is optional param. First of variadic parameters is used, else are ignored.
func New(message string, errcode ...interface{}) error {
	return WrapAnnotated_skipstack(1, errors.New(message), "", errcode...)
}

// Errorf returns error with formatted message and without error code.
func Errorf(format string, args ...interface{}) error {
	return WrapAnnotated_skipstack(1, fmt.Errorf(format, args...), "")
}

// Codef returns error with formatted message and specified error code.
func Codef(errcode interface{}, format string, args ...interface{}) error {
	return WrapAnnotated_skipstack(1, fmt.Errorf(format, args...), "", errcode)
}

// Wrap returns error with underlying err and error code, if specified.
// If err is nil returns nil.
// If err already has error code(it means it was created earlier by one of this package's functions) and
// errcode is specified, it is overridden.
//
// errcode is optional param. First of variadic parameters is used, else are ignored.
func Wrap(err error, errcode ...interface{}) error {
	return WrapAnnotated_skipstack(1, err, "", errcode...)
}

// WrapAnnotated returns annotated error with underlying err and error code, if specified.
// If err is nil returns nil.
// If err already has error code(it means it was created earlier by one of this package's functions) and
// errcode is specified, it is overridden.
// See https://godoc.org/github.com/pashaosipyants/errors/#hdr-Annotations to learn about annotations.
//
// errcode is optional param. First of variadic parameters is used, else are ignored.
func WrapAnnotated(err error, annotation string, errcode ...interface{}) error {
	return WrapAnnotated_skipstack(1, err, annotation, errcode...)
}

// WrapAndExtend returns error with underlying extender(err) and error code, if specified.
// If err is nil returns nil.
// If err already has error code(it means it was created earlier by one of this package's functions) and
// errcode is specified, it is overridden.
//
// errcode is optional param. First of variadic parameters is used, else are ignored.
func WrapAndExtend(err error, extender func(error) error, errcode ...interface{}) error {
	if err == nil {
		return nil
	}
	reterr := WrapAnnotated_skipstack(1, err, "", errcode...).(*_error)
	reterr.error = extender(reterr.error)
	return reterr
}

// Suppress returns error with underlying newerr, error code, if specified, and putting suppressed param as a
// suppressed error of returned error.
// Suppressed error is printed and can be obtained by Suppressed function of this package.
// If newerr is nil returns nil.
// If newerr already has error code(it means it was created earlier by one of this package's functions) and
// errcode is specified, it is overridden.
//
// errcode is optional param. First of variadic parameters is used, else are ignored.
func Suppress(suppressed, newerr error, errcode ...interface{}) error {
	return Suppress_skipstack(1, suppressed, newerr, errcode...)
}

/* for wrappers of this package */
// with these functions one can specify correct first stack frame to print in stack trace to skip stack frames of
// wrapper objects

// see https://godoc.org/github.com/pashaosipyants/errors/#hdr-Skipstack_management
func New_skipstack(message string, skip int, errcode ...interface{}) error {
	return WrapAnnotated_skipstack(skip+1, errors.New(message), "", errcode...)
}

// see https://godoc.org/github.com/pashaosipyants/errors/#hdr-Skipstack_management
func Errorf_skipstack(skip int, format string, args ...interface{}) error {
	return WrapAnnotated_skipstack(skip+1, fmt.Errorf(format, args...), "")
}

// see https://godoc.org/github.com/pashaosipyants/errors/#hdr-Skipstack_management
func Codef_skipstack(skip int, errcode interface{}, format string, args ...interface{}) error {
	return WrapAnnotated_skipstack(skip+1, fmt.Errorf(format, args...), "", errcode)
}

// see https://godoc.org/github.com/pashaosipyants/errors/#hdr-Skipstack_management
func Wrap_skipstack(skip int, err error, errcode ...interface{}) error {
	return WrapAnnotated_skipstack(skip+1, err, "", errcode...)
}

// see https://godoc.org/github.com/pashaosipyants/errors/#hdr-Skipstack_management
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
		reterr.addAnnotation(skip + 1, annotation)
	}
	return reterr
}

// see https://godoc.org/github.com/pashaosipyants/errors/#hdr-Skipstack_management
func Suppress_skipstack(skip int, suppressed, newerr error, errcode ...interface{}) error {
	if newerr == nil {
		return nil
	}

	reterr := WrapAnnotated_skipstack(skip+1, newerr, "", errcode...).(*_error)
	reterr.suppressed = suppressed
	return reterr
}

/* methods */

func (f *_error) Format(s fmt.State, verb rune) {
	printerr := func(err error) {
		if ferr, ok := err.(fmt.Formatter); ok {
			ferr.Format(s, verb)
		} else {
			fmt.Fprintf(s, "%"+string(verb), err)
		}
		fmt.Fprint(s, "\n")
	}

	// print error
	fmt.Fprintln(s, "ERROR:")
	printerr(f.error)
	fmt.Fprint(s, "\n")

	// print stack
	if verb == 'v' {
		fmt.Fprintln(s, "STACK:")
		fmt.Fprint(s, f.annotatedStack)
		fmt.Fprint(s, "\n")
	}

	// print suppressed
	if f.suppressed != nil {
		fmt.Fprintln(s, "SUPPRESSED:")
		printerr(f.suppressed)
		fmt.Fprint(s, "\n")
	}
}

func (w *_error) Cause() error { return w.error }

func (w *_error) ErrCode() interface{} { return w.errcode }

func (w *_error) Suppressed() error { return w.suppressed }

/* getters */

// If it is possible, gets errcode of the error.
// Otherwise returns nil.
func ErrCode(err error) interface{} {
	if err, ok := err.(*_error); ok {
		return err.ErrCode()
	}
	return nil
}

// If it is possible, gets underlying error wrapped with this package's error type.
// Otherwise returns err itself.
func Cause(err error) error {
	if err, ok := err.(*_error); ok {
		return err.Cause()
	}
	return err
}

// If it is possible, gets suppressed error.
// Otherwise returns nil.
func Suppressed(err error) error {
	if err, ok := err.(*_error); ok {
		return err.Suppressed()
	}
	return nil
}

var _ Handleable = (*_error)(nil)
