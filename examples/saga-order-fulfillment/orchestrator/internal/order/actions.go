package order

import "github.com/janmbaco/go-redux/v2/actions"

// Actions for saga orchestration
var (
	CreateOrderAction   = actions.NewAction[Order]("CREATE_ORDER")
	StartSagaAction     = actions.NewAction[string]("START_SAGA")
	StepCompletedAction = actions.NewAction[StepCompletedPayload]("STEP_COMPLETED")
	StepFailedAction    = actions.NewAction[StepFailedPayload]("STEP_FAILED")
	CompleteOrderAction = actions.NewAction[string]("COMPLETE_ORDER")
	FailOrderAction     = actions.NewAction[FailOrderPayload]("FAIL_ORDER")
	RollbackOrderAction = actions.NewAction[string]("ROLLBACK_ORDER")
)

// StepCompletedPayload represents the payload when a saga step completes
type StepCompletedPayload struct {
	OrderID string   `json:"orderId"`
	Step    SagaStep `json:"step"`
	StepID  string   `json:"stepId"` // Payment/Reservation/Shipment ID
}

// StepFailedPayload represents the payload when a saga step fails
type StepFailedPayload struct {
	OrderID string   `json:"orderId"`
	Step    SagaStep `json:"step"`
	Error   string   `json:"error"`
}

// FailOrderPayload represents the payload for failing an order
type FailOrderPayload struct {
	OrderID string `json:"orderId"`
	Error   string `json:"error"`
}
