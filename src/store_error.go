package redux

import (
	"github.com/janmbaco/go-infrastructure/errors"
)

type StoreErrorType uint8

const (
	UnexpectedStoreError StoreErrorType = iota
	AnyReducerForThisActionError
	EmptySelectorError
	AnyStateBySelectorError
	MultipleReducerForSelectorError
	MultipleReducerForActionsObjectError
)

type StoreError interface {
	errors.CustomError
	GetErrorType() StoreErrorType
}

type storeError struct {
	errors.CustomizableError
	ErrorType StoreErrorType
}

func newStoreError(errorType StoreErrorType, message string) StoreError {
	return &storeError{
		CustomizableError: errors.CustomizableError{
			Message: message,
			InternalError: nil,
		},
		ErrorType: errorType, 
	}
}

func (e *storeError) GetErrorType() StoreErrorType {
	return e.ErrorType
}
