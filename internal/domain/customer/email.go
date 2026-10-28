package customer

import (
	"errors"
	"net/mail"
	"strings"
)

var ErrInvalidEmail = errors.New("invalid email")

const maxEmailLength = 254

type Email struct {
	address string
}

func NewEmail(rawEmail string) (Email, error) {
	address := strings.ToLower(strings.TrimSpace(rawEmail))

	if address == "" || len(address) > maxEmailLength {
		return Email{}, ErrInvalidEmail
	}

	parsed, err := mail.ParseAddress(address)
	if err != nil || parsed.Address != address {
		return Email{}, ErrInvalidEmail
	}

	domain := address[strings.LastIndex(address, "@")+1:]
	if !strings.Contains(domain, ".") {
		return Email{}, ErrInvalidEmail
	}

	return Email{
		address: address,
	}, nil
}

func (e Email) String() string {
	return e.address
}
