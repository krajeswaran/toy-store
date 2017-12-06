package errors

import (
	"fmt"
)

// ApplicationError captures errors related to blaance service
type ApplicationError struct {
	ErrorType int
	Message   string
}

// NewError to create new error type
func NewError(eype int, msg string) ApplicationError {
	err, ok := ErrorMap[eype]
	if !ok {
		err = ErrorMap[UnknownError] // unknown err
	}

	if msg != "" {
		err.Message += msg
	}

	return err
}

var (
	// ErrorMap mapping of error codes and messages
	ErrorMap = map[int]ApplicationError{
		UnknownError:      {999, "Unknown Error - "},
		ReplenishError:    {100, "Can't replenish inventory - "},
		OrderError:        {101, "Error placing order - "},
		PurchaseDoneBreak: {200, "All done, place order - "},
	}
)

const (
	UnknownError = iota
	ReplenishError
	OrderError
	PurchaseDoneBreak
)

// Error to format errors
func (e ApplicationError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}
