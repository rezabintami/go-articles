package constants

import "errors"

const (
	DUPLICATE_DATA_UNIQUE = "duplicate key value violates unique constraint"
	PARAMETER_MUST_BE_FILLED = "Parameter must be filled"
)

var (
	ErrInternalServerError = errors.New("Internal Server error")
	ErrRecordNotFound      = errors.New("Record not found")
	ErrDoNotHavePermission = errors.New("You Do not have permission to do the action")
	ErrDuplicateData = errors.New("Duplicate data")
	ErrEmailPasswordNotFound = errors.New("(Email) or (Password) empty")
)