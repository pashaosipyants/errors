package errors

/*
	helper functions to create new error with several additional contexts right away
*/

import "fmt"

// Creates new error. Use options to fill error's context.
func Error(text string, opts ...OptionE) error {
	return ToSkipE(1).WrapE(NewE(text), opts...)
}

// Creates new error.
// E.g.: Errorf("user with id %s is absent", userID).Opts(OStack(), OAnno(text))
func Errorf(format string, a ...interface{}) errf {
	return errf{
		err: fmt.Errorf(format, a...),
	}
}

func (e errf) Opts(opts ...OptionE) error {
	return ToSkipE(1).WrapE(e.err, opts...)
}

type errf struct {
	err error
}
