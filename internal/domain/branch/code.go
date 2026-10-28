package branch

import (
	"errors"
	"strings"
)

var ErrInvalidCode = errors.New("invalid branch code")

type Code struct {
	value string
}

func (c Code) String() string {
	return c.value
}

func NewCode(raw string) (Code, error) {
	value, err := normalizeCode(raw)
	if err != nil {
		return Code{}, err
	}

	return Code{value: value}, nil
}

func normalizeCode(raw string) (string, error) {
	code := strings.TrimSpace(raw)
	if len(code) != 4 {
		return "", ErrInvalidCode
	}

	for _, r := range code {
		if r < '0' || r > '9' {
			return "", ErrInvalidCode
		}
	}

	return code, nil
}
