package models

import "github.com/go-playground/validator"

type SortOrder int64

const (
	Asc  SortOrder = 1
	Desc SortOrder = -1
)

func EmptyValidate(fl validator.FieldLevel) bool {
	return true
}
