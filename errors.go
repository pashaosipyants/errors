package errors

import (
	"fmt"
)

// immutable type. none of functions changes it
type _error struct {
	error
	annotatedStack
	suppressed error
}

// New returns error with error code, if specified.
//
// errcode is optional param. First of variadic parameters is used, else are ignored.
func New(message string, opts ...Option) error {
	opts = append([]Option{Skip(1)}, opts...)
	return Wrap(Error(message), opts...)
}

func Newf(format string, args []interface{}, opts ...Option) error {
	opts = append([]Option{Skip(1)}, opts...)
	return Wrap(Errorf(format, args...), opts...)
}

// Wrap returns error with underlying err and error code, if specified.
// If err is nil returns nil.
// If err already has error code(it means it was created earlier by one of this package's functions) and
// errcode is specified, it is overridden.
//
// errcode is optional param. First of variadic parameters is used, else are ignored.
func Wrap(err error, opts ...Option) error {
	if err == nil {
		return nil
	}

	setup := new(setup).setup(opts)

	reterr := new(_error)
	if err1, ok := err.(*_error); ok {
		reterr.error = err1.error
		reterr.suppressed = err1.suppressed
		reterr.stack = err1.stack
		copy_map_annots(err1.annotations, &reterr.annotations)
	} else {
		reterr.error = err
		reterr.stack = callers(setup.skip + 1)
	}

	if setup.annotation != nil {
		if reterr.annotations == nil {
			reterr.annotations = make(map[string]string)
		}
		reterr.addAnnotation(setup.skip + 1, *setup.annotation)
	}

	if setup.extender != nil {
		reterr.error = setup.extender(reterr.error)
	}

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
func Suppress(suppressed, newerr error, opts ...Option) error {
	if newerr == nil {
		return nil
	}

	opts = append(opts, SkipAdd(1))
	reterr := Wrap(newerr, opts...).(*_error)
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

func (w *_error) Handleable82aad239749f494ea1b6459518249fef() {}

/* getters */

func Is(err1, err2 error) bool {
	return Cause(err1) == Cause(err2)
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
