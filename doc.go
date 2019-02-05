/*
Package errors provides detailed stack-traced description for error type and exception(handle-check) style handling

This package has two goals.

1. Provide convenient error type which has error code to easily distinguish errors,
can be extended by additional messages(annotations here) while error flows through the program,
can be printed with stack trace.

2. Provide convenient way to handle errors in exception style.

With the first point everything is clear, let's talk a little about second.
Go's error handling idiom usually looks like

    x, err := DoSmth()
    if err != nil {
        logger.Error(err)
        // may be some err processing
        return err
    }

It is very "holy war" question which way to handle errors, exception like or "returning err" like, is better.
Everyone chooses himself which way is more appropriate to him.
But doubtlessly "returning err" style has some disadvantages.

1. Error handling logic is mixed with "business" logic so it shadows it

2. Code repetitions, which can cause a bug

3. Verbosity

That's why Go 2 draft https://go.googlesource.com/proposal/+/master/design/go2draft.md contains new Handle-Check
error handling approach, which is quite similar to exception approach.
This package provides similar to this Handle-Check mechanism, based on panics.
Details are below =)

Errors

All errors created by this package are returned as golang native error interface.
Underlying type is immutable, so one can't change smth in already existing error, but
can use wrap functions, which create new one based on underlying.
New, Errorf, Codef - creates new error. Wrap, WrapAnnotated - wrap existing.

Error codes

Error can be created with error code, which has interface{} type. Error code can be obtained from error by ErrCode
function.

    err := errors.New("description", 1)
    if errors.ErrCode(err) == 1 {
      fmt.Println("Hooray, it works!")
    }

Missing error code in creation is equal to nil error code.

    err := errors.New("description")
    if errors.ErrCode(err) == nil {
      fmt.Println("Hooray, it works!")
    }

When one wraps error without providing error code, it is preserved. Otherwise - overridden.

    err := errors.New("description", 1)
    err = errors.Wrap(err)
    if errors.ErrCode(err) == 1 {
      fmt.Println("Hooray, it works!")
    }
    err = errors.Wrap(err, 2)
    if errors.ErrCode(err) == 2 {
      fmt.Println("Hooray, it still works!")
    }

Annotations



Format



Skipstack management
// with these functions one can specify correct first stack frame to print in stack trace to skip stack frames of
// wrapper objects
skip==0 это что

Handle-Check mechanism

not thread safe

*/
package errors
