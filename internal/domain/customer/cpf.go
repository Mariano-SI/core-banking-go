package customer

import (
	"errors"
	"strings"
)

var ErrInvalidCPF = errors.New("invalid cpf")

type CPF struct {
	digits string
}

func (c CPF) String() string {
	return c.digits
}

func (c CPF) Formatted() string {
	d := c.digits
	return d[0:3] + "." + d[3:6] + "." + d[6:9] + "-" + d[9:11]
}

func NewCPF(rawCpf string) (CPF, error) {
	digits, ok := normalize(rawCpf)
	if !ok || len(digits) != 11 {
		return CPF{}, ErrInvalidCPF
	}

	if strings.Count(digits, digits[:1]) == len(digits) {
		return CPF{}, ErrInvalidCPF
	}

	if checkDigit(digits[:9]) != digits[9] || checkDigit(digits[:10]) != digits[10] {
		return CPF{}, ErrInvalidCPF
	}

	return CPF{digits: digits}, nil
}

func normalize(raw string) (string, bool) {
	var b strings.Builder

	for _, r := range strings.TrimSpace(raw) {
		switch {
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '.' || r == '-':
		default:
			return "", false
		}
	}

	return b.String(), true
}

func checkDigit(digits string) byte {
	digitsLength := len(digits) + 1
	acc := 0

	for _, digit := range digits {
		acc += int(digit-'0') * digitsLength
		digitsLength--
	}

	remainder := acc % 11

	if remainder < 2 {
		return '0'
	}

	return byte('0' + 11 - remainder)
}
