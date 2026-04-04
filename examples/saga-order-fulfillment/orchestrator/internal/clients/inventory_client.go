package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/janmbaco/go-infrastructure/v2/logs"
)

// InventoryClient handles communication with inventory service
type InventoryClient struct {
	baseURL string
	client  *http.Client
	logger  logs.Logger
}

// ReservationRequest represents an inventory reservation request
type ReservationRequest struct {
	OrderID       string `json:"orderId"`
	ProductID     string `json:"productId"`
	Quantity      int    `json:"quantity"`
	CorrelationID string `json:"correlationId"`
}

// ReservationResponse represents an inventory reservation response
type ReservationResponse struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"orderId"`
	ProductID     string    `json:"productId"`
	Quantity      int       `json:"quantity"`
	Status        string    `json:"status"`
	CorrelationID string    `json:"correlationId"`
	CreatedAt     time.Time `json:"createdAt"`
}

// NewInventoryClient creates a new inventory service client
func NewInventoryClient(baseURL string, logger logs.Logger) *InventoryClient {
	return &InventoryClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// ReserveInventory sends an inventory reservation request
func (c *InventoryClient) ReserveInventory(orderID, productID string, quantity int, correlationID string) (*ReservationResponse, error) {
	req := ReservationRequest{
		OrderID:       orderID,
		ProductID:     productID,
		Quantity:      quantity,
		CorrelationID: correlationID,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal reservation request: %w", err)
	}

	c.logger.Info(fmt.Sprintf("[InventoryClient] Reserving inventory for order %s (correlation: %s)", orderID, correlationID))

	resp, err := c.client.Post(c.baseURL+"/reservations", "application/json", bytes.NewBuffer(data))
	if err != nil {
		c.logger.Error(fmt.Sprintf("[InventoryClient] Failed to reserve inventory: %v", err))
		return nil, fmt.Errorf("inventory service error: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.logger.Warning(fmt.Sprintf("[InventoryClient] Reservation failed with status %d: %s", resp.StatusCode, string(body)))
		return nil, fmt.Errorf("reservation failed: %s", string(body))
	}

	var reservation ReservationResponse
	if err := json.Unmarshal(body, &reservation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reservation response: %w", err)
	}

	c.logger.Info(fmt.Sprintf("[InventoryClient] Reservation %s created successfully", reservation.ID))
	return &reservation, nil
}

// ReleaseReservation sends a release request to the inventory service
func (c *InventoryClient) ReleaseReservation(reservationID, correlationID string) error {
	c.logger.Info(fmt.Sprintf("[InventoryClient] Releasing reservation %s (correlation: %s)", reservationID, correlationID))

	req, _ := http.NewRequest("POST", c.baseURL+"/reservations/"+reservationID+"/release", nil)
	req.Header.Set("X-Correlation-ID", correlationID)

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(fmt.Sprintf("[InventoryClient] Failed to release reservation: %v", err))
		return fmt.Errorf("release service error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Warning(fmt.Sprintf("[InventoryClient] Release failed with status %d: %s", resp.StatusCode, string(body)))
		return fmt.Errorf("release failed: %s", string(body))
	}

	c.logger.Info(fmt.Sprintf("[InventoryClient] Reservation %s released successfully", reservationID))
	return nil
}

// GetReservationByOrderID retrieves a reservation by order ID
func (c *InventoryClient) GetReservationByOrderID(orderID string) (*ReservationResponse, error) {
	resp, err := c.client.Get(c.baseURL + "/reservations/order/" + orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reservation not found")
	}

	var reservation ReservationResponse
	if err := json.NewDecoder(resp.Body).Decode(&reservation); err != nil {
		return nil, fmt.Errorf("failed to decode reservation: %w", err)
	}

	return &reservation, nil
}
