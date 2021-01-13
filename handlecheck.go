package errors

// Check panics with an error.
// One can handle this panic defering Handler func.
// opts allows to add context to the error. OStack is enabled by default.
func Check(err error, opts ...OptionE) {
	if err != nil {
		opts = append(opts, OStack())
		panic(handleable{ToSkipE(1).WrapE(err, opts...)})
	}
}

// CheckIf panics with an error if ifErr is true.
// One can handle this panic defering Handler func.
// opts allows to add context to the error. OStack is enabled by default.
func CheckIf(ifErr bool, err error, opts ...OptionE) {
	if ifErr {
		opts = append(opts, OStack())
		panic(handleable{ToSkipE(1).WrapE(err, opts...)})
	}
}

// The same as Check, but OStack is not enabled by default.
func CheckNoStack(err error, opts ...OptionE) {
	if err != nil {
		panic(handleable{ToSkipE(1).WrapE(err, opts...)})
	}
}

// The same as CheckIf, but OStack is not enabled by default.
func CheckIfNoStack(ifErr bool, err error, opts ...OptionE) {
	if ifErr {
		panic(handleable{ToSkipE(1).WrapE(err, opts...)})
	}
}

// Type to handle Check...() produced panic.
// not to confuse with ordinary panic
type handleable struct {
	err error
}

// defer Handler func and provide handler to process panics made by Check...().
// Panics with type different to this package's one are just forwarded.
// elseDefer - deferred functions which will be executed if no error arisen.
// other panics are not intercepted.
// don't forget that Handler may not obtain error if u recover earlier.
func Handler(handle func(err error), elseDefer ...func()) {
	switch r := recover().(type) {
	case nil:
		for _, f := range elseDefer {
			f()
		}
	case handleable:
		handle(r.err)
	default:
		panic(r)
	}
}

// as Handler, but instead of running your handle function put arisen error in reterr.
func DefaultHandler(reterr *error, elseDefer ...func()) {
	switch r := recover().(type) {
	case nil:
		for _, f := range elseDefer {
			f()
		}
	case handleable:
		*reterr = r.err
	default:
		panic(r)
	}
}
