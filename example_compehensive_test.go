package errors_test

import (
	"github.com/pashaosipyants/errors/example_auxiliary"
	"github.com/pashaosipyants/errors"
	"fmt"
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
	// loop to work out different cases
	for i := 0; i < 4; i++ {
		// func is smth like try block here
		func() {
			defer fmt.Printf("Case %d finished\n-------------------------------------------\n\n\n", i)
			// smth like catch block
			defer errors.Handler(func(err errors.Handleable) {
				switch errors.ErrCode(err) {
				case errcode_apicreatetaskfailed:
					fmt.Printf("%s\n", err) // log
					// may be some specific actions
				case errcode_apiuserloginfailed:
					// may be some specific actions
					panic("Assertion failed") // but in our example can't be here
				default:
					panic("Assertion failed")
				}
			})

			errors.Check(
				apiUserLogin(), errcode_apiuserloginfailed)

			errors.Check(
				apiCreateTask(i), errcode_apicreatetaskfailed) // override errcode

			fmt.Println("Success!!!") // log
		}()
	}

	// Output:
	// wrong output specially to make this function be executed and see output of this example
}

const errcode_apicreatetaskfailed = "api_create_task_failed"
const errcode_apiuserloginfailed = "api_user_login_failed"

func apiCreateTask(i int) (reterr error) {
	defer errors.Handler(func(err errors.Handleable) {
		// do some specific logic - e.g. mark task in db as done
		if errors.ErrCode(err) == example_auxiliary.ErrCode_TaskAlreadyExistAndDone {
			reterr = errors.AnyErr(
				errors.Suppress(err,
					markTaskAsDone(), example_auxiliary.ErrCode_TaskAlreadyExistAndDone),
				err,
			)
		} else {
			reterr = err
		}
		// common logic
		fmt.Print(reterr, "\n") // log
	})

	err := example_auxiliary.CreateTaskInitedByUser1(i)
	err = errors.ExtendCause(
		err,
		func(e error) error { return WrapUserError(e, 1) },
	)
	errors.Check(err)

	return nil
}

// pretends to be always success
func apiUserLogin() error {
	// some work

	return nil
}

// pretends that there is an error an task can not be marked as done
func markTaskAsDone() error {
	return errors.New("task can not be marked as done")
}

func WrapUserError(err error, user int) error {
	if err == nil {
		return nil
	}
	return UserError{
		user: user,
		err: err,
	}
}

type UserError struct {
	user int
	err error
}

func (x UserError) Error() string {
	return x.err.Error()
}

func (x UserError) Format(f fmt.State, verb rune) {
	if verb == 's' {
		fmt.Fprint(f, x.err.Error())
	} else {
		fmt.Fprintf(f, "User: %d; Error: %v", x.user, x.err)
	}
}
