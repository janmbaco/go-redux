package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/orchestrator/internal/order"
)

// OrderHandler handles HTTP requests for order operations
type OrderHandler struct {
	orchestrator *order.SagaOrchestrator
	logger       logs.Logger
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orchestrator *order.SagaOrchestrator, logger logs.Logger) *OrderHandler {
	return &OrderHandler{
		orchestrator: orchestrator,
		logger:       logger,
	}
}

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	CustomerID string  `json:"customerId"`
	ProductID  string  `json:"productId"`
	Quantity   int     `json:"quantity"`
	Amount     float64 `json:"amount"`
	Address    string  `json:"address"`
}

// CreateOrder handles POST /orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var request CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error(fmt.Sprintf("[Orchestrator] Invalid request body: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate request
	if request.CustomerID == "" || request.ProductID == "" || request.Quantity <= 0 || request.Amount <= 0 || request.Address == "" {
		h.logger.Error("[Orchestrator] Missing required fields")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing required fields"})
		return
	}

	h.logger.Info(fmt.Sprintf("[Orchestrator] Creating order for customer %s: %d x %s = $%.2f",
		request.CustomerID, request.Quantity, request.ProductID, request.Amount))

	// Create order (saga executes asynchronously)
	newOrder, err := h.orchestrator.CreateOrder(
		request.CustomerID,
		request.ProductID,
		request.Quantity,
		request.Amount,
		request.Address,
	)
	if err != nil {
		h.logger.Error(fmt.Sprintf("[Orchestrator] Failed to create order: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	h.logger.Info(fmt.Sprintf("[Orchestrator] Order %s created successfully (saga started)", newOrder.ID))
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(newOrder)
}

// GetOrder handles GET /orders/:id
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("id")
	if orderID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order ID is required"})
		return
	}

	orderData := h.orchestrator.GetOrder(orderID)
	if orderData == nil {
		h.logger.Warning(fmt.Sprintf("[Orchestrator] Order %s not found", orderID))
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Order not found: %s", orderID)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderData)
}

// GetAllOrders handles GET /orders
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders := h.orchestrator.GetAllOrders()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"orders": orders,
		"count":  len(orders),
	})
}

// GetStatistics handles GET /stats
func (h *OrderHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	stats := h.orchestrator.GetStatistics()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

// HealthCheck handles GET /health
func (h *OrderHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "orchestrator",
	})
}
