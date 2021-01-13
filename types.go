package errors

/*
	error types which allow to add context to an error:
		- stack - stacktrace
		- annotation - additional msg
		- suppressed - error which is suppressed by newer one
		- value - additional value that can be accessed by corresponding key any time later
*/

type _errorStack struct {
	error
	stack
}

type _errorAnnotation struct {
	error
	where, annotation string // where defines at which level of stacktrace this annotation was added
}

type _errorSuppressed struct {
	error
	suppressed error
}

type _errorValue struct {
	error
	key, value interface{}
}

func (e *_errorStack) Unwrap() error {
	return e.error
}

func (e *_errorAnnotation) Unwrap() error {
	return e.error
}

func (e *_errorSuppressed) Unwrap() error {
	return e.error
}

func (e *_errorValue) Unwrap() error {
	return e.error
}
