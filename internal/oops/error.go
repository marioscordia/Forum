package oops

import "errors"

var (
	ErrFormInvalid = errors.New("form: invalid")

	ErrNoRecord = errors.New("models: no matching record found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")
)