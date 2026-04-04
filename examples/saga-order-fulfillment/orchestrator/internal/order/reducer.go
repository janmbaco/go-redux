package order

// OrderState represents the application state
type OrderState struct {
	Orders          map[string]*Order `json:"orders"`
	TotalOrders     int               `json:"totalOrders"`
	CompletedOrders int               `json:"completedOrders"`
	FailedOrders    int               `json:"failedOrders"`
	RolledBack      int               `json:"rolledBack"`
}

// NewOrderState creates a new order state
func NewOrderState() OrderState {
	return OrderState{
		Orders: make(map[string]*Order),
	}
}

// Pure reducers - each handles one use case

func CreateOrderReducer(state OrderState, order Order) OrderState {
	newState := state
	newState.Orders = copyOrders(state.Orders)
	newState.Orders[order.ID] = &order
	newState.TotalOrders++
	return newState
}

func StartSagaReducer(state OrderState, orderID string) OrderState {
	newState := state
	newState.Orders = copyOrders(state.Orders)
	if order, exists := newState.Orders[orderID]; exists {
		updatedOrder := *order
		updatedOrder.Status = StatusProcessing
		newState.Orders[orderID] = &updatedOrder
	}
	return newState
}

func StepCompletedReducer(state OrderState, payload StepCompletedPayload) OrderState {
	newState := state
	newState.Orders = copyOrders(state.Orders)
	if order, exists := newState.Orders[payload.OrderID]; exists {
		updatedOrder := *order
		switch payload.Step {
		case StepPayment:
			updatedOrder.PaymentID = payload.StepID
		case StepInventory:
			updatedOrder.ReservationID = payload.StepID
		case StepShipping:
			updatedOrder.ShipmentID = payload.StepID
		}
		newState.Orders[payload.OrderID] = &updatedOrder
	}
	return newState
}

func StepFailedReducer(state OrderState, payload StepFailedPayload) OrderState {
	newState := state
	newState.Orders = copyOrders(state.Orders)
	if order, exists := newState.Orders[payload.OrderID]; exists {
		updatedOrder := *order
		updatedOrder.ErrorMessage = payload.Error
		newState.Orders[payload.OrderID] = &updatedOrder
	}
	return newState
}

func CompleteOrderReducer(state OrderState, orderID string) OrderState {
	newState := state
	newState.Orders = copyOrders(state.Orders)
	if order, exists := newState.Orders[orderID]; exists {
		updatedOrder := *order
		updatedOrder.Status = StatusCompleted
		newState.Orders[orderID] = &updatedOrder
		newState.CompletedOrders++
	}
	return newState
}

func FailOrderReducer(state OrderState, payload FailOrderPayload) OrderState {
	newState := state
	newState.Orders = copyOrders(state.Orders)
	if order, exists := newState.Orders[payload.OrderID]; exists {
		updatedOrder := *order
		updatedOrder.Status = StatusFailed
		updatedOrder.ErrorMessage = payload.Error
		newState.Orders[payload.OrderID] = &updatedOrder
		newState.FailedOrders++
	}
	return newState
}

func RollbackOrderReducer(state OrderState, orderID string) OrderState {
	newState := state
	newState.Orders = copyOrders(state.Orders)
	if order, exists := newState.Orders[orderID]; exists {
		updatedOrder := *order
		updatedOrder.Status = StatusRolledBack
		newState.Orders[orderID] = &updatedOrder
		newState.RolledBack++
	}
	return newState
}

// copyOrders creates a shallow copy of the orders map
func copyOrders(orders map[string]*Order) map[string]*Order {
	newOrders := make(map[string]*Order, len(orders))
	for k, v := range orders {
		newOrders[k] = v
	}
	return newOrders
}
