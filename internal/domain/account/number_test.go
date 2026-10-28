package account

import (
	"errors"
	"testing"
)

func TestNewNumber_Valid(t *testing.T) {
	tests := []struct {
		sequence int64
		want     string
	}{
		{sequence: 1, want: "000001-8"},
		{sequence: 123, want: "000123-0"},
		{sequence: 1234, want: "001234-4"},
		{sequence: 123456, want: "123456-6"},
		{sequence: 1000000, want: "1000000-8"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			number, err := NewNumber(tt.sequence)
			if err != nil {
				t.Fatalf("NewNumber(%d) error = %v", tt.sequence, err)
			}
			if number.String() != tt.want {
				t.Errorf("NewNumber(%d) = %q, want %q", tt.sequence, number.String(), tt.want)
			}
		})
	}
}

func TestNewNumber_Invalid(t *testing.T) {
	for _, sequence := range []int64{0, -1} {
		if _, err := NewNumber(sequence); !errors.Is(err, ErrInvalidNumber) {
			t.Errorf("NewNumber(%d) error = %v, want %v", sequence, err, ErrInvalidNumber)
		}
	}
}

func TestParseNumber_Valid(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "001234-4", want: "001234-4"},
		{input: " 123456-6 ", want: "123456-6"},
		{input: "7992739871-3", want: "7992739871-3"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			number, err := ParseNumber(tt.input)
			if err != nil {
				t.Fatalf("ParseNumber(%q) error = %v", tt.input, err)
			}
			if number.String() != tt.want {
				t.Errorf("ParseNumber(%q) = %q, want %q", tt.input, number.String(), tt.want)
			}
		})
	}
}

func TestParseNumber_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "missing hyphen", input: "0012344"},
		{name: "wrong check digit", input: "001234-5"},
		{name: "base too short", input: "01234-4"},
		{name: "letter in base", input: "00A234-4"},
		{name: "missing check digit", input: "001234-"},
		{name: "two check digits", input: "001234-44"},
		{name: "extra hyphen", input: "001234-4-4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ParseNumber(tt.input); !errors.Is(err, ErrInvalidNumber) {
				t.Errorf("ParseNumber(%q) error = %v, want %v", tt.input, err, ErrInvalidNumber)
			}
		})
	}
}

func TestNumber_RoundTrip(t *testing.T) {
	created, err := NewNumber(42)
	if err != nil {
		t.Fatalf("NewNumber error = %v", err)
	}

	parsed, err := ParseNumber(created.String())
	if err != nil {
		t.Fatalf("ParseNumber(%q) error = %v", created.String(), err)
	}
	if parsed != created {
		t.Errorf("ParseNumber(%q) = %v, want %v", created.String(), parsed, created)
	}
}
