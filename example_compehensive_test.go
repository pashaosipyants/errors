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
// Imagine there is an api to create a task that executes another service.
// Besides user of this api wants to hold info whether this task is already done.
// This api, ofc, can return error. E.g. certain task may be already created.
// If so, error should report whether it is done or not.
func Example() {
	// loop to work out different cases
	for i := 0; i < 4; i++ {
		// func is smth like try block here
		func() {
			// smth like catch block
			defer errors.Handler(func(err errors.Handleable) {
				switch errors.ErrCode(err) {
				case errcode_apicreatetaskfailed:
					fmt.Printf("%s\n\n\n", err) // log
					// may be some specific actions
				case errcode_apiuserloginfailed:
					// may be some specific actions
					panic("Assertion failed") // but in our example can't be here
				default:
					panic("Assertion failed")
				}
			})

			err := apiUserLogin()
			errors.Check(err, errcode_apiuserloginfailed)

			err = apiCreateTask(i)
			errors.Check(err, errcode_apicreatetaskfailed) // override errcode

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
		switch errors.ErrCode(err) {
		case example_auxiliary.ErrCode_ConnectionFailed, example_auxiliary.ErrCode_TaskAlreadyExistButNotDone:
			fmt.Print(err, "\n\n\n\n") // log
			reterr = err
		case example_auxiliary.ErrCode_TaskAlreadyExistAndDone:
			fmt.Print(err, "\n\n\n\n") // log
			// do some specific logic - e.g. mark task in db as done
			reterr = err
		default:
			panic("Assertion failed")
		}
	})

	err := example_auxiliary.CreateTaskInitedByUser1(i)
	errors.Check(err)

	return nil
}

// pretends to be always success
func apiUserLogin() error {
	// some work

	return nil
}
