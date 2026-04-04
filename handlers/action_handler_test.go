package handlers_test

import (
	"errors"
	"testing"

	"github.com/janmbaco/go-redux/v2/actions"
	"github.com/janmbaco/go-redux/v2/handlers"
)

type counterState struct {
	Count int
}

func TestActionHandlerHandle_ShouldUpdateState_WhenActionMatches(t *testing.T) {
	// Arrange
	incrementAction := actions.NewAction[int]("counter/increment")
	handler := handlers.NewActionHandlerBuilder[counterState]().
		On(incrementAction, func(state counterState, amount int) counterState {
			state.Count += amount
			return state
		}).
		Build()

	// Act
	newState, err := handler.Handle(counterState{Count: 2}, incrementAction.With(3))

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if newState.Count != 5 {
		t.Fatalf("expected count %d, got %d", 5, newState.Count)
	}
}

func TestActionHandlerHandle_ShouldPassNilInterfacePayload_WhenPayloadIsNil(t *testing.T) {
	// Arrange
	resetAction := actions.NewAction[any]("counter/reset")
	handler := handlers.NewActionHandlerBuilder[counterState]().
		On(resetAction, func(state counterState, payload any) counterState {
			if payload == nil {
				state.Count = 0
			}
			return state
		}).
		Build()

	// Act
	newState, err := handler.Handle(counterState{Count: 9}, resetAction.With(nil))

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if newState.Count != 0 {
		t.Fatalf("expected count %d, got %d", 0, newState.Count)
	}
}

func TestActionHandlerHandle_ShouldReturnReducerError_WhenReducerFails(t *testing.T) {
	// Arrange
	expectedErr := errors.New("boom")
	incrementAction := actions.NewAction[int]("counter/error")
	handler := handlers.NewActionHandlerBuilder[counterState]().
		On(incrementAction, func(state counterState, amount int) (counterState, error) {
			return state, expectedErr
		}).
		Build()

	// Act
	_, err := handler.Handle(counterState{Count: 1}, incrementAction.With(1))

	// Assert
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
