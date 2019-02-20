package example_auxiliary

import "github.com/pashaosipyants/errors"

// Error code to distinguish errors.
const (
	ErrCode_ConnectionFailed = "connection_failed"
	ErrCode_TaskAlreadyExistButNotDone = "already_exist_but_not_done"
	ErrCode_TaskAlreadyExistAndDone = "already_exist_and_done"
)

// It mocks case when one tries to save task to db, but connection error occurs.
//
// It returns error with error code ErrCode_ConnectionFailed using:
//     errors.New("connection failed", ErrCode_ConnectionFailed)
func SaveTaskToDbMockConnectionError() error {
	// db work

	return errors.New("connection failed", ErrCode_ConnectionFailed)
}

// It mocks case when one tries to save task to db, but it already exists.
// Then it asks if task is already done and it occurs it is not.
//
// It returns error with error code ErrCode_TaskAlreadyExistButNotDone using:
//     return errors.New("already exist", ErrCode_TaskAlreadyExistButNotDone)
func SaveTaskToDbMockExistButNotDone() error {
	// db work
	// ask another service if task is already done - false

	return errors.New("already exist", ErrCode_TaskAlreadyExistButNotDone)
}

// It mocks case when one tries to save task to db, but it already exists.
// Then it asks if task is already done and it occurs it is done.
//
// It returns error with error code ErrCode_TaskAlreadyExistAndDone using:
//     return errors.New("already exist and done", ErrCode_TaskAlreadyExistAndDone)
func SaveTaskToDbMockExistAndDone() error {
	// db work
	// ask another service if task is already done - true

	return errors.New("already exist and done", ErrCode_TaskAlreadyExistAndDone)
}

// It mocks case of successful saving task.
func SaveTaskToDbMockSuccess() error {
	return nil
}

// Pretends to be service logic function that invokes repository db function SaveTaskToDb...().
//
// 0 - calls SaveTaskToDbMockConnectionError
// 1 - calls SaveTaskToDbMockExistButNotDone
// 2 - calls SaveTaskToDbMockExistAndDone
// 3 - calls SaveTaskToDbMockSuccess
//
// All errors occured it annotates with information that this task id inited by user 1 using:
//     errors.WrapAnnotated(SaveTaskToDb...(), "Inited by user 1")
func CreateTaskInitedByUser1(i int) error {
	switch i {
	case 0:
		return errors.WrapAnnotated(
			SaveTaskToDbMockConnectionError(), "Inited by user 1")
	case 1:
		return errors.WrapAnnotated(
			SaveTaskToDbMockExistButNotDone(), "Inited by user 1")
	case 2:
		return errors.WrapAnnotated(
			SaveTaskToDbMockExistAndDone(), "Inited by user 1")
	case 3:
		return errors.WrapAnnotated(
			SaveTaskToDbMockSuccess(), "Inited by user 1")
	default:
		panic("Assertion")
	}
}
