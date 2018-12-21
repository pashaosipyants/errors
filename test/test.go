package test

import "github.com/pashaosipyants/errors"

const (
	ErrCode_ConnectionFailed = "connection_failed"
	ErrCode_TaskAlreadyExistButNotExecuted = "already_exist_but_not_executed"
	ErrCode_TaskAlreadyExistAndExecuted = "already_exist_and_executed"
)

func SaveTaskToDbMockConnectionError() error {
	// db work

	return errors.New("connection failed", ErrCode_ConnectionFailed)
}

func SaveTaskToDbMockExistButNotExecuted() error {
	// db work
	// ask another service if task is already executed - false

	return errors.New("already exist", ErrCode_TaskAlreadyExistButNotExecuted)
}

func SaveTaskToDbMockExistAndExecuted() error {
	// db work
	// ask another service if task is already executed - true

	return errors.New("already exist and executed", ErrCode_TaskAlreadyExistAndExecuted)
}

func SaveTaskToDbMockSuccess() error {
	return nil
}

func CreateTaskInitedByUser1(i int) error {
	// pretends to be service logic function that invokes repository db function

	switch i {
	case 0:
		return errors.WrapAnnotated(SaveTaskToDbMockConnectionError(), "Inited by user 1")
	case 1:
		return errors.WrapAnnotated(SaveTaskToDbMockExistButNotExecuted(), "Inited by user 1")
	case 2:
		return errors.WrapAnnotated(SaveTaskToDbMockExistAndExecuted(), "Inited by user 1")
	case 3:
		return errors.WrapAnnotated(SaveTaskToDbMockSuccess(), "Inited by user 1")
	default:
		panic("Assertion")
	}
}
