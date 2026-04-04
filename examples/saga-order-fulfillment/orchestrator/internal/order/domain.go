package order

import (
	"time"
)

// Order represents an order in the saga
type Order struct {
	ID            string    `json:"id"`
	CustomerID    string    `json:"customerId"`
	ProductID     string    `json:"productId"`
	Quantity      int       `json:"quantity"`
	Amount        float64   `json:"amount"`
	Address       string    `json:"address"`
	Status        string    `json:"status"`
	PaymentID     string    `json:"paymentId,omitempty"`
	ReservationID string    `json:"reservationId,omitempty"`
	ShipmentID    string    `json:"shipmentId,omitempty"`
	CorrelationID string    `json:"correlationId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ErrorMessage  string    `json:"errorMessage,omitempty"`
}

// OrderStatus constants
const (
	StatusPending    = "PENDING"
	StatusProcessing = "PROCESSING"
	StatusCompleted  = "COMPLETED"
	StatusFailed     = "FAILED"
	StatusRolledBack = "ROLLED_BACK"
)

// SagaStep represents a step in the saga
type SagaStep string

const (
	StepPayment   SagaStep = "PAYMENT"
	StepInventory SagaStep = "INVENTORY"
	StepShipping  SagaStep = "SHIPPING"
)
