package order

import (
	"fmt"
	"time"

	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/orchestrator/internal/clients"
	redux "github.com/janmbaco/go-redux/v2"
)

// SagaOrchestrator coordinates the saga execution and compensation
type SagaOrchestrator struct {
	store           redux.Store[map[string]any]
	paymentClient   *clients.PaymentClient
	inventoryClient *clients.InventoryClient
	shippingClient  *clients.ShippingClient
	logger          logs.Logger
}

// NewSagaOrchestrator creates a new saga orchestrator
func NewSagaOrchestrator(
	store redux.Store[map[string]any],
	paymentClient *clients.PaymentClient,
	inventoryClient *clients.InventoryClient,
	shippingClient *clients.ShippingClient,
	logger logs.Logger,
) *SagaOrchestrator {
	return &SagaOrchestrator{
		store:           store,
		paymentClient:   paymentClient,
		inventoryClient: inventoryClient,
		shippingClient:  shippingClient,
		logger:          logger,
	}
}

// CreateOrder creates a new order and starts the saga
func (s *SagaOrchestrator) CreateOrder(customerID, productID string, quantity int, amount float64, address string) (*Order, error) {
	order := &Order{
		ID:            generateID(),
		CustomerID:    customerID,
		ProductID:     productID,
		Quantity:      quantity,
		Amount:        amount,
		Address:       address,
		Status:        StatusPending,
		CorrelationID: generateID(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Dispatch CREATE_ORDER action
	s.store.Dispatch(CreateOrderAction.With(*order))

	s.logger.Info(fmt.Sprintf("[Saga] Order %s created (correlation: %s)", order.ID, order.CorrelationID))

	// Start saga execution
	go s.ExecuteSaga(order.ID)

	return order, nil
}

// ExecuteSaga executes the saga steps
func (s *SagaOrchestrator) ExecuteSaga(orderID string) {
	s.logger.Info(fmt.Sprintf("[Saga] Starting saga for order %s", orderID))

	// Dispatch START_SAGA action
	s.store.Dispatch(StartSagaAction.With(orderID))

	order := s.getOrder(orderID)
	if order == nil {
		s.logger.Error(fmt.Sprintf("[Saga] Order %s not found", orderID))
		return
	}

	// Step 1: Process Payment
	s.logger.Info(fmt.Sprintf("[Saga] Step 1/3 - Processing payment for order %s", orderID))
	payment, err := s.paymentClient.ProcessPayment(order.ID, order.CustomerID, order.Amount, order.CorrelationID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("[Saga] Payment failed for order %s: %v", orderID, err))
		s.handleStepFailure(orderID, StepPayment, err.Error())
		return
	}

	s.handleStepSuccess(orderID, StepPayment, payment.ID)
	s.logger.Info(fmt.Sprintf("[Saga] Payment %s completed for order %s", payment.ID, orderID))

	// Step 2: Reserve Inventory
	s.logger.Info(fmt.Sprintf("[Saga] Step 2/3 - Reserving inventory for order %s", orderID))
	reservation, err := s.inventoryClient.ReserveInventory(order.ID, order.ProductID, order.Quantity, order.CorrelationID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("[Saga] Inventory reservation failed for order %s: %v", orderID, err))
		s.handleStepFailure(orderID, StepInventory, err.Error())
		s.compensate(orderID, StepInventory)
		return
	}

	s.handleStepSuccess(orderID, StepInventory, reservation.ID)
	s.logger.Info(fmt.Sprintf("[Saga] Reservation %s completed for order %s", reservation.ID, orderID))

	// Step 3: Create Shipment
	s.logger.Info(fmt.Sprintf("[Saga] Step 3/3 - Creating shipment for order %s", orderID))
	shipment, err := s.shippingClient.CreateShipment(order.ID, order.CustomerID, order.Address, order.CorrelationID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("[Saga] Shipment creation failed for order %s: %v", orderID, err))
		s.handleStepFailure(orderID, StepShipping, err.Error())
		s.compensate(orderID, StepShipping)
		return
	}

	s.handleStepSuccess(orderID, StepShipping, shipment.ID)
	s.logger.Info(fmt.Sprintf("[Saga] Shipment %s created with tracking %s for order %s", shipment.ID, shipment.TrackingCode, orderID))

	// All steps completed successfully
	s.completeOrder(orderID)
	s.logger.Info(fmt.Sprintf("[Saga] Order %s completed successfully ✓", orderID))
}

