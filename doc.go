/*
Package errors follows 2 goals:

1. Provide a simple way to enrich an error with some context details, using go13 Wrap/Unwrap mechanism.
Such types of context are:
- stacktrace
- annotation - additional msg
- error, which is suppressed by newer one
- any value that can be accessed by corresponding key any time later

2. Provide a convenient way to handle errors in exception style.

About first point:

Error might look like a context vice versa. It also flows through the function call stack, but, opposite to context,
it flows upwards. While it flows we want to enrich it with some information relevant to current stack level.
This package provides several wrapping constructs providing this ability.

One can wrap an existing error with an error with a stacktrace, so that this error will be possible to print
with it's stacktrace and easily find a place error occured in. Only the most deep stacktrace is saved and printed.

One can add an annotation to an error. Moreover, if error already has a stacktrace, this annotation will be printed
along with the corresponding line in stacktrace, so it's easy to find in which function this annotations was added.

One can suppress an existing error with newer one. Both errors will be printed.

One can add a value to an error. Later this value can be retrieved with the key used to add this value. If several
values are added with the same key, only the first one (from the most deep stack level) is saved. It might be
very useful to pass logger to the function upper in the stack. Look at the example to see this approach.

About second point:

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

Example of error format, when printed with SprintE:

    ERROR:
    task can not be marked as done

    STACK:
    github.com/pashaosipyants/errors/v2_test.markTaskAsDone
    	/work/go/src/github.com/pashaosipyants/errors/example_compehensive_test.go:105
    github.com/pashaosipyants/errors/v2_test.apiCreateTask.func1
    	/work/go/src/github.com/pashaosipyants/errors/example_compehensive_test.go:78
    github.com/pashaosipyants/errors/v2.Handler
    	/work/go/src/github.com/pashaosipyants/errors/handlecheck.go:55
    github.com/pashaosipyants/errors/v2_test.apiCreateTask
    	/work/go/src/github.com/pashaosipyants/errors/example_compehensive_test.go:90
    github.com/pashaosipyants/errors/v2_test.Example.func1
    	work/go/src/github.com/pashaosipyants/errors/example_compehensive_test.go:60
    github.com/pashaosipyants/errors/v2_test.Example
    	/work/go/src/github.com/pashaosipyants/errors/example_compehensive_test.go:63
    testing.runExample
    	/Go/src/testing/run_example.go:62
    testing.runExamples
    	/Go/src/testing/example.go:44
    testing.(*M).Run
    	/Go/src/testing/testing.go:1200
    main.main
    	_testmain.go:58
    runtime.main
    	/Go/src/runtime/proc.go:203
    runtime.goexit
    	/Go/src/runtime/asm_amd64.s:1373

    SUPPRESSED:
    ERROR:
    User: 239; Error: already_exist_but_not_done

    STACK:
    github.com/pashaosipyants/errors/v2/example_auxiliary.SaveTaskToDbMockExistButNotDone
    	/work/go/src/github.com/pashaosipyants/errors/example_auxiliary/example_auxiliary.go:36
    github.com/pashaosipyants/errors/v2/example_auxiliary.CreateTaskInitedByUser
    	/work/go/src/github.com/pashaosipyants/errors/example_auxiliary/example_auxiliary.go:78
    	ANNOTATION: Inited by user 239
    github.com/pashaosipyants/errors/v2_test.apiCreateTask
    	/work/go/src/github.com/pashaosipyants/errors/example_compehensive_test.go:89
    github.com/pashaosipyants/errors/v2_test.Example.func1
    	/work/go/src/github.com/pashaosipyants/errors/example_compehensive_test.go:60
    github.com/pashaosipyants/errors/v2_test.Example
    	/work/go/src/github.com/pashaosipyants/errors/example_compehensive_test.go:63
    testing.runExample
    	/Go/src/testing/run_example.go:62
    testing.runExamples
    	/Go/src/testing/example.go:44
    testing.(*M).Run
    	/Go/src/testing/testing.go:1200
    main.main
    	_testmain.go:58
    runtime.main
    	/Go/src/runtime/proc.go:203
    runtime.goexit
    	/Go/src/runtime/asm_amd64.s:1373
*/
package errors
