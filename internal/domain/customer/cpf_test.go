package customer

import (
	"errors"
	"testing"
)

func TestNewCPF_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "digits only", input: "52998224725", want: "52998224725"},
		{name: "formatted", input: "529.982.247-25", want: "52998224725"},
		{name: "surrounding spaces", input: "  111.444.777-35 ", want: "11144477735"},
		{name: "first check digit zero", input: "123.456.789-09", want: "12345678909"},
		{name: "second check digit zero", input: "935.411.347-80", want: "93541134780"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpf, err := NewCPF(tt.input)
			if err != nil {
				t.Fatalf("NewCPF(%q) error = %v", tt.input, err)
			}
			if cpf.String() != tt.want {
				t.Errorf("String() = %q, want %q", cpf.String(), tt.want)
			}
		})
	}
}

func TestNewCPF_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "too short", input: "5299822472"},
		{name: "too long", input: "529982247251"},
		{name: "letters", input: "529.982.247-2A"},
		{name: "all digits equal", input: "111.111.111-11"},
		{name: "wrong first check digit", input: "529.982.247-35"},
		{name: "wrong second check digit", input: "529.982.247-26"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCPF(tt.input)
			if !errors.Is(err, ErrInvalidCPF) {
				t.Errorf("NewCPF(%q) error = %v, want %v", tt.input, err, ErrInvalidCPF)
			}
		})
	}
}

func TestCPF_Formatted(t *testing.T) {
	cpf, err := NewCPF("52998224725")
	if err != nil {
		t.Fatalf("NewCPF error = %v", err)
	}

	if got := cpf.Formatted(); got != "529.982.247-25" {
		t.Errorf("Formatted() = %q, want %q", got, "529.982.247-25")
	}
}
