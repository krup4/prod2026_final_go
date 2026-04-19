package exceptions

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrExperimentExists   = errors.New("experiment already exists")
	ErrActiveExperiment   = errors.New("active experiment with this flag already exists")
	ErrExperimentNotFound = errors.New("experiment not found")
)
