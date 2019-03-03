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

type setup struct {
	errcode *interface{}
	annotation *string
	extender func(error) error
}

type option func(*setup)

func AddErrCode(errcode interface{}) option {
	return func(o *setup) {
		o.errcode = &errcode
	}
}

func AddAnnotation(annotation string) option {
	return func(o *setup) {
		o.annotation = &annotation
	}
}

func AddExtender(extender func(error) error) option {
	return func(o *setup) {
		o.extender = extender
	}
}

// New returns error with error code, if specified.
//
// errcode is optional param. First of variadic parameters is used, else are ignored.
func New(message string, opts ...option) error {
	return Wrap_skipstack(1, errors.New(message), opts...)
}

// Wrap returns error with underlying err and error code, if specified.
// If err is nil returns nil.
// If err already has error code(it means it was created earlier by one of this package's functions) and
// errcode is specified, it is overridden.
//
// errcode is optional param. First of variadic parameters is used, else are ignored.
func Wrap(err error, opts ...option) error {
	return Wrap_skipstack(1, err, opts...)
}

// Suppress returns error with underlying newerr, error code, if specified, and putting suppressed param as a
// suppressed error of returned error.
// Suppressed error is printed and can be obtained by Suppressed function of this package.
// If newerr is nil returns nil.
// If newerr already has error code(it means it was created earlier by one of this package's functions) and
// errcode is specified, it is overridden.
//
// errcode is optional param. First of variadic parameters is used, else are ignored.
func Suppress(suppressed, newerr error, opts ...option) error {
	return Suppress_skipstack(1, suppressed, newerr, opts...)
}

/* for wrappers of this package */
// with these functions one can specify correct first stack frame to print in stack trace to skip stack frames of
// wrapper objects

// see https://godoc.org/github.com/pashaosipyants/errors/#hdr-Skipstack_management
func New_skipstack(message string, skip int, opts ...option) error {
	return Wrap_skipstack(skip+1, errors.New(message), opts...)
}
// see https://godoc.org/github.com/pashaosipyants/errors/#hdr-Skipstack_management
func Wrap_skipstack(skip int, err error, opts ...option) error {
	if err == nil {
		return nil
	}

	setup := new(setup)
	for _, opt := range opts {
		opt(setup)
	}

	reterr := new(_error)
	if err1, ok := err.(*_error); ok {
		reterr.error = err1.error
		reterr.errcode = err1.errcode
		reterr.stack = err1.stack
		copy_map_annots(err1.annotations, &reterr.annotations)
	}

	if setup.annotation != nil {
		if reterr.annotations == nil {
			reterr.annotations = make(map[string]string)
		}
		reterr.addAnnotation(skip + 1, *setup.annotation)
	}

	if setup.errcode != nil {
		reterr.errcode = *setup.errcode
	}

	if setup.extender != nil {
		reterr.error = setup.extender(reterr.error)
	}

	return reterr
}

// see https://godoc.org/github.com/pashaosipyants/errors/#hdr-Skipstack_management
func Suppress_skipstack(skip int, suppressed, newerr error, opts ...option) error {
	if newerr == nil {
		return nil
	}

	reterr := Wrap_skipstack(skip+1, newerr, opts...).(*_error)
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
