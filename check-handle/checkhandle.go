package check_handle

// Check panics with an error, first wrapping it with errcode.
// One can handle this panic defering Handler func.
func Check(err error) {
	if err != nil {
		panic(Handleable{err})
	}
}

// CheckIf panics with an error if ifErr is true, first wrapping err it with errcode.
// One can handle this panic defering Handler func.
func CheckIf(ifErr bool, err error) {
	if ifErr {
		panic(Handleable{err})
	}
}

// Type to handle Check...() produced panic.
// not to confuse with ordinary panic
type Handleable struct {
	err error
}

// defer Handle and provide handler to process panics made by Check...().
// Panics with type different to this package's one are just forwarded.
// elseDefer - ...
// other panics are not intercepted
// don't forget that Handler may not obtain error if u recover earlier, check type
func Handler(handle func(err error), elseDefer ...func()) {
	switch r := recover().(type) {
	case nil:
		for _, f := range elseDefer {
			f()
		}
	case Handleable:
		handle(r.err)
	default:
		panic(r)
	}
}
