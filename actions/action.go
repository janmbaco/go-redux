// Package actions provides action types for Redux.
//
// Actions represent events that trigger state changes in the Redux store.
//
// Basic usage:
//
//	// Define actions
//	var IncrementAction = actions.NewAction[int]("INCREMENT")
//
//	// Dispatch with payload
//	store.Dispatch(IncrementAction.With(5))
package actions

import (
	"fmt"
	"sync"
)

// Action represents a dispatchable event in the Redux system.
// Actions are identified by their type string and carry a typed payload.
// Actions should be declared as package-level variables for identity.
//
// Example:
//
//	var AddTodoAction = redux.NewAction[Todo]("ADD_TODO")
//	var IncrementAction = redux.NewAction[int]("INCREMENT")
type Action[P any] struct {
	actionType string
	payload    *P
	mu         sync.RWMutex
}

var actionRegistry = &struct {
	mu      sync.RWMutex
	actions map[string]any
}{
	actions: make(map[string]any),
}

// NewAction creates a new Action with the given type.
// The action is registered globally and can be used to dispatch payloads.
//
// Example:
//
//	var MoveAction = redux.NewAction[MovePayload]("PLAYER_MOVE")
func NewAction[P any](actionType string) *Action[P] {
	action := &Action[P]{
		actionType: actionType,
	}

	actionRegistry.mu.Lock()
	actionRegistry.actions[actionType] = action
	actionRegistry.mu.Unlock()

	return action
}

// With creates an action instance with the given payload.
// This is the method used to dispatch actions with data.
//
// Example:
//
//	store.Dispatch(IncrementAction.With(5))
//	store.Dispatch(AddTodoAction.With(Todo{ID: 1, Text: "Buy milk"}))
func (a *Action[P]) With(payload P) *ActionInstance[P] {
	return &ActionInstance[P]{
		action:  a,
		payload: payload,
	}
}

// Type returns the action type string
func (a *Action[P]) Type() string {
	return a.actionType
}

// ActionInstance represents a dispatched action with its payload.
// This is what gets passed through the Redux system.
type ActionInstance[P any] struct {
	action  *Action[P]
	payload P
}

// GetAction returns the action template
func (ai *ActionInstance[P]) GetAction() *Action[P] {
	return ai.action
}

// GetPayload returns the action payload
func (ai *ActionInstance[P]) GetPayload() P {
	return ai.payload
}

// Type returns the action type string
func (ai *ActionInstance[P]) Type() string {
	return ai.action.actionType
}

// String returns a string representation of the action instance
func (ai *ActionInstance[P]) String() string {
	return fmt.Sprintf("Action{type: %s, payload: %v}", ai.action.actionType, ai.payload)
}

// GetActionByType retrieves a registered action by its type string
func GetActionByType(actionType string) (any, bool) {
	actionRegistry.mu.RLock()
	defer actionRegistry.mu.RUnlock()
	action, ok := actionRegistry.actions[actionType]
	return action, ok
}
