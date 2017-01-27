package validation

import (
	"fmt"
)

// FieldError represents a field error that has occured during the validation step.
type FieldError struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Code     string `json:"code"`
}

func (err *FieldError) Error() string {
	return fmt.Sprintf("Validation failed for resource \"%s\", field \"%s\" with the reason \"%s\"",
		err.Resource, err.Field, err.Code)
}

func newFieldError(resource string, field string, code string) *FieldError {
	return &FieldError{
		Resource: resource,
		Field:    field,
		Code:     code,
	}
}
