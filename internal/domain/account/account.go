package account

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrMissingBranch   = errors.New("account branch is required")
	ErrMissingCustomer = errors.New("account customer is required")
)

type Account struct {
	id         uuid.UUID
	number     Number
	branchID   uuid.UUID
	customerID uuid.UUID
	status     Status
	createdAt  time.Time
	updatedAt  time.Time
}

type NewAccountParams struct {
	Number     Number
	BranchID   uuid.UUID
	CustomerID uuid.UUID
}

func NewAccount(p NewAccountParams) (Account, error) {
	if p.Number == (Number{}) {
		return Account{}, ErrInvalidNumber
	}
	if p.BranchID == uuid.Nil {
		return Account{}, ErrMissingBranch
	}
	if p.CustomerID == uuid.Nil {
		return Account{}, ErrMissingCustomer
	}

	id, err := uuid.NewV7()
	if err != nil {
		return Account{}, fmt.Errorf("generate account id: %w", err)
	}

	now := time.Now().UTC()

	return Account{
		id:         id,
		number:     p.Number,
		branchID:   p.BranchID,
		customerID: p.CustomerID,
		status:     StatusActive,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

func (a Account) ID() uuid.UUID { return a.id }

func (a Account) Number() Number { return a.number }

func (a Account) BranchID() uuid.UUID { return a.branchID }

func (a Account) CustomerID() uuid.UUID { return a.customerID }

func (a Account) Status() Status { return a.status }

func (a Account) IsActive() bool { return a.status == StatusActive }

func (a Account) CreatedAt() time.Time { return a.createdAt }

func (a Account) UpdatedAt() time.Time { return a.updatedAt }
