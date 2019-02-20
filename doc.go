/*
Package errors provides detailed stack-traced description for error type and exception(handle-check) style handling

This package has two goals.

1. Provide convenient error type which has error code to easily distinguish errors,
which can be extended by additional messages(annotations here) while error flows through the program,
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

With WrapAnnotated one can add additional messages to an error. Annotations are printed in stacktrace, bound
to corresponding function(where it was added).
If annotation was added out of functions in stacktrace, it will be printed in separate section.

Format

Errors printout consists of 3 parts: ERROR, STACK, SUPPRESSED.
ERROR part:

Prints underlying error. Uses Format if underlying error implements fmt.Formatter.

STACK part:

Prints stack trace with annotations if modifier is %v.

SUPPRESSED part:
Prints suppressed error if it exists. Uses Format if it implements fmt.Formatter.

E.g. with %v:

    ERROR:
    connection failed

    STACK:
    github.com/pashaosipyants/errors/example_auxiliary.SaveTaskToDbMockConnectionError
    	D:/work/go/src/github.com/pashaosipyants/errors/example_auxiliary/example_auxiliary.go:19
    github.com/pashaosipyants/errors/example_auxiliary.CreateTaskInitedByUser1
    	D:/work/go/src/github.com/pashaosipyants/errors/example_auxiliary/example_auxiliary.go:63
    ANNOTATIONS:
    Inited by user 1
    github.com/pashaosipyants/errors_test.apiCreateTask
    	D:/work/go/src/github.com/pashaosipyants/errors/example_compehensive_test.go:70
    github.com/pashaosipyants/errors_test.Example.func1
    	D:/work/go/src/github.com/pashaosipyants/errors/example_compehensive_test.go:41
    github.com/pashaosipyants/errors_test.Example
    	D:/work/go/src/github.com/pashaosipyants/errors/example_compehensive_test.go:45
    testing.runExample
    	C:/Go/src/testing/example.go:121
    testing.runExamples
    	C:/Go/src/testing/example.go:45
    testing.(*M).Run
    	C:/Go/src/testing/testing.go:1035
    main.main
    	_testmain.go:70
    runtime.main
    	C:/Go/src/runtime/proc.go:201
    runtime.goexit
    	C:/Go/src/runtime/asm_amd64.s:1333

Cause

To get underlying error use Cause func.

    err := errors.Wrap(io.EOF)
    if reflect.DeepEqual(errors.Cause(err), io.EOF) {
      fmt.Println("Hooray, it works!")
    }

Skipstack management

All creation & wrapping functions have their duplicates with _skipstack suffix and skip argument.
With those functions one can specify correct first stack frame to print in stack trace.
skip==0 means starting from caller of this function.
It can be useful for example to skip stack frames of wrapper object, if you want your own wrapper of this package.

    type MyErr struct {
      error
      X AddInfo
    }

    func NewMyErr(message string, errcode ...interface{}, x AddInfo) MyErr {
      return MyErr{New_skipstack(message, 1, errcode...), x}
    }


Handle check mechanism

example:

    defer errors.Handler(func(err errors.Handleable) {
      log.Error(err)
      switch err.ErrCode() {
        case 1:
          repairError1()
        case 2:
          repairError2()
        default:
          panic("Unknown err code")
      }
    })

    err := task1()
    errors.Check(err, 1)

    x := twoPlusTwo()
    errors.CheckIfNew(x != 4, "twoPlusTwo is wrong", 2)
*/
package errors
