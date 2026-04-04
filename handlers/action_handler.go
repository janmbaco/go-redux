package handlers

import (
	"fmt"
	"reflect"

	"github.com/janmbaco/go-redux/v2/actions"
)

// ActionHandler represents a unit that encapsulates business logic for a specific domain.
// ActionHandlers connect actions to handlers and maintain their own state slice.
//
// Example:
//
//	builder := redux.NewActionHandlerBuilder[CounterState]()
//	builder.On(IncrementAction, incrementHandler)
//	builder.SetSelector("counter")
//	handler := builder.Build()
type ActionHandler[S any] interface {
	// GetInitialState returns the initial state for this handler
	GetInitialState() S

	// Handle processes an action and returns the new state
	Handle(state S, action any) (S, error)

	// GetSelector returns the selector key for this handler's state slice
	GetSelector() string

	// CanHandle checks if this handler can handle the given action
	CanHandle(action any) bool
}

// actionHandler is the concrete implementation of ActionHandler
type actionHandler[S any] struct {
	initialState S
	handlers     map[string]any
	selector     string
}

// ActionHandlerBuilder builds an ActionHandler with action handlers.
//
// Example:
//
//	builder := redux.NewActionHandlerBuilder[int]()
//	builder.SetInitialState(0)
//	builder.On(IncrementAction, func(state int, payload int) (int, error) {
//	    return state + payload, nil
//	})
//	builder.SetSelector("counter")
//	handler := builder.Build()
type ActionHandlerBuilder[S any] struct {
	initialState S
	handlers     map[string]any
	selector     string
}

// NewActionHandlerBuilder creates a new ActionHandlerBuilder for the given state type
func NewActionHandlerBuilder[S any]() *ActionHandlerBuilder[S] {
	return &ActionHandlerBuilder[S]{
		handlers: make(map[string]any),
	}
}

// SetInitialState sets the initial state for this handler
func (mb *ActionHandlerBuilder[S]) SetInitialState(state S) *ActionHandlerBuilder[S] {
	mb.initialState = state
	return mb
}

// On registers a type-safe handler for the given action with payload.
// The action type parameter P ensures compile-time type safety.
//
// Example:
//
//	builder.On(IncrementAction, func(state int, amount int) int {
//	    return state + amount
//	})
func (mb *ActionHandlerBuilder[S]) On(action any, reducer any) *ActionHandlerBuilder[S] {
	reducerValue := reflect.ValueOf(reducer)

	// Extract action type
	var actionType string
	switch a := action.(type) {
	case *actions.Action[int]:
		actionType = a.Type()
	case *actions.Action[string]:
		actionType = a.Type()
	case *actions.Action[struct{}]:
		actionType = a.Type()
	default:
		// Use reflection to get Type() method
		typeMethod := reflect.ValueOf(action).MethodByName("Type")
		if !typeMethod.IsValid() {
			panic(fmt.Sprintf("action must have Type() method, got %T", action))
		}
		results := typeMethod.Call(nil)
		if len(results) != 1 {
			panic("action.Type() must return exactly one value")
		}
		actionType = results[0].String()
	}

	// Validate reducer signature
	reducerType := reducerValue.Type()
	if reducerType.Kind() != reflect.Func {
		panic(fmt.Sprintf("reducer must be a function, got %v", reducerType))
	}

	// Check if it's a simple reducer (1 param) or payload reducer (2 params)
	numIn := reducerType.NumIn()
	if numIn != 1 && numIn != 2 {
		panic(fmt.Sprintf("reducer must accept 1 parameter (state) or 2 parameters (state, payload), got %d", numIn))
	}

	numOut := reducerType.NumOut()
	if numOut != 1 && numOut != 2 {
		panic(fmt.Sprintf("reducer must return 1 value (state) or 2 values (state, error), got %d", numOut))
	}

	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if numOut == 2 && !reducerType.Out(1).Implements(errorType) {
		panic("reducer second return value must implement error")
	}

	// Wrap reducer to match internal handler signature (state, payload) -> (state, error)
	if numIn == 1 {
		// Simple reducer: func(S) S
		wrapper := func(state S, payload any) (S, error) {
			results := reducerValue.Call([]reflect.Value{reflect.ValueOf(state)})
			return extractReducerResults[S](results)
		}
		mb.handlers[actionType] = wrapper
	} else {
		// Payload reducer: func(S, P) S
		wrapper := func(state S, payload any) (S, error) {
			results := reducerValue.Call([]reflect.Value{
				reflect.ValueOf(state),
				buildReducerArgument(reducerType.In(1), payload),
			})
			return extractReducerResults[S](results)
		}
		mb.handlers[actionType] = wrapper
	}

	return mb
}

// SetSelector sets the selector key for this handler's state slice in the global state
func (mb *ActionHandlerBuilder[S]) SetSelector(selector string) *ActionHandlerBuilder[S] {
	mb.selector = selector
	return mb
}

// Build creates the ActionHandler from the builder configuration
func (mb *ActionHandlerBuilder[S]) Build() ActionHandler[S] {
	return &actionHandler[S]{
		initialState: mb.initialState,
		handlers:     mb.handlers,
		selector:     mb.selector,
	}
}

// GetInitialState returns the initial state
func (m *actionHandler[S]) GetInitialState() S {
	return m.initialState
}

// Handle processes an action and returns the new state
func (m *actionHandler[S]) Handle(state S, action any) (S, error) {
	actionType := getActionType(action)
	if actionType == "" {
		return state, fmt.Errorf("invalid action: no type")
	}

	handler, ok := m.handlers[actionType]
	if !ok {
		// ActionHandler doesn't handle this action, return state unchanged
		return state, nil
	}

	// Get payload from action instance
	payload := getActionPayload(action)

	typedHandler, ok := handler.(func(S, any) (S, error))
	if !ok {
		return state, fmt.Errorf("invalid handler signature for action %q", actionType)
	}

	return typedHandler(state, payload)
}

// GetSelector returns the selector key
func (m *actionHandler[S]) GetSelector() string {
	return m.selector
}

// CanHandle checks if this handler can handle the given action
func (m *actionHandler[S]) CanHandle(action any) bool {
	actionType := getActionType(action)
	_, ok := m.handlers[actionType]
	return ok
}

// Helper function to extract action type from action instance
func getActionType(action any) string {
	if action == nil {
		return ""
	}

	typeMethod := reflect.ValueOf(action).MethodByName("Type")
	if !typeMethod.IsValid() {
		return ""
	}
	results := typeMethod.Call(nil)
	if len(results) != 1 {
		return ""
	}
	return results[0].String()
}

// Helper function to extract payload from action instance
func getActionPayload(action any) any {
	if action == nil {
		return nil
	}

	payloadMethod := reflect.ValueOf(action).MethodByName("GetPayload")
	if !payloadMethod.IsValid() {
		return nil
	}
	results := payloadMethod.Call(nil)
	if len(results) != 1 {
		return nil
	}
	return results[0].Interface()
}

func buildReducerArgument(expectedType reflect.Type, payload any) reflect.Value {
	if payload == nil {
		return reflect.Zero(expectedType)
	}

	value := reflect.ValueOf(payload)
	if value.Type().AssignableTo(expectedType) {
		return value
	}

	if value.Type().ConvertibleTo(expectedType) {
		return value.Convert(expectedType)
	}

	return value
}

func extractReducerResults[S any](results []reflect.Value) (S, error) {
	state := results[0].Interface().(S)
	if len(results) == 1 || results[1].IsNil() {
		return state, nil
	}

	return state, results[1].Interface().(error)
}
