package customer

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const maxNameLength = 150

var (
	ErrInvalidName      = errors.New("invalid name")
	ErrInvalidBirthDate = errors.New("invalid birth date")
)

type Customer struct {
	id        uuid.UUID
	name      string
	cpf       CPF
	email     Email
	birthDate *time.Time
	createdAt time.Time
	updatedAt time.Time
}

type NewCustomerParams struct {
	Name      string
	CPF       CPF
	Email     Email
	BirthDate *time.Time
}

func NewCustomer(params NewCustomerParams) (Customer, error) {
	name := strings.TrimSpace(params.Name)

	if name == "" || utf8.RuneCountInString(name) > maxNameLength {
		return Customer{}, ErrInvalidName
	}

	if params.CPF == (CPF{}) {
		return Customer{}, ErrInvalidCPF
	}
	if params.Email == (Email{}) {
		return Customer{}, ErrInvalidEmail
	}

	now := time.Now().UTC()

	var birthDate *time.Time
	if params.BirthDate != nil {
		if params.BirthDate.After(now) {
			return Customer{}, ErrInvalidBirthDate
		}
		date := *params.BirthDate
		birthDate = &date
	}

	id, err := uuid.NewV7()
	if err != nil {
		return Customer{}, fmt.Errorf("generate customer id: %w", err)
	}

	return Customer{
		id:        id,
		name:      name,
		cpf:       params.CPF,
		email:     params.Email,
		birthDate: birthDate,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (c Customer) ID() uuid.UUID { return c.id }

func (c Customer) Name() string { return c.name }

func (c Customer) CPF() CPF { return c.cpf }

func (c Customer) Email() Email { return c.email }

func (c Customer) BirthDate() *time.Time {
	if c.birthDate == nil {
		return nil
	}
	date := *c.birthDate
	return &date
}

func (c Customer) CreatedAt() time.Time { return c.createdAt }

func (c Customer) UpdatedAt() time.Time { return c.updatedAt }
