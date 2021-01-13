package errors

/*
	common function to wrap an error with several additional contexts right away
*/

type OptionE func(error, int) error

// WrapE - one function that does anything other Wrap... functions can.
// Instead of doing smth like WrapStackE(WrapAnnotationE(WrapValueE(err, key, value), "text")) just
// use several options in one WrapE. Like this: WrapE(err, OStack(), OAnno(text), OValue(key, value)).
func WrapE(err error, opts ...OptionE) error {
	return ToSkipE(1).WrapE(err, opts...)
}

func OStack() OptionE {
	return func(err error, skip int) error {
		return WrapStackE(err, skip)
	}
}

func OAnno(annotation string) OptionE {
	return func(err error, skip int) error {
		return WrapAnnotationE(err, annotation, skip)
	}
}

func OSupp(suppressed error) OptionE {
	return func(err error, _ int) error {
		return WrapSuppressedE(err, suppressed)
	}
}

func OValue(key, value interface{}) OptionE {
	return func(err error, _ int) error {
		return WrapValueE(err, key, value)
	}
}

type ToSkipE int

// The same as WrapE, but skip can be specified to skip several stacktrace levels.
// This is useful when WrapE is used inside some utility function which will be useless in stacktrace.
// E.g. see implementation of Error in this package. It utilizes WrapE, but skips itself from the stacktrace.
func (skip ToSkipE) WrapE(err error, opts ...OptionE) error {
	for _, opt := range opts {
		err = opt(err, int(skip)+2)
	}
	return err
}
