package branch

import (
	"errors"
	"strings"
	"testing"
)

// mustCode builds a valid branch code for the tests.
func mustCode(t *testing.T, raw string) Code {
	t.Helper()
	code, err := NewCode(raw)
	if err != nil {
		t.Fatalf("NewCode(%q) error = %v", raw, err)
	}
	return code
}

func validParams(t *testing.T) NewBranchParams {
	t.Helper()
	return NewBranchParams{
		Name: "  Agência Central ",
		Code: mustCode(t, "0001"),
	}
}

func TestNewBranch_Valid(t *testing.T) {
	b, err := NewBranch(validParams(t))
	if err != nil {
		t.Fatalf("NewBranch error = %v", err)
	}

	if b.ID().Version() != 7 {
		t.Errorf("ID() = %v, want a UUID v7", b.ID())
	}
	if b.Name() != "Agência Central" {
		t.Errorf("Name() = %q, want %q", b.Name(), "Agência Central")
	}
	if b.Code().String() != "0001" {
		t.Errorf("Code() = %q, want %q", b.Code().String(), "0001")
	}
	if !b.IsActive() {
		t.Errorf("IsActive() = false, want true")
	}
	if b.CreatedAt().IsZero() {
		t.Errorf("CreatedAt() is zero, want the creation time")
	}
}

func TestNewBranch_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(p *NewBranchParams)
		wantErr error
	}{
		{name: "empty name", modify: func(p *NewBranchParams) { p.Name = "" }, wantErr: ErrInvalidBranchName},
		{name: "spaces only name", modify: func(p *NewBranchParams) { p.Name = "   " }, wantErr: ErrInvalidBranchName},
		{name: "name too long", modify: func(p *NewBranchParams) { p.Name = strings.Repeat("a", 101) }, wantErr: ErrInvalidBranchName},
		{name: "missing code", modify: func(p *NewBranchParams) { p.Code = Code{} }, wantErr: ErrInvalidCode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := validParams(t)
			tt.modify(&params)

			_, err := NewBranch(params)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewBranch error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestBranch_EnsureCanOpenAccount(t *testing.T) {
	b, err := NewBranch(validParams(t))
	if err != nil {
		t.Fatalf("NewBranch error = %v", err)
	}

	if err := b.EnsureCanOpenAccount(); err != nil {
		t.Errorf("active branch: EnsureCanOpenAccount() = %v, want nil", err)
	}

	b.Deactivate()

	if b.IsActive() {
		t.Errorf("IsActive() after Deactivate() = true, want false")
	}
	if err := b.EnsureCanOpenAccount(); !errors.Is(err, ErrBranchInactive) {
		t.Errorf("inactive branch: EnsureCanOpenAccount() = %v, want %v", err, ErrBranchInactive)
	}
}
