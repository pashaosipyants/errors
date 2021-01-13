package errors_test

import (
	"bytes"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	. "github.com/pashaosipyants/errors/v2"
	"github.com/pashaosipyants/errors/v2/example_auxiliary"
)

// This is comprehensive, pretending to be close to real-life example of using this package.
// It's easier to see it in code, but if you use godoc, please, notice
// https://godoc.org/github.com/pashaosipyants/errors/example_auxiliary package,
// which is used here.
//
// Imagine there is an api to create a task that is executed by another service.
// Besides user of this api wants to hold info whether this task is already done.
// This api, ofc, can return error. E.g. certain task may be already created.
// If so, error should report whether it is done or not.
func Example() {
	l := logrus.New().WithField("application", "tasks")
	l.Logger.SetFormatter(TerminalVerbose{})
	l.Logger.SetOutput(os.Stdout)

	// loop to work out different cases
	for i := 0; i < 4; i++ {
		// func is smth like try block here
		func() {
			defer l.Infof("\n\tCase %d finished\n-------------------------------------------\n\n\n", i)

			// smth like catch block
			defer Handler(func(err error) {
				switch ValueE(err, "api") {
				case errcode_apicreatetaskfailed:
					logger, ok := ValueE(err, "logger").(*logrus.Entry) // logger with relevant fields of functions deeper in the call stack
					if !ok {
						logger = l
					}
					var ue *UserError
					if AsE(err, &ue) {
						logger = logger.WithField("user", ue.user)
					}
					logger.Error(SprintE(err)) // log
					// may be some specific actions
				case errcode_apiuserloginfailed:
					// may be some specific actions
					panic("Assertion failed") // but in our example can't be here
				default:
					panic("Assertion failed")
				}
			})

			Check(
				apiUserLogin(l), OValue("api", errcode_apiuserloginfailed))

			Check(
				apiCreateTask(l, i), OValue("api", errcode_apicreatetaskfailed))

			l.Info("Success!!!\n") // log
		}()
	}

	// Output:
	// wrong output specially to make this function be executed and see output of this example
}

const errcode_apicreatetaskfailed = "api_create_task_failed"
const errcode_apiuserloginfailed = "api_user_login_failed"

func apiCreateTask(l *logrus.Entry, i int) (reterr error) {
	defer Handler(func(err error) {
		// do some specific logic - e.g. mark task in db as done
		switch {
		case IsE(err, example_auxiliary.ErrTaskAlreadyExistButNotDone):
			errOnMark := markTaskAsDone()
			reterr = AnyE(
				WrapE(errOnMark, OSupp(err)),
				err,
			)
		default:
			reterr = err
		}
		// common logic
	})

	err := example_auxiliary.CreateTaskInitedByUser(l, i, 239)
	Check(WrapUserError(err, 239))

	return nil
}

// pretends to be always success
func apiUserLogin(l *logrus.Entry) error {
	// some work

	return nil
}

// pretends that there is an error an task can not be marked as done
func markTaskAsDone() (reterr error) {
	defer DefaultHandler(&reterr)
	CheckIf(true, Error("task can not be marked as done"))
	return
}

func WrapUserError(err error, user int) error {
	if err == nil {
		return nil
	}
	return &UserError{
		user: user,
		err:  err,
	}
}

type UserError struct {
	user int
	err  error
}

func (x *UserError) Error() string {
	return x.err.Error()
}

func (x *UserError) Unwrap() error {
	return x.err
}

func (x *UserError) Format(f fmt.State, verb rune) {
	if verb == 's' {
		fmt.Fprint(f, x.err.Error())
	} else {
		fmt.Fprintf(f, "User: %d; Error: %v", x.user, x.err)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////

type TerminalVerbose struct {
}

const breaker = "........................................................\n"

func (t TerminalVerbose) Format(e *logrus.Entry) ([]byte, error) {
	if e.Level == logrus.InfoLevel {
		msg := lineBreaker(e.Message)
		return []byte(msg), nil
	}

	msg := lineBreaker(e.Message)

	var fields string
	for k, v := range e.Data {
		msg := fmt.Sprint(v)
		msg = lineBreaker(msg)
		fields += "--" + k + ":\n" + msg
	}

	var b *bytes.Buffer
	if e.Buffer != nil {
		b = e.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	_, err := b.WriteString("\n" + breaker + msg + "\n" + fields)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func lineBreaker(in string) string {
	if len(in) > 0 && in[len(in)-1] != '\n' {
		in += "\n"
	}
	return in
}
