package account

import (
	"errors"
	"testing"
)

func TestParseStatus_Valid(t *testing.T) {
	tests := []struct {
		input string
		want  Status
	}{
		{input: "active", want: StatusActive},
		{input: "blocked", want: StatusBlocked},
		{input: "closed", want: StatusClosed},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseStatus(tt.input)
			if err != nil {
				t.Fatalf("ParseStatus(%q) error = %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("ParseStatus(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseStatus_Invalid(t *testing.T) {
	for _, input := range []string{"", "ACTIVE", "frozen", " active"} {
		t.Run(input, func(t *testing.T) {
			if _, err := ParseStatus(input); !errors.Is(err, ErrInvalidStatus) {
				t.Errorf("ParseStatus(%q) error = %v, want %v", input, err, ErrInvalidStatus)
			}
		})
	}
}
