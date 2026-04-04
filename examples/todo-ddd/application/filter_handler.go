package application

import (
	"fmt"

	"github.com/janmbaco/go-redux/v2/handlers"
)

// FilterHandler implements the application logic for filter operations
type FilterHandler struct{}

// NewFilterHandler creates a new FilterHandler
func NewFilterHandler() handlers.ActionHandler[AppState] {
	builder := handlers.NewActionHandlerBuilder[AppState]()

	builder.On(SetFilterAction, func(state AppState, filter string) AppState {
		filterState := state["filter"].(FilterState)
		filterState.Filter = filter
		state["filter"] = filterState

		fmt.Printf("🔍 Filter changed to: %s\n", filter)
		return state
	})

	builder.SetSelector("filter")
	builder.SetInitialState(AppState{
		"filter": FilterState{Filter: "all"},
	})

	return builder.Build()
}
