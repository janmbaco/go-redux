package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/services/payment-service/internal/domain"
)

// PaymentHandler handles payment HTTP requests
type PaymentHandler struct {
	dataAccess dataaccess.DataAccess
	logger     logs.Logger
	rng        *rand.Rand
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(
	dataAccess dataaccess.DataAccess,
	logger logs.Logger,
) *PaymentHandler {
	return &PaymentHandler{
		dataAccess: dataAccess,
		logger:     logger,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ProcessPaymentRequest represents the request body
type ProcessPaymentRequest struct {
	OrderID       string  `json:"order_id"`
	CustomerID    string  `json:"customer_id"`
	Amount        float64 `json:"amount"`
	CorrelationID string  `json:"correlation_id"`
}

// ProcessPayment handles POST /payments (30% random failure)
func (h *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var paymentReq ProcessPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentReq); err != nil {
		h.logger.Error(fmt.Sprintf("[Payment] Invalid request: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	h.logger.Info(fmt.Sprintf("[Payment:%s] Processing payment for order %s (amount: %.2f)",
		paymentReq.CorrelationID, paymentReq.OrderID, paymentReq.Amount))

	// Simulate processing delay
	time.Sleep(time.Millisecond * time.Duration(100+h.rng.Intn(400)))

	// 30% random failure rate
	if h.rng.Float64() < 0.30 {
		payment := &domain.Payment{
			ID:            uuid.New().String(),
			OrderID:       paymentReq.OrderID,
			CustomerID:    paymentReq.CustomerID,
			Amount:        paymentReq.Amount,
			Status:        domain.PaymentFailed,
			CorrelationID: paymentReq.CorrelationID,
			FailureReason: "insufficient funds",
		}

		if err := dataaccess.InsertRow(h.dataAccess, payment); err != nil {
			h.logger.Error(fmt.Sprintf("[Payment:%s] DB error: %v", paymentReq.CorrelationID, err))
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
			return
		}

		h.logger.Error(fmt.Sprintf("[Payment:%s] Payment failed: insufficient funds", paymentReq.CorrelationID))
		w.WriteHeader(http.StatusPaymentRequired)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "payment failed",
			"reason":  "insufficient funds",
			"payment": payment,
		})
		return
	}

	// Create successful payment
	payment := &domain.Payment{
		ID:            uuid.New().String(),
		OrderID:       paymentReq.OrderID,
		CustomerID:    paymentReq.CustomerID,
		Amount:        paymentReq.Amount,
		Status:        domain.PaymentCompleted,
		CorrelationID: paymentReq.CorrelationID,
	}

	if err := dataaccess.InsertRow(h.dataAccess, payment); err != nil {
		h.logger.Error(fmt.Sprintf("[Payment:%s] DB error: %v", paymentReq.CorrelationID, err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		return
	}

	h.logger.Info(fmt.Sprintf("[Payment:%s] Payment completed: %s", paymentReq.CorrelationID, payment.ID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

// RefundPayment handles POST /payments/:id/refund (compensation)
func (h *PaymentHandler) RefundPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("id")

	filter := &domain.Payment{ID: paymentID}
	payments, err := dataaccess.SelectRows[domain.Payment](h.dataAccess, filter)
	if err != nil || len(payments) == 0 {
		h.logger.Error(fmt.Sprintf("[Payment] Payment not found: %s", paymentID))
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "payment not found"})
		return
	}

	payment := payments[0]
	h.logger.Info(fmt.Sprintf("[Payment:%s] Refunding payment: %s", payment.CorrelationID, paymentID))

	// Simulate processing delay
	time.Sleep(time.Millisecond * time.Duration(50+h.rng.Intn(200)))

	payment.MarkAsRefunded()
	update := &domain.Payment{Status: payment.Status}
	if err := dataaccess.UpdateRow(h.dataAccess, filter, update); err != nil {
		h.logger.Error(fmt.Sprintf("[Payment:%s] DB error on refund: %v", payment.CorrelationID, err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		return
	}

	h.logger.Info(fmt.Sprintf("[Payment:%s] Refund completed: %s", payment.CorrelationID, paymentID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

// GetPayment handles GET /payments/:id
func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("id")

	filter := &domain.Payment{ID: paymentID}
	payments, err := dataaccess.SelectRows[domain.Payment](h.dataAccess, filter)
	if err != nil || len(payments) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "payment not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payments[0])
}

// GetPaymentByOrderID handles GET /payments/order/:orderId
func (h *PaymentHandler) GetPaymentByOrderID(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("orderId")

	filter := &domain.Payment{OrderID: orderID}
	payments, err := dataaccess.SelectRows[domain.Payment](h.dataAccess, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payments)
}

// HealthCheck handles GET /health
func (h *PaymentHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "payment-service",
	})
}
