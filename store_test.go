package redux_test

import (
	"testing"

	"github.com/janmbaco/go-infrastructure/v2/logs"
	redux "github.com/janmbaco/go-redux/v2"
	"github.com/janmbaco/go-redux/v2/handlers"
)

type counterState struct {
	Count int
}

func TestStoreDispatch_ShouldUpdateState_WhenHandlerMatches(t *testing.T) {
	// Arrange
	logger := logs.NewLogger()
	logger.Mute()
	store := redux.NewStore(counterState{Count: 1}, logger)
	defer store.Close()

	incrementAction := redux.NewAction[int]("counter/increment")
	handler := handlers.NewActionHandlerBuilder[counterState]().
		On(incrementAction, func(state counterState, amount int) counterState {
			state.Count += amount
			return state
		}).
		Build()

	if err := store.AddModule(handler); err != nil {
		t.Fatalf("expected no add module error, got %v", err)
	}

	// Act
	err := store.Dispatch(incrementAction.With(2))

	// Assert
	if err != nil {
		t.Fatalf("expected no dispatch error, got %v", err)
	}

	if state := store.GetState(); state.Count != 3 {
		t.Fatalf("expected count %d, got %d", 3, state.Count)
	}
}

func TestStoreDispatch_ShouldUpdateSelectedSlice_WhenStateIsMap(t *testing.T) {
	// Arrange
	logger := logs.NewLogger()
	logger.Mute()
	store := redux.NewStore(map[string]any{}, logger)
	defer store.Close()

	incrementAction := redux.NewAction[int]("counter/map/increment")
	handler := handlers.NewActionHandlerBuilder[counterState]().
		SetInitialState(counterState{Count: 0}).
		SetSelector("counter").
		On(incrementAction, func(state counterState, amount int) counterState {
			state.Count += amount
			return state
		}).
		Build()

	if err := store.AddModule(handler); err != nil {
		t.Fatalf("expected no add module error, got %v", err)
	}

	// Act
	err := store.Dispatch(incrementAction.With(5))

	// Assert
	if err != nil {
		t.Fatalf("expected no dispatch error, got %v", err)
	}

	state := store.GetState()
	counter, ok := state["counter"].(counterState)
	if !ok {
		t.Fatalf("expected counter slice to be present")
	}

	if counter.Count != 5 {
		t.Fatalf("expected count %d, got %d", 5, counter.Count)
	}
}

func TestStoreDispatch_ShouldPublishState_WhenActionChangesState(t *testing.T) {
	// Arrange
	logger := logs.NewLogger()
	logger.Mute()
	store := redux.NewStore(counterState{Count: 0}, logger)
	defer store.Close()

	incrementAction := redux.NewAction[int]("counter/publish")
	handler := handlers.NewActionHandlerBuilder[counterState]().
		On(incrementAction, func(state counterState, amount int) counterState {
			state.Count += amount
			return state
		}).
		Build()
	if err := store.AddModule(handler); err != nil {
		t.Fatalf("expected no add module error, got %v", err)
	}

	var received counterState
	callback := func(state counterState) {
		received = state
	}
	if err := store.Subscribe(&callback); err != nil {
		t.Fatalf("expected no subscribe error, got %v", err)
	}

	// Act
	err := store.Dispatch(incrementAction.With(4))

	// Assert
	if err != nil {
		t.Fatalf("expected no dispatch error, got %v", err)
	}

	if received.Count != 4 {
		t.Fatalf("expected subscriber to receive count %d, got %d", 4, received.Count)
	}
}

func TestStoreDispatch_ShouldReturnError_WhenStoreIsClosed(t *testing.T) {
	// Arrange
	logger := logs.NewLogger()
	logger.Mute()
	store := redux.NewStore(counterState{Count: 0}, logger)
	store.Close()

	incrementAction := redux.NewAction[int]("counter/closed")

	// Act
	err := store.Dispatch(incrementAction.With(1))

	// Assert
	if err == nil {
		t.Fatalf("expected dispatch to fail when store is closed")
	}
}
