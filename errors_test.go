package errors

import (
	"errors"
	"io"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		err      error
		wantMsg  string
		wantCode interface{}
	}{
		{
			New("msg"), "msg", nil,
		},
		{
			New("msg with format specifiers %v %s"), "msg with format specifiers %v %s", nil,
		},
		{
			New("msg and err code", 1), "msg and err code", 1,
		},
		{
			New("msg and err codes variadic too many", 1, 2, 3),
			"msg and err codes variadic too many", 1,
		},
	}

	for _, tt := range tests {
		if tt.err.Error() != tt.wantMsg {
			t.Errorf("New Error: got: %q, want %q", tt.err, tt.wantMsg)
		}
		if !reflect.DeepEqual(Cause(tt.err), errors.New(tt.wantMsg)) {
			t.Errorf("New Cause: got: %q, want %q", Cause(tt.err), tt.wantMsg)
		}
		if ErrCode(tt.err) != tt.wantCode {
			t.Errorf("New ErrCode: got: %v, want %v", ErrCode(tt.err), tt.wantCode)
		}
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{Errorf("read error without format specifiers"), "read error without format specifiers"},
		{Errorf("read error with %d format specifier", 1), "read error with 1 format specifier"},
	}

	for _, tt := range tests {
		if tt.err.Error() != tt.want {
			t.Errorf("New Error: got: %q, want %q", tt.err, tt.want)
		}
		if !reflect.DeepEqual(Cause(tt.err), errors.New(tt.want)) {
			t.Errorf("New Cause: got: %q, want %q", Cause(tt.err), tt.want)
		}
		if ErrCode(tt.err) != nil {
			t.Errorf("New ErrCode: got: %v, want %v", ErrCode(tt.err), nil)
		}
	}
}

func TestCodef(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{Codef(1, "read error without format specifiers"), "read error without format specifiers"},
		{Codef(1, "read error with %d format specifier", 1), "read error with 1 format specifier"},
	}

	for _, tt := range tests {
		if tt.err.Error() != tt.want {
			t.Errorf("New Error: got: %q, want %q", tt.err, tt.want)
		}
		if !reflect.DeepEqual(Cause(tt.err), errors.New(tt.want)) {
			t.Errorf("New Cause: got: %q, want %q", Cause(tt.err), tt.want)
		}
		if ErrCode(tt.err) != 1 {
			t.Errorf("New ErrCode: got: %v, want %v", ErrCode(tt.err), 1)
		}
	}
}

func TestWrapNil(t *testing.T) {
	got := Wrap(nil, "no error")
	if got != nil {
		t.Errorf("Wrap(nil, \"no error\"): got %v, expected nil", got)
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		err      error
		wantMsg  string
		wantCode interface{}
	}{
		{Wrap(io.EOF), "EOF", nil},
		{Wrap(Wrap(io.EOF)), "EOF", nil},
		{Wrap(io.EOF, 1), "EOF", 1},
		{Wrap(io.EOF, 1, 2, 3), "EOF", 1},
		{Wrap(Wrap(io.EOF, 1)), "EOF", 1},
		{Wrap(Wrap(io.EOF), 1), "EOF", 1},
		{Wrap(Wrap(io.EOF, 2), 1), "EOF", 1},
	}

	for _, tt := range tests {
		if tt.err.Error() != tt.wantMsg {
			t.Errorf("New Error: got: %q, want %q", tt.err, tt.wantMsg)
		}
		if Cause(tt.err) != io.EOF {
			t.Errorf("New Cause: got: %q, want %q", Cause(tt.err), io.EOF)
		}
		if ErrCode(tt.err) != tt.wantCode {
			t.Errorf("New ErrCode: got: %v, want %v", ErrCode(tt.err), tt.wantCode)
		}
	}
}

type nilError struct{}

func (nilError) Error() string { return "nil error" }

func TestCause(t *testing.T) {
	tests := []struct {
		err  error
		want error
	}{{
		// nil error is nil
		err:  nil,
		want: nil,
	}, {
		// explicit nil error is nil
		err:  (error)(nil),
		want: nil,
	}, {
		// typed nil is nil
		err:  (*nilError)(nil),
		want: (*nilError)(nil),
	}, {
		// uncaused error is unaffected
		err:  io.EOF,
		want: io.EOF,
	}, {
		err:  New("AAA"),
		want: errors.New("AAA"),
	}, {
		err:  Wrap(io.EOF),
		want: io.EOF,
	}}

	for i, tt := range tests {
		got := Cause(tt.err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("example_auxiliary %d: got %v, want %v", i+1, got, tt.want)
		}
	}
}

func TestErrCode(t *testing.T) {
	tests := []struct {
		err  error
		wantCode interface{}
	}{{
		err:  nil,
		wantCode: nil,
	}, {
		err:  (error)(nil),
		wantCode: nil,
	}, {
		err:  &_error{errcode: "code"},
		wantCode: "code",
	}}

	for i, tt := range tests {
		if ErrCode(tt.err) != tt.wantCode {
			t.Errorf("example_auxiliary %d: got %v, want %v", i+1, ErrCode(tt.err), tt.wantCode)
		}
	}
}