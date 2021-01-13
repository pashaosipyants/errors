package errors

import (
	"database/sql"
	"fmt"
	"io"
	"testing"
)

func TestWrapNil(t *testing.T) {
	checkTest := func(funcTest string, got error) {
		if got != nil {
			t.Errorf("%s: got %v, expected nil", funcTest, got)
		}
	}

	checkTest("WrapStackE(nil)",
		WrapStackE(nil))
	checkTest("WrapAnnotationE(nil)",
		WrapAnnotationE(nil, ""))
	checkTest("WrapSuppressedE(nil)",
		WrapSuppressedE(nil, io.EOF))
	checkTest("WrapValueE(nil)",
		WrapValueE(nil, "", ""))
}

func errofWrap() OptionE {
	return func(err error, _ int) error {
		return fmt.Errorf("wrapped: %w", err)
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		err            error
		opts           []OptionE
		key, value     interface{}
		wantMsg        string
		wantSuppressed error
	}{
		{io.EOF, nil, nil, nil, "EOF", nil},
		{io.EOF,
			[]OptionE{OStack(), OSupp(sql.ErrNoRows)},
			nil, nil, "EOF", sql.ErrNoRows},
		{io.EOF,
			[]OptionE{OSupp(sql.ErrNoRows), OStack(), OAnno("sdfsdf"), OStack()},
			nil, nil, "EOF", sql.ErrNoRows},
		{io.EOF,
			[]OptionE{OSupp(sql.ErrNoRows), OAnno("sdfsdf"), OStack(), OValue("a", "b")},
			"a", "b", "EOF", sql.ErrNoRows},
		{io.EOF,
			[]OptionE{OSupp(sql.ErrNoRows), OValue("a", "b"), OAnno("sdfsdf"), OStack()},
			"a", "b", "EOF", sql.ErrNoRows},
		{io.EOF,
			[]OptionE{OValue("a", "b"), OValue("c", "d")},
			"a", "b", "EOF", nil},
		{io.EOF,
			[]OptionE{OValue("a", "b"), OAnno("sdfsdf"), OStack(), errofWrap()},
			"a", "b", "wrapped: EOF", nil},
	}

	for _, tt := range tests {
		err := WrapE(tt.err, tt.opts...)

		if err.Error() != tt.wantMsg {
			t.Errorf("Wrong err.Error(): got: %q, want %q", err.Error(), tt.wantMsg)
		}
		if !IsE(err, tt.err) {
			t.Error("Is failed")
		}
		if tt.key != nil && ValueE(err, tt.key) != tt.value {
			t.Errorf("Wrong Value in err: got: %q, want %q", ValueE(err, tt.key), tt.value)
		}
		if tt.wantSuppressed != nil {
			s := SuppressedE(err)
			if len(s) != 1 {
				t.Errorf("Wrong number of Suppressed in err: got: %q, want %q", len(s), 1)
			} else if s[0] != tt.wantSuppressed {
				t.Errorf("Wrong Suppressed in err: got: %q, want %q", s[0], tt.wantSuppressed)
			}
		}
	}
}
