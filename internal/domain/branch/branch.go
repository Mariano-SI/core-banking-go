package branch

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const branchNameMaxSize = 100

var (
	ErrInvalidBranchName = errors.New("invalid branch name")
	ErrBranchInactive    = errors.New("branch is inactive")
)

type Branch struct {
	id        uuid.UUID
	name      string
	code      Code
	active    bool
	createdAt time.Time
}

type NewBranchParams struct {
	Name string
	Code Code
}

func NewBranch(params NewBranchParams) (Branch, error) {
	name := strings.TrimSpace(params.Name)

	if name == "" || utf8.RuneCountInString(name) > branchNameMaxSize {
		return Branch{}, ErrInvalidBranchName
	}

	if params.Code == (Code{}) {
		return Branch{}, ErrInvalidCode
	}

	now := time.Now().UTC()
	id, err := uuid.NewV7()
	if err != nil {
		return Branch{}, fmt.Errorf("generate branch id: %w", err)
	}

	return Branch{
		id:        id,
		name:      name,
		code:      params.Code,
		active:    true,
		createdAt: now,
	}, nil
}

func (b *Branch) ID() uuid.UUID {
	return b.id
}

func (b *Branch) Name() string {
	return b.name
}

func (b *Branch) Code() Code {
	return b.code
}

func (b *Branch) CreatedAt() time.Time {
	return b.createdAt
}

func (b *Branch) Deactivate() {
	b.active = false
}

func (b *Branch) IsActive() bool {
	return b.active
}

func (b *Branch) EnsureCanOpenAccount() error {
	if !b.active {
		return ErrBranchInactive
	}
	return nil
}
