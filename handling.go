package errors

import c_h "github.com/pashaosipyants/errors/check-handle"

// Check panics with an error, first wrapping it with errcode.
// One can handle this panic defering Handler func.
func Check(err error, opts ...Option) {
	c_h.Check(
		Wrap(err, append(opts, SkipAdd(1))...),
		)
}

// CheckIf panics with an error if ifErr is true, first wrapping err it with errcode.
// One can handle this panic defering Handler func.
func CheckIf(ifErr bool, err error, opts ...Option) {
	c_h.CheckIf(ifErr,
		Wrap(err, append(opts, SkipAdd(1))...),
		)
}

// CheckIfNew panics with an error if ifErr is true.
// error is made from message and errcode.
// One can handle this panic defering Handler func.
func CheckIfNew(ifErr bool, message string, opts ...Option) {
	c_h.CheckIf(ifErr,
		New(message, append(opts, SkipAdd(1))...),
		)
}

func CheckIfNewf(ifErr bool, format string, args []interface{}, opts ...Option) {
	c_h.CheckIf(ifErr,
		Newf(format, args, append(opts, SkipAdd(1))...),
	)
}

type Handleable = c_h.Handleable

// defer Handle and provide handler to process panics made by Check...().
// Panics with type different to this package's one are just forwarded.
// elseDefer - ...
func Handler(handle func(err error), elseDefer ...func()) {
	c_h.Handler(handle, elseDefer...)
}
