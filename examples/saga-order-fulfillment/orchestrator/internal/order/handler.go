package order

import (
	"github.com/janmbaco/go-redux/v2/handlers"
)

// NewOrderHandler creates the ActionHandler for order saga orchestration
func NewOrderHandler() handlers.ActionHandler[OrderState] {
	builder := handlers.NewActionHandlerBuilder[OrderState]()

	// Register pure reducers for each action
	builder.On(CreateOrderAction, CreateOrderReducer)
	builder.On(StartSagaAction, StartSagaReducer)
	builder.On(StepCompletedAction, StepCompletedReducer)
	builder.On(StepFailedAction, StepFailedReducer)
	builder.On(CompleteOrderAction, CompleteOrderReducer)
	builder.On(FailOrderAction, FailOrderReducer)
	builder.On(RollbackOrderAction, RollbackOrderReducer)

	builder.SetInitialState(NewOrderState())
	builder.SetSelector("orders")

	return builder.Build()
}
