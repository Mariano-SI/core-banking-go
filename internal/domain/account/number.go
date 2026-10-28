package account

import (
	"errors"
	"fmt"
	"strings"
)

const minBaseDigits = 6

var ErrInvalidNumber = errors.New("invalid account number")

type Number struct {
	value string
}

func NewNumber(sequence int64) (Number, error) {
	if sequence <= 0 {
		return Number{}, ErrInvalidNumber
	}

	base := fmt.Sprintf("%0*d", minBaseDigits, sequence)

	return Number{value: fmt.Sprintf("%s-%c", base, luhnCheckDigit(base))}, nil
}

func ParseNumber(raw string) (Number, error) {
	base, check, found := strings.Cut(strings.TrimSpace(raw), "-")
	if !found || len(base) < minBaseDigits || len(check) != 1 {
		return Number{}, ErrInvalidNumber
	}

	if !isDigits(base) || !isDigits(check) {
		return Number{}, ErrInvalidNumber
	}

	if luhnCheckDigit(base) != check[0] {
		return Number{}, ErrInvalidNumber
	}

	return Number{value: base + "-" + check}, nil
}

func (n Number) String() string {
	return n.value
}

func luhnCheckDigit(base string) byte {
	sum := 0
	double := true

	for i := len(base) - 1; i >= 0; i-- {
		digit := int(base[i] - '0')
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}

	return byte('0' + (10-sum%10)%10)
}

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
