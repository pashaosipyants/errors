package errors

// function panics with an error, first wrapping it with errcode
// one can handle this panic defering Handler func
func Check(err error, errcode ...interface{}) {
	if err != nil {
		panic(Wrap(err, errcode...))
	}
}

// function panics with an error if ifErr is true, first wrapping err it with errcode
// one can handle this panic defering Handler func
func CheckIf(ifErr bool, err error, errcode ...interface{}) {
	if ifErr {
		panic(Wrap(err, errcode...))
	}
}

// function panics with an error if ifErr is true. error is made from message and errcode
// one can handle this panic defering Handler func
func CheckIfNew(ifErr bool, message string, errcode ...interface{}) {
	if ifErr {
		panic(New(message, errcode...))
	}
}

// type to handle on Check...() produced panic
type Handleable interface {
	Error() string
	ErrCode() interface{}
	Cause() error
}

// defer Handle and provide handler to process panics made by Check...()
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