// compensate executes compensation logic (reverse order)
func (s *SagaOrchestrator) compensate(orderID string, failedStep SagaStep) {
	s.logger.Warning(fmt.Sprintf("[Saga] Starting compensation for order %s (failed at %s)", orderID, failedStep))

	order := s.getOrder(orderID)
	if order == nil {
		s.logger.Error(fmt.Sprintf("[Saga] Order %s not found for compensation", orderID))
		return
	}

	// Compensate in reverse order
	switch failedStep {
	case StepShipping:
		// Failed at shipping, need to release inventory and refund payment
		if order.ReservationID != "" {
			s.logger.Info(fmt.Sprintf("[Saga] Compensating: Releasing inventory reservation %s", order.ReservationID))
			if err := s.inventoryClient.ReleaseReservation(order.ReservationID, order.CorrelationID); err != nil {
				s.logger.Error(fmt.Sprintf("[Saga] Failed to release reservation: %v", err))
			}
		}
		if order.PaymentID != "" {
			s.logger.Info(fmt.Sprintf("[Saga] Compensating: Refunding payment %s", order.PaymentID))
			if err := s.paymentClient.RefundPayment(order.PaymentID, order.CorrelationID); err != nil {
				s.logger.Error(fmt.Sprintf("[Saga] Failed to refund payment: %v", err))
			}
		}

	case StepInventory:
		// Failed at inventory, need to refund payment
		if order.PaymentID != "" {
			s.logger.Info(fmt.Sprintf("[Saga] Compensating: Refunding payment %s", order.PaymentID))
			if err := s.paymentClient.RefundPayment(order.PaymentID, order.CorrelationID); err != nil {
				s.logger.Error(fmt.Sprintf("[Saga] Failed to refund payment: %v", err))
			}
		}

	case StepPayment:
		// Failed at payment, no compensation needed
		s.logger.Info(fmt.Sprintf("[Saga] No compensation needed for payment failure on order %s", orderID))
	}

	// Dispatch ROLLBACK_ORDER action
	s.store.Dispatch(RollbackOrderAction.With(orderID))

	s.logger.Warning(fmt.Sprintf("[Saga] Order %s rolled back successfully ↺", orderID))
}

// handleStepSuccess dispatches STEP_COMPLETED action
func (s *SagaOrchestrator) handleStepSuccess(orderID string, step SagaStep, stepID string) {
	s.store.Dispatch(StepCompletedAction.With(StepCompletedPayload{
		OrderID: orderID,
		Step:    step,
		StepID:  stepID,
	}))
}

// handleStepFailure dispatches STEP_FAILED action
func (s *SagaOrchestrator) handleStepFailure(orderID string, step SagaStep, errorMsg string) {
	s.store.Dispatch(StepFailedAction.With(StepFailedPayload{
		OrderID: orderID,
		Step:    step,
		Error:   errorMsg,
	}))

	// Dispatch FAIL_ORDER action
	s.store.Dispatch(FailOrderAction.With(FailOrderPayload{
		OrderID: orderID,
		Error:   errorMsg,
	}))
}

// completeOrder dispatches COMPLETE_ORDER action
func (s *SagaOrchestrator) completeOrder(orderID string) {
	s.store.Dispatch(CompleteOrderAction.With(orderID))
}

// getOrder retrieves an order from the store
func (s *SagaOrchestrator) getOrder(orderID string) *Order {
	state := s.store.GetState()
	orderState := state["orders"].(OrderState)
	return orderState.Orders[orderID]
}

// GetOrder returns an order by ID
func (s *SagaOrchestrator) GetOrder(orderID string) *Order {
	return s.getOrder(orderID)
}

// GetAllOrders returns all orders
func (s *SagaOrchestrator) GetAllOrders() []*Order {
	state := s.store.GetState()
	orderState := state["orders"].(OrderState)
	orders := make([]*Order, 0, len(orderState.Orders))
	for _, order := range orderState.Orders {
		orders = append(orders, order)
	}
	return orders
}

// GetStatistics returns saga statistics
func (s *SagaOrchestrator) GetStatistics() map[string]interface{} {
	state := s.store.GetState()
	orderState := state["orders"].(OrderState)
	return map[string]interface{}{
		"total_orders":     orderState.TotalOrders,
		"completed_orders": orderState.CompletedOrders,
		"failed_orders":    orderState.FailedOrders,
		"rolled_back":      orderState.RolledBack,
		"success_rate":     calculateSuccessRate(orderState.CompletedOrders, orderState.TotalOrders),
	}
}

func calculateSuccessRate(completed, total int) float64 {
	if total == 0 {
		return 0.0
	}
	return float64(completed) / float64(total) * 100.0
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
