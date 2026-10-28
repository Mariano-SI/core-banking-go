package account

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func validParams(t *testing.T) NewAccountParams {
	t.Helper()
	number, err := NewNumber(1234)
	if err != nil {
		t.Fatalf("NewNumber error = %v", err)
	}
	return NewAccountParams{
		Number:     number,
		BranchID:   uuid.New(),
		CustomerID: uuid.New(),
	}
}

func TestNewAccount_Valid(t *testing.T) {
	params := validParams(t)

	a, err := NewAccount(params)
	if err != nil {
		t.Fatalf("NewAccount error = %v", err)
	}

	if a.ID().Version() != 7 {
		t.Errorf("ID() = %v, want a UUID v7", a.ID())
	}
	if a.Number() != params.Number {
		t.Errorf("Number() = %v, want %v", a.Number(), params.Number)
	}
	if a.BranchID() != params.BranchID {
		t.Errorf("BranchID() = %v, want %v", a.BranchID(), params.BranchID)
	}
	if a.CustomerID() != params.CustomerID {
		t.Errorf("CustomerID() = %v, want %v", a.CustomerID(), params.CustomerID)
	}
	if a.Status() != StatusActive || !a.IsActive() {
		t.Errorf("Status() = %q, IsActive() = %v, want active", a.Status(), a.IsActive())
	}
	if a.CreatedAt().IsZero() || !a.CreatedAt().Equal(a.UpdatedAt()) {
		t.Errorf("CreatedAt() = %v, UpdatedAt() = %v, want equal non-zero times", a.CreatedAt(), a.UpdatedAt())
	}
}

func TestNewAccount_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(p *NewAccountParams)
		wantErr error
	}{
		{name: "missing number", modify: func(p *NewAccountParams) { p.Number = Number{} }, wantErr: ErrInvalidNumber},
		{name: "missing branch", modify: func(p *NewAccountParams) { p.BranchID = uuid.Nil }, wantErr: ErrMissingBranch},
		{name: "missing customer", modify: func(p *NewAccountParams) { p.CustomerID = uuid.Nil }, wantErr: ErrMissingCustomer},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := validParams(t)
			tt.modify(&params)

			_, err := NewAccount(params)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewAccount error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAccount_GeneratesUniqueIDs(t *testing.T) {
	params := validParams(t)

	first, _ := NewAccount(params)
	second, _ := NewAccount(params)

	if first.ID() == second.ID() {
		t.Errorf("two accounts got the same ID: %v", first.ID())
	}
}
