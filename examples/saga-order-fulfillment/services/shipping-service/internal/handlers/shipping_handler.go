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
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/services/shipping-service/internal/domain"
)

// ShippingHandler handles shipping HTTP requests
type ShippingHandler struct {
	dataAccess dataaccess.DataAccess
	logger     logs.Logger
	rng        *rand.Rand
}

// NewShippingHandler creates a new shipping handler
func NewShippingHandler(
	dataAccess dataaccess.DataAccess,
	logger logs.Logger,
) *ShippingHandler {
	return &ShippingHandler{
		dataAccess: dataAccess,
		logger:     logger,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CreateShipmentRequest represents the request body
type CreateShipmentRequest struct {
	OrderID       string `json:"order_id"`
	Address       string `json:"address"`
	CorrelationID string `json:"correlation_id"`
}

// CreateShipment handles POST /shipments (30% random failure)
func (h *ShippingHandler) CreateShipment(w http.ResponseWriter, r *http.Request) {
	var shipmentReq CreateShipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&shipmentReq); err != nil {
		h.logger.Error(fmt.Sprintf("[Shipping] Invalid request: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	h.logger.Info(fmt.Sprintf("[Shipping:%s] Creating shipment for order %s",
		shipmentReq.CorrelationID, shipmentReq.OrderID))

	// Simulate processing delay
	time.Sleep(time.Millisecond * time.Duration(100+h.rng.Intn(400)))

	// 30% random failure rate
	if h.rng.Float64() < 0.30 {
		shipment := &domain.Shipment{
			ID:            uuid.New().String(),
			OrderID:       shipmentReq.OrderID,
			Address:       shipmentReq.Address,
			Status:        domain.ShipmentFailed,
			CorrelationID: shipmentReq.CorrelationID,
			FailureReason: "carrier unavailable",
		}

		if err := dataaccess.InsertRow(h.dataAccess, shipment); err != nil {
			h.logger.Error(fmt.Sprintf("[Shipping:%s] DB error: %v", shipmentReq.CorrelationID, err))
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
			return
		}

		h.logger.Error(fmt.Sprintf("[Shipping:%s] Shipment creation failed: carrier unavailable", shipmentReq.CorrelationID))
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":    "shipment creation failed",
			"reason":   "carrier unavailable",
			"shipment": shipment,
		})
		return
	}

	// Create successful shipment
	trackingCode := fmt.Sprintf("TRACK-%d", time.Now().UnixNano())
	shipment := &domain.Shipment{
		ID:            uuid.New().String(),
		OrderID:       shipmentReq.OrderID,
		Address:       shipmentReq.Address,
		Status:        domain.ShipmentCreated,
		CorrelationID: shipmentReq.CorrelationID,
		TrackingCode:  trackingCode,
	}

	if err := dataaccess.InsertRow(h.dataAccess, shipment); err != nil {
		h.logger.Error(fmt.Sprintf("[Shipping:%s] DB error: %v", shipmentReq.CorrelationID, err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		return
	}

	h.logger.Info(fmt.Sprintf("[Shipping:%s] Shipment created: %s (tracking: %s)", shipmentReq.CorrelationID, shipment.ID, trackingCode))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shipment)
}

// CancelShipment handles POST /shipments/:id/cancel (compensation)
func (h *ShippingHandler) CancelShipment(w http.ResponseWriter, r *http.Request) {
	shipmentID := r.PathValue("id")

	filter := &domain.Shipment{ID: shipmentID}
	shipments, err := dataaccess.SelectRows[domain.Shipment](h.dataAccess, filter)
	if err != nil || len(shipments) == 0 {
		h.logger.Error(fmt.Sprintf("[Shipping] Shipment not found: %s", shipmentID))
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "shipment not found"})
		return
	}

	shipment := shipments[0]
	h.logger.Info(fmt.Sprintf("[Shipping:%s] Cancelling shipment: %s", shipment.CorrelationID, shipmentID))

	// Simulate processing delay
	time.Sleep(time.Millisecond * time.Duration(50+h.rng.Intn(200)))

	shipment.MarkAsCancelled()
	update := &domain.Shipment{Status: shipment.Status}
	if err := dataaccess.UpdateRow(h.dataAccess, filter, update); err != nil {
		h.logger.Error(fmt.Sprintf("[Shipping:%s] DB error on cancel: %v", shipment.CorrelationID, err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		return
	}

	h.logger.Info(fmt.Sprintf("[Shipping:%s] Shipment cancelled: %s", shipment.CorrelationID, shipmentID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(shipment)
}

// GetShipment handles GET /shipments/:id
func (h *ShippingHandler) GetShipment(w http.ResponseWriter, r *http.Request) {
	shipmentID := r.PathValue("id")

	filter := &domain.Shipment{ID: shipmentID}
	shipments, err := dataaccess.SelectRows[domain.Shipment](h.dataAccess, filter)
	if err != nil || len(shipments) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "shipment not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(shipments[0])
}

// GetShipmentByOrderID handles GET /shipments/order/:orderId
func (h *ShippingHandler) GetShipmentByOrderID(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("orderId")

	filter := &domain.Shipment{OrderID: orderID}
	shipments, err := dataaccess.SelectRows[domain.Shipment](h.dataAccess, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(shipments)
}

// HealthCheck handles GET /health
func (h *ShippingHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "shipping-service",
	})
}
