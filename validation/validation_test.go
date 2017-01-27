package validation

import (
	"testing"

	"github.com/YuukanOO/go-toolbelt/errors"
)

func TestFluentValidationFailure(t *testing.T) {
	err := Validate("User").
		Field("username", "bob", "required").
		Field("password", "", "required").Errors()

	if err == nil {
		t.Error("Error should not be nil")
	}

	domErr := err.(*errors.DomainError)

	if domErr == nil {
		t.Error("Error should be of type DomainError")
	}

	if domErr.Code != FailedErrCode {
		t.Error("Code should be equal to", FailedErrCode)
	}

	if len(domErr.Errors) != 1 {
		t.Error("It should only have one error")
	}

	fieldErr := domErr.Errors[0].(*FieldError)

	if fieldErr == nil {
		t.Error("Inner error should be of type FieldError")
	}

	if fieldErr.Resource != "User" ||
		fieldErr.Field != "password" ||
		fieldErr.Code != "required" {
		t.Error("It should complains about the password field required")
	}
}

func TestFluentValidationManyTagsFailure(t *testing.T) {
	err := Validate("User").
		Field("username", "bobismyusername", "required,max=10,min=2").Errors()

	if err == nil {
		t.Error("It should find an error")
	}

	domErr := err.(*errors.DomainError)

	if len(domErr.Errors) != 1 {
		t.Error("It should find 1 error")
	}

	fieldErr := domErr.Errors[0].(*FieldError)

	if fieldErr.Field != "username" || fieldErr.Code != "max" {
		t.Error("It should fail on the max tag")
	}
}

func TestFluentValidationPass(t *testing.T) {
	if err := Validate("User").
		Field("username", "bob", "required").
		Field("password", "bobpwd", "required").Errors(); err != nil {
		t.Error("Everything should be fine here")
	}
}
