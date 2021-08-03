package redux

import (
	"github.com/janmbaco/go-infrastructure/errors"
	"reflect"
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

type StoreError struct {
	errors.CustomError
	ErrorType StoreErrorType
}

func newStoreError(errorType StoreErrorType, message string) *StoreError {
	return &StoreError{ErrorType: errorType, CustomError: errors.CustomError{Message: message}}
}

type storeErrorPipe struct{}

func (storeErrorPipe *storeErrorPipe) Pipe(err error) error {
	resultError := err

	if errType := reflect.Indirect(reflect.ValueOf(err)).Type(); errType != reflect.TypeOf(&StoreError{}) {
		errorType := UnexpectedStoreError
		resultError = &StoreError{
			CustomError: errors.CustomError{
				Message:       err.Error(),
				InternalError: err,
			},
			ErrorType: errorType,
		}
	}
	return resultError
}
