package redux

import "github.com/janmbaco/go-redux/v2/handlers"

// ActionHandler re-exports the handler contract from the handlers package.
type ActionHandler[S any] = handlers.ActionHandler[S]

// ActionHandlerBuilder re-exports the fluent builder used to compose handlers.
type ActionHandlerBuilder[S any] = handlers.ActionHandlerBuilder[S]

// NewActionHandlerBuilder creates a new handler builder from the root package API.
func NewActionHandlerBuilder[S any]() *handlers.ActionHandlerBuilder[S] {
	return handlers.NewActionHandlerBuilder[S]()
}
