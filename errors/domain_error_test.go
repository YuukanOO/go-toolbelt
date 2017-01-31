package errors

import (
	"errors"
	"testing"
)

const (
	TestErrCode = "TestingError"
	TestErrDesc = "An error has occured"
)

var TestInnerErr = errors.New("An inner error")

func TestDomainError(t *testing.T) {
	err := NewDomainError(TestErrCode, TestErrDesc, TestInnerErr)

	domErr := err.(*DomainError)

	if domErr == nil {
		t.Error("Error should be of type DomainError")
	}

	if domErr.Code != TestErrCode || domErr.Message != TestErrDesc {
		t.Error("Code and message should match those provided")
	}
}
