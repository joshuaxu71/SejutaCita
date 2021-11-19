package models

import "github.com/go-playground/validator"

// A boolean value that is returned in the response to denote success
// swagger:response booleanResponse
type booleanResponseWrapper struct {
	// in:body
	Body struct {
		Success bool
	}
}

// swagger:enum SortOrder
type SortOrder int

const (
	Asc  SortOrder = 2
	Desc SortOrder = 1
)

func EmptyValidate(fl validator.FieldLevel) bool {
	return true
}
