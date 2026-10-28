package account

import "errors"

type Status string

var ErrInvalidStatus = errors.New("invalid account status")

const (
	StatusActive  Status = "active"
	StatusBlocked Status = "blocked"
	StatusClosed  Status = "closed"
)

func ParseStatus(raw string) (Status, error) {
	parsed := Status(raw)

	switch parsed {
	case StatusActive, StatusBlocked, StatusClosed:
		return parsed, nil
	default:
		return "", ErrInvalidStatus
	}
}
