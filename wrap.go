package errors

import "runtime"

/*
	functions to wrap an error with additional context
*/

// WrapStackE returns error with err wrapped in and stacktrace recorded.
// returnederr.Error() will be the same as err.Error(), but one can use SprintE(returnederr) to print it with stacktrace.
// If err is nil returns nil.
// UnwrapE(returnederr) == err.
//
// skip is optional param. First of variadic parameters is used, else are ignored.
// skip specifies the number of stacktrace levels to skip. By default stacktrace starts with the caller of WrapStackE.
func WrapStackE(err error, skip ...int) error {
	if err == nil {
		return nil
	}

	var errStack *_errorStack
	if AsE(err, &errStack) {
		return err

	} else {
		return &_errorStack{
			error: err,
			stack: callers(getSkip(skip) + 1),
		}
	}
}

// WrapAnnotationE returns error with err wrapped in and annotation(additional message) added.
// returnederr.Error() will be the same as err.Error(), but one can use SprintE(returnederr) to print it with annotations.
// If err is nil returns nil.
// UnwrapE(returnederr) == err.
// The sense of annotation is that by SprintE it's printed along with corresponding stacktrace level.
// So it's obvious in what place in code this annotation was added.
//
// skip is optional param. First of variadic parameters is used, else are ignored.
// skip specifies the number of stacktrace levels to skip. By default stacktrace starts with the caller of WrapAnnotationE.
func WrapAnnotationE(err error, annotation string, skip ...int) error {
	if err == nil {
		return nil
	}

	pc, _, _, _ := runtime.Caller(getSkip(skip) + 1)
	return &_errorAnnotation{
		error:      err,
		where:      runtime.FuncForPC(pc).Name(),
		annotation: annotation,
	}
}

// WrapSuppressedE returns error with err wrapped in and suppressed error added.
// returnederr.Error() will be the same as err.Error(), but one can use SprintE(returnederr) to print it with suppressed error.
// If err is nil returns nil.
// UnwrapE(returnederr) == err.
// List of suppressed errors can be retrieved with the SuppressedE.
func WrapSuppressedE(err, suppressed error) error {
	if err == nil {
		return nil
	}

	return &_errorSuppressed{
		error:      err,
		suppressed: suppressed,
	}
}

// WrapValueE returns error with err wrapped in and value added.
// returnederr.Error() will be the same as err.Error().
// If err is nil returns nil.
// UnwrapE(returnederr) == err.
// Value can be retrieved from returnederr(or any wrappers of it) by specified key with the ValueE function.
func WrapValueE(err error, key, value interface{}) error {
	if err == nil {
		return nil
	}

	var errVal *_errorValue
	for errIteration := err; AsE(errIteration, &errVal); errIteration = errVal.error {
		if errVal.key == key {
			return err // if a value for the given key already exists, do nothing
		}
	}

	return &_errorValue{
		error: err,
		key:   key,
		value: value,
	}
}

func getSkip(skip []int) (skipV int) {
	if len(skip) > 0 {
		skipV = skip[0]
	}
	return skipV
}
