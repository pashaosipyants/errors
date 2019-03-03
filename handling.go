package errors

// Check panics with an error, first wrapping it with errcode.
// One can handle this panic defering Handler func.
func Check(err error, opts ...option) {
	if err != nil {
		panic(Wrap_skipstack(1, err, opts...))
	}
}

// CheckIf panics with an error if ifErr is true, first wrapping err it with errcode.
// One can handle this panic defering Handler func.
func CheckIf(ifErr bool, err error, opts ...option) {
	if ifErr {
		panic(Wrap_skipstack(1, err, opts...))
	}
}

// CheckIfNew panics with an error if ifErr is true.
// error is made from message and errcode.
// One can handle this panic defering Handler func.
func CheckIfNew(ifErr bool, message string, opts ...option) {
	if ifErr {
		panic(New_skipstack(message, 1, opts...))
	}
}

// Type to handle Check...() produced panic.
type Handleable interface {
	Error() string
	ErrCode() interface{}
	Cause() error
	Suppressed() error
}

// defer Handle and provide handler to process panics made by Check...().
// Panics with type different to this package's one are just forwarded.
// elseDefer - ...
func Handler(handle func(err Handleable), elseDefer ...func()) {
	switch r := recover().(type) {
	case nil:
		for _, f := range elseDefer {
			f()
		}
	case *_error:
		handle(r)
	default:
		panic(r)
	}
}

/* for wrappers of this package */
// with these functions one can specify correct first stack frame to print in stack trace to skip stack frames of
// wrapper objects

func Check_skipstack(skip int, err error, opts ...option) {
	if err != nil {
		panic(Wrap_skipstack(skip+1, err, opts...))
	}
}

func CheckIf_skipstack(skip int, ifErr bool, err error, opts ...option) {
	if ifErr {
		panic(Wrap_skipstack(skip+1, err, opts...))
	}
}

func CheckIfNew_skipstack(skip int, ifErr bool, message string, opts ...option) {
	if ifErr {
		panic(New_skipstack(message, skip+1, opts...))
	}
}
