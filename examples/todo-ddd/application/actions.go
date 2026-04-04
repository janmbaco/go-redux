package application

import "github.com/janmbaco/go-redux/v2/actions"

// Todo actions
var (
	AddTodoAction        = actions.NewAction[string]("ADD_TODO")
	ToggleTodoAction     = actions.NewAction[int]("TOGGLE_TODO")
	RemoveTodoAction     = actions.NewAction[int]("REMOVE_TODO")
	ClearCompletedAction = actions.NewAction[struct{}]("CLEAR_COMPLETED")
)

// Filter actions
var (
	SetFilterAction = actions.NewAction[string]("SET_FILTER")
)
