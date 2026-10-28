package customer

import (
	"errors"
	"strings"
	"testing"
)

func TestNewEmail_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "simple", input: "ana@email.com", want: "ana@email.com"},
		{name: "uppercase and spaces", input: "  Ana.Souza@Email.COM ", want: "ana.souza@email.com"},
		{name: "plus and subdomain", input: "ana+bank@email.com.br", want: "ana+bank@email.com.br"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.input)
			if err != nil {
				t.Fatalf("NewEmail(%q) error = %v", tt.input, err)
			}
			if email.String() != tt.want {
				t.Errorf("String() = %q, want %q", email.String(), tt.want)
			}
		})
	}
}

func TestNewEmail_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "spaces only", input: "   "},
		{name: "missing at", input: "ana"},
		{name: "missing domain", input: "ana@"},
		{name: "missing local part", input: "@email.com"},
		{name: "domain without dot", input: "ana@email"},
		{name: "display name", input: "Ana <ana@email.com>"},
		{name: "space inside", input: "ana souza@email.com"},
		{name: "too long", input: strings.Repeat("a", 250) + "@email.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEmail(tt.input)
			if !errors.Is(err, ErrInvalidEmail) {
				t.Errorf("NewEmail(%q) error = %v, want %v", tt.input, err, ErrInvalidEmail)
			}
		})
	}
}
