package errors

import "testing"

func TestCheck(t *testing.T) {
	for i := 0; i < 3; i++ {
		func() {
			defer Handler(func(err Handleable) {
				switch x := ErrCode(err); x {
				case 0:
				case 1:
				case 2:
				default:
					t.Errorf("wrong error code %v", x)
				}
			})

			Check(New("AAA", 100), i)
			t.Error("can't execute this code")
		}()
	}
}

func twoPlusTwo() int {
	return 5
}

func TestCheckIf(t *testing.T) {
	for i := 0; i < 3; i++ {
		func() {
			defer Handler(func(err Handleable) {
				switch x := ErrCode(err); x {
				case 0:
				case 1:
				case 2:
				default:
					t.Errorf("wrong error code %v", x)
				}
			})

			CheckIf(twoPlusTwo() != 4, New("twoPlusTwo is wrong", 100), i)
			t.Error("can't execute this code")
		}()
	}
}

func TestCheckIfNew(t *testing.T) {
	for i := 0; i < 3; i++ {
		func() {
			defer Handler(func(err Handleable) {
				switch x := ErrCode(err); x {
				case 0:
				case 1:
				case 2:
				default:
					t.Errorf("wrong error code %v", x)
				}
			})

			CheckIfNew(twoPlusTwo() != 4, "twoPlusTwo is wrong", i)
			t.Error("can't execute this code")
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
	defer Handler(func(err Handleable) {
		t.Error("can't be handled here")
	})
	panic("AAA")
}

func TestWithoutAnyPanic(t *testing.T) {
	defer Handler(func(err Handleable) {
		t.Error("can't be handled here")
	})
}

func TestPanicNil(t *testing.T) {
	defer Handler(func(err Handleable) {
		t.Error("can't be handled here")
	})
	panic(nil)
}
