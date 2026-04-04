package redux

import "github.com/janmbaco/go-redux/v2/actions"

// Action re-exports the generic action type from the actions package.
type Action[P any] = actions.Action[P]

// ActionInstance re-exports a dispatched action with its typed payload.
type ActionInstance[P any] = actions.ActionInstance[P]

// NewAction creates a typed action in the root package API.
func NewAction[P any](actionType string) *actions.Action[P] {
	return actions.NewAction[P](actionType)
}
