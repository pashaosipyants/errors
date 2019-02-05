package errors

// Check panics with an error, first wrapping it with errcode.
// One can handle this panic defering Handler func.
func Check(err error, errcode ...interface{}) {
	if err != nil {
		panic(Wrap(err, errcode...))
	}
}

// CheckIf panics with an error if ifErr is true, first wrapping err it with errcode.
// One can handle this panic defering Handler func.
func CheckIf(ifErr bool, err error, errcode ...interface{}) {
	if ifErr {
		panic(Wrap(err, errcode...))
	}
}

// CheckIfNew panics with an error if ifErr is true.
// error is made from message and errcode.
// One can handle this panic defering Handler func.
func CheckIfNew(ifErr bool, message string, errcode ...interface{}) {
	if ifErr {
		panic(New(message, errcode...))
	}
}

// Type to handle Check...() produced panic.
type Handleable interface {
	Error() string
	ErrCode() interface{}
	Cause() error
}

// defer Handle and provide handler to process panics made by Check...().
// Panics with type different to this package's one are just forwarded.
func Handler(handle func(err Handleable)) {
	switch r := recover().(type) {
	case nil:
		return
	case *_error:
		handle(r)
	default:
		panic(r)
	}
}
