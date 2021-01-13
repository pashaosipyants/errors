package example_auxiliary

import (
	"fmt"

	. "github.com/pashaosipyants/errors/v2"
	"github.com/sirupsen/logrus"
)

// Different errors
var (
	ErrConnectionFailed           = NewE("connection_failed")
	ErrTaskAlreadyExistButNotDone = NewE("already_exist_but_not_done")
	ErrTaskAlreadyExistAndDone    = NewE("already_exist_and_done")
)

// It mocks case when one tries to save task to db, but connection error occurs.
//
// It returns specific error with stacktrace and with logger inside, which contains relevant fields.
func SaveTaskToDbMockConnectionError(l *logrus.Entry) error {
	// db work

	l = l.WithField("reason", "connection")
	return WrapE(ErrConnectionFailed, OStack(), OValue("logger", l))
}

// It mocks case when one tries to save task to db, but it already exists.
// Then it asks if task is already done and it occurs it is not.
//
// It returns specific error with stacktrace and with logger inside, which contains relevant fields.
func SaveTaskToDbMockExistButNotDone(l *logrus.Entry) error {
	// db work
	// ask another service if task is already done - false

	l = l.WithField("reason", "duplicate/notDone")
	return WrapE(ErrTaskAlreadyExistButNotDone, OStack(), OValue("logger", l))
}

// It mocks case when one tries to save task to db, but it already exists.
// Then it asks if task is already done and it occurs it is done.
//
// It returns specific error with stacktrace and with logger inside, which contains relevant fields.
func SaveTaskToDbMockExistAndDone(l *logrus.Entry) error {
	// db work
	// ask another service if task is already done - true

	l = l.WithField("reason", "duplicate/Done")
	return WrapE(ErrTaskAlreadyExistAndDone, OStack(), OValue("logger", l))
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
// All errors occured it annotates with information that this task is inited by userID
func CreateTaskInitedByUser(l *logrus.Entry, i, userID int) (reterr error) {
	l = l.WithField("service", "database")

	defer Handler(func(err error) {
		reterr = WrapE(err, OValue("logger", l)) // actually not set, because it's set inside further functions with more specific fields
	})

	switch i {
	case 0:
		Check(
			SaveTaskToDbMockConnectionError(l),
			OAnno(fmt.Sprintf("Inited by user %v", userID)))
	case 1:
		Check(
			SaveTaskToDbMockExistButNotDone(l),
			OAnno(fmt.Sprintf("Inited by user %v", userID)))
	case 2:
		Check(
			SaveTaskToDbMockExistAndDone(l),
			OAnno(fmt.Sprintf("Inited by user %v", userID)))
	case 3:
		Check(
			SaveTaskToDbMockSuccess(),
			OAnno(fmt.Sprintf("Inited by user %v", userID)))
	default:
		panic("Assertion")
	}
	return
}
