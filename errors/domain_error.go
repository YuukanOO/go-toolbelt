package errors

import "fmt"

// DomainError represents an error thrown by the domain.
type DomainError struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Errors  []error `json:"errors"`
}

func (err DomainError) Error() string {
	msg := fmt.Sprintf("%s - %s", err.Code, err.Message)

	for _, e := range err.Errors {
		msg = msg + "\n\t" + e.Error()
	}

	return msg
}

// NewDomainError instantiates a new domain error.
func NewDomainError(code string, message string, errors ...error) error {
	return &DomainError{
		Code:    code,
		Message: message,
		Errors:  errors,
	}
}
