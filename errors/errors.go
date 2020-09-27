package errors

import (
	"fmt"
	"strings"
)

// Custom errors code
const (
	InternalErrorCode = iota
	NotFoundCode
	BadRequestCode
)

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewCustomError(code int, message string) *CustomError {
	return &CustomError{
		code,
		message,
	}
}

func (c CustomError) WithParams(params []string) *CustomError {
	return &CustomError{
		c.Code,
		fmt.Sprintf(c.Message, strings.Join(params, ",")),
	}
}

func (errMsg CustomError) Error() string {
	return errMsg.Message
}

var (
	InternalError  = NewCustomError(InternalErrorCode, "Internal server error.")
	NotFoundItems  = NewCustomError(NotFoundCode, "Items not found: %s.")
	InvalidItems   = NewCustomError(BadRequestCode, "Invalid items: %s.")
	AtLeastOneItem = NewCustomError(BadRequestCode, "You must provide at least one item code.")
	InvalidFormat  = NewCustomError(BadRequestCode, "Request invalid format.")
	MaxItemsExceded  = NewCustomError(BadRequestCode, "Max items quantity exceded.")
)
