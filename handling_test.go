package errors

import (
	"database/sql"
	"io"
	"testing"
)

func TestCheck(t *testing.T) {
	failures := []func(){
		func() { Check(nil) },
		func() { Check(io.EOF) },
		func() { Check(sql.ErrNoRows, OValue("key", "value")) },
	}
	for i, f := range failures {
		func() {
			defer Handler(func(err error) {
				switch {
				case IsE(err, io.EOF):
				case ValueE(err, "key") == "value":
				default:
					t.Error("can't be here")
				}
			})

			f()
			if i != 0 {
				t.Errorf("can't execute this code. iteration %d", i)
			}
		}()
	}
}

func TestRealPanicHandling(t *testing.T) {
	ok := false
	defer func() {
		if !ok {
			t.Error("real panic haven't been recovered")
		}
	}()
	defer func() {
		if r := recover(); r != nil {
			if r != "AAA" {
				t.Error("wrong recovered")
			} else {
				ok = true
			}
		}
	}()
	defer Handler(func(err error) {
		t.Error("can't be handled here")
	})
	panic("AAA")
}

func TestWithoutAnyPanic(t *testing.T) {
	defer Handler(func(err error) {
		t.Error("can't be handled here")
	})
}

func TestPanicNil(t *testing.T) {
	defer Handler(func(err error) {
		t.Error("can't be handled here")
	})
	panic(nil)
}
