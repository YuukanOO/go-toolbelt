package validation

import (
	"github.com/YuukanOO/go-toolbelt/errors"
	"gopkg.in/go-playground/validator.v9"
)

// Validator instance exposed just in case you need more control over validation.
var Validator = validator.New()

// FailedErrCode error code constant for validation errors.
const FailedErrCode = "ValidationFailed"

// Chain represents a validation facility to ease the validation process
// by providing a Fluent like API.
type Chain struct {
	resource string
	fields   []*chainField
}

type chainField struct {
	name  string
	value interface{}
	other interface{}
	tag   string
}

// Validate instantiate a new validation chain for the given resource. Then, you can
// use a fluent like API to constructs your validations and call Errors() to actually
// evaluates the chain.
func Validate(resource string) *Chain {
	return &Chain{
		resource: resource,
		fields:   []*chainField{},
	}
}

// FieldWithValue adds a validation against another value to the chain.
func (chain *Chain) FieldWithValue(name string, value interface{}, other interface{}, tag string) *Chain {
	chain.fields = append(chain.fields, &chainField{
		name:  name,
		other: other,
		value: value,
		tag:   tag,
	})
	return chain
}

// Field adds a validation for the given field name.
func (chain *Chain) Field(name string, value interface{}, tag string) *Chain {
	return chain.FieldWithValue(name, value, nil, tag)
}

// Errors evaluates the entire chain and returns any errors as a DomainError.
func (chain *Chain) Errors() error {
	var fieldErrs []error

	for _, f := range chain.fields {
		var err error

		if f.other == nil {
			err = Validator.Var(f.value, f.tag)
		} else {
			err = Validator.VarWithValue(f.value, f.other, f.tag)
		}

		if err != nil {
			if valErr := err.(validator.ValidationErrors); valErr != nil {
				for _, fieldErr := range valErr {
					fieldErrs = append(fieldErrs, newFieldError(chain.resource, f.name, fieldErr.ActualTag()))
				}
			}
		}
	}

	if len(fieldErrs) > 0 {
		return errors.NewDomainError(FailedErrCode, "Validation failed", fieldErrs...)
	}

	return nil
}
