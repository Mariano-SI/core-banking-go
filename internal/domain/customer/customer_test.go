package customer

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// mustCPF and mustEmail build valid value objects for the tests.
func mustCPF(t *testing.T, raw string) CPF {
	t.Helper()
	cpf, err := NewCPF(raw)
	if err != nil {
		t.Fatalf("NewCPF(%q) error = %v", raw, err)
	}
	return cpf
}

func mustEmail(t *testing.T, raw string) Email {
	t.Helper()
	email, err := NewEmail(raw)
	if err != nil {
		t.Fatalf("NewEmail(%q) error = %v", raw, err)
	}
	return email
}

func validParams(t *testing.T) NewCustomerParams {
	t.Helper()
	birthDate := time.Date(1990, time.May, 10, 0, 0, 0, 0, time.UTC)
	return NewCustomerParams{
		Name:      "  Ana Souza ",
		CPF:       mustCPF(t, "529.982.247-25"),
		Email:     mustEmail(t, "ana@email.com"),
		BirthDate: &birthDate,
	}
}

func TestNewCustomer_Valid(t *testing.T) {
	c, err := NewCustomer(validParams(t))
	if err != nil {
		t.Fatalf("NewCustomer error = %v", err)
	}

	if c.ID() == uuid.Nil || c.ID().Version() != 7 {
		t.Errorf("ID() = %v, want a UUID v7", c.ID())
	}
	if c.Name() != "Ana Souza" {
		t.Errorf("Name() = %q, want %q", c.Name(), "Ana Souza")
	}
	if c.CPF().String() != "52998224725" {
		t.Errorf("CPF() = %q, want %q", c.CPF().String(), "52998224725")
	}
	if c.Email().String() != "ana@email.com" {
		t.Errorf("Email() = %q, want %q", c.Email().String(), "ana@email.com")
	}
	if c.BirthDate() == nil {
		t.Errorf("BirthDate() = nil, want a date")
	}
	if c.CreatedAt().IsZero() || !c.CreatedAt().Equal(c.UpdatedAt()) {
		t.Errorf("CreatedAt() = %v, UpdatedAt() = %v, want equal non-zero times", c.CreatedAt(), c.UpdatedAt())
	}
}

func TestNewCustomer_WithoutBirthDate(t *testing.T) {
	params := validParams(t)
	params.BirthDate = nil

	c, err := NewCustomer(params)
	if err != nil {
		t.Fatalf("NewCustomer error = %v", err)
	}
	if c.BirthDate() != nil {
		t.Errorf("BirthDate() = %v, want nil", c.BirthDate())
	}
}

func TestNewCustomer_GeneratesUniqueIDs(t *testing.T) {
	first, _ := NewCustomer(validParams(t))
	second, _ := NewCustomer(validParams(t))

	if first.ID() == second.ID() {
		t.Errorf("two customers got the same ID: %v", first.ID())
	}
}

func TestNewCustomer_Invalid(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name    string
		modify  func(p *NewCustomerParams)
		wantErr error
	}{
		{name: "empty name", modify: func(p *NewCustomerParams) { p.Name = "" }, wantErr: ErrInvalidName},
		{name: "spaces only name", modify: func(p *NewCustomerParams) { p.Name = "   " }, wantErr: ErrInvalidName},
		{name: "name too long", modify: func(p *NewCustomerParams) { p.Name = strings.Repeat("a", 151) }, wantErr: ErrInvalidName},
		{name: "missing cpf", modify: func(p *NewCustomerParams) { p.CPF = CPF{} }, wantErr: ErrInvalidCPF},
		{name: "missing email", modify: func(p *NewCustomerParams) { p.Email = Email{} }, wantErr: ErrInvalidEmail},
		{name: "birth date in the future", modify: func(p *NewCustomerParams) { p.BirthDate = &future }, wantErr: ErrInvalidBirthDate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := validParams(t)
			tt.modify(&params)

			_, err := NewCustomer(params)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewCustomer error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
