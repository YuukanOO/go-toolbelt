# go-toolbelt

This is my personal toolbelt for the golang awesome programming language. It is built around ddd practices.

## Installation

`go get github.com/YuukanOO/go-toolbelt`

## Usage

### Errors

A `DomainError` struct is defined and encapsulates a Domain error. A domain error is an expected error in case some input were wrong and will eventually be displayed to the user.

```go
err := errors.NewDomainError("AConstantCode", "A friendly message for the developper", errors.New("Any number of errors"))
```

### Validation

This is a very simple fluent like API that uses [the go-playground validator](https://github.com/go-playground/validator) under the hood.

```go
err := validation.Validate("User").
  Field("username", "mytoolongusername", "required,max=10,min=1").
  Field("password", "aS3cretP@ssw0rd", "required,min=10").
  Errors() // Will trigger the evaluation

// If it has an error, a domain error will be returned
domErr := err.(*errors.DomainError)

// And inner Errors will be of type FieldError
fieldErr := domErr.Errors[0].(*validation.FieldError)
```

Don't hesitate to check the tests for more examples.