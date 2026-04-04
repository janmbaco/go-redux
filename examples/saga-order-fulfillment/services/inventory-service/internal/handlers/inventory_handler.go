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
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/services/inventory-service/internal/domain"
)

// InventoryHandler handles inventory HTTP requests
type InventoryHandler struct {
	dataAccess dataaccess.DataAccess
	logger     logs.Logger
	rng        *rand.Rand
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(
	dataAccess dataaccess.DataAccess,
	logger logs.Logger,
) *InventoryHandler {
	return &InventoryHandler{
		dataAccess: dataAccess,
		logger:     logger,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ReserveInventoryRequest represents the request body
type ReserveInventoryRequest struct {
	OrderID       string   `json:"order_id"`
	ProductIDs    []string `json:"product_ids"`
	TotalItems    int      `json:"total_items"`
	CorrelationID string   `json:"correlation_id"`
}

// ReserveInventory handles POST /reservations (30% random failure)
func (h *InventoryHandler) ReserveInventory(w http.ResponseWriter, r *http.Request) {
	var reserveReq ReserveInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&reserveReq); err != nil {
		h.logger.Error(fmt.Sprintf("[Inventory] Invalid request: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	h.logger.Info(fmt.Sprintf("[Inventory:%s] Reserving %d items for order %s",
		reserveReq.CorrelationID, reserveReq.TotalItems, reserveReq.OrderID))

	// Simulate processing delay
	time.Sleep(time.Millisecond * time.Duration(100+h.rng.Intn(400)))

	// Convert product IDs to JSON string
	productIDsJSON, _ := json.Marshal(reserveReq.ProductIDs)

	// 30% random failure rate
	if h.rng.Float64() < 0.30 {
		reservation := &domain.InventoryReservation{
			ID:            uuid.New().String(),
			OrderID:       reserveReq.OrderID,
			ProductIDs:    string(productIDsJSON),
			TotalItems:    reserveReq.TotalItems,
			Status:        domain.ReservationFailed,
			CorrelationID: reserveReq.CorrelationID,
			FailureReason: "out of stock",
		}

		if err := dataaccess.InsertRow(h.dataAccess, reservation); err != nil {
			h.logger.Error(fmt.Sprintf("[Inventory:%s] DB error: %v", reserveReq.CorrelationID, err))
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
			return
		}

		h.logger.Error(fmt.Sprintf("[Inventory:%s] Reservation failed: out of stock", reserveReq.CorrelationID))
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":       "reservation failed",
			"reason":      "out of stock",
			"reservation": reservation,
		})
		return
	}

	// Create successful reservation
	reservation := &domain.InventoryReservation{
		ID:            uuid.New().String(),
		OrderID:       reserveReq.OrderID,
		ProductIDs:    string(productIDsJSON),
		TotalItems:    reserveReq.TotalItems,
		Status:        domain.ReservationActive,
		CorrelationID: reserveReq.CorrelationID,
	}

	if err := dataaccess.InsertRow(h.dataAccess, reservation); err != nil {
		h.logger.Error(fmt.Sprintf("[Inventory:%s] DB error: %v", reserveReq.CorrelationID, err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		return
	}

	h.logger.Info(fmt.Sprintf("[Inventory:%s] Reservation completed: %s", reserveReq.CorrelationID, reservation.ID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reservation)
}

// ReleaseReservation handles POST /reservations/:id/release (compensation)
func (h *InventoryHandler) ReleaseReservation(w http.ResponseWriter, r *http.Request) {
	reservationID := r.PathValue("id")

	filter := &domain.InventoryReservation{ID: reservationID}
	reservations, err := dataaccess.SelectRows[domain.InventoryReservation](h.dataAccess, filter)
	if err != nil || len(reservations) == 0 {
		h.logger.Error(fmt.Sprintf("[Inventory] Reservation not found: %s", reservationID))
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "reservation not found"})
		return
	}

	reservation := reservations[0]
	h.logger.Info(fmt.Sprintf("[Inventory:%s] Releasing reservation: %s", reservation.CorrelationID, reservationID))

	// Simulate processing delay
	time.Sleep(time.Millisecond * time.Duration(50+h.rng.Intn(200)))

	reservation.MarkAsReleased()
	update := &domain.InventoryReservation{Status: reservation.Status}
	if err := dataaccess.UpdateRow(h.dataAccess, filter, update); err != nil {
		h.logger.Error(fmt.Sprintf("[Inventory:%s] DB error on release: %v", reservation.CorrelationID, err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		return
	}

	h.logger.Info(fmt.Sprintf("[Inventory:%s] Reservation released: %s", reservation.CorrelationID, reservationID))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservation)
}

// GetReservation handles GET /reservations/:id
func (h *InventoryHandler) GetReservation(w http.ResponseWriter, r *http.Request) {
	reservationID := r.PathValue("id")

	filter := &domain.InventoryReservation{ID: reservationID}
	reservations, err := dataaccess.SelectRows[domain.InventoryReservation](h.dataAccess, filter)
	if err != nil || len(reservations) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "reservation not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservations[0])
}

// GetReservationByOrderID handles GET /reservations/order/:orderId
func (h *InventoryHandler) GetReservationByOrderID(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("orderId")

	filter := &domain.InventoryReservation{OrderID: orderID}
	reservations, err := dataaccess.SelectRows[domain.InventoryReservation](h.dataAccess, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservations)
}

// HealthCheck handles GET /health
func (h *InventoryHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "inventory-service",
	})
}
