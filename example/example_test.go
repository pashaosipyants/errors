package error_test

import (
	"github.com/pashaosipyants/errors/test"
	"github.com/pashaosipyants/errors"
	"fmt"
)

// imagine there is an api to create a task that execute another service
func ExamplePackage() {
	for i := 0; i < 4; i++ {
		func() {
			defer errors.Handler(func(err errors.Handleable) {
				switch errors.ErrCode(err) {
				case errcode_apicreatetaskfailed:
					fmt.Printf("%s\n\n\n", err) // log
					// may be some specific actions
				case errcode_apiuserloginfailed:
					// may be some specific actions
					panic("Assertion failed") // but in our case can't be here
				default:
					panic("Assertion failed")
				}
			})

			err := apiUserLogin()
			errors.Check(err, errcode_apiuserloginfailed)

			err = apiCreateTask(i)
			errors.Check(err, errcode_apicreatetaskfailed) // override errcode
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
		case test.ErrCode_ConnectionFailed, test.ErrCode_TaskAlreadyExistButNotExecuted:
			fmt.Print(err, "\n\n\n\n") // log
			reterr = err
		case test.ErrCode_TaskAlreadyExistAndExecuted:
			fmt.Print(err, "\n\n\n\n") // log
			// do some specific logic - e.g. mark task in db as executed
			reterr = err
		default:
			panic("Assertion failed")
		}
	})

	err := test.CreateTaskInitedByUser1(i)
	errors.Check(err)

	return nil
}

// pretends to be always success
func apiUserLogin() error {
	// some work

	return nil
}
