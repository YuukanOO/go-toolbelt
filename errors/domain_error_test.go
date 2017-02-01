package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

const (
	TestErrCode = "TestingError"
	TestErrDesc = "An error has occured"
)

var ErrTest = errors.New("An inner error")

func TestDomainError(t *testing.T) {
	err := NewDomainError(TestErrCode, TestErrDesc, ErrTest)

	domErr := err.(*DomainError)

	if domErr == nil {
		t.Error("Error should be of type DomainError")
	}

	if domErr.Code != TestErrCode || domErr.Message != TestErrDesc {
		t.Error("Code and message should match those provided")
	}

	if len(domErr.Errors) == 0 {
		t.Error("It should have one inner error")
	}

	errStr := domErr.Error()

	if !strings.HasPrefix(errStr, fmt.Sprintf("%s - %s", TestErrCode, TestErrDesc)) {
		t.Error("Wrong error message")
	}
}
