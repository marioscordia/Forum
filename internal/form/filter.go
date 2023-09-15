package form

import "newforum/internal/validator"

type Filter struct {
	Category []string
	validator.Validator
}
