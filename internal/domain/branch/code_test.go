package branch

import (
	"errors"
	"testing"
)

func TestNewCode_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "four digits", input: "0001", want: "0001"},
		{name: "surrounding spaces", input: " 0002 ", want: "0002"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := NewCode(tt.input)
			if err != nil {
				t.Fatalf("NewCode(%q) error = %v", tt.input, err)
			}
			if code.String() != tt.want {
				t.Errorf("String() = %q, want %q", code.String(), tt.want)
			}
		})
	}
}

func TestNewCode_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "spaces only", input: "    "},
		{name: "too short", input: "001"},
		{name: "too long", input: "00001"},
		{name: "letter", input: "00A1"},
		{name: "symbol", input: "00-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCode(tt.input)
			if !errors.Is(err, ErrInvalidCode) {
				t.Errorf("NewCode(%q) error = %v, want %v", tt.input, err, ErrInvalidCode)
			}
		})
	}
}
