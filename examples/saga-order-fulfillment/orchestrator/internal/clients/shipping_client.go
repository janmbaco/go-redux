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

// ShippingClient handles communication with shipping service
type ShippingClient struct {
	baseURL string
	client  *http.Client
	logger  logs.Logger
}

// ShipmentRequest represents a shipment creation request
type ShipmentRequest struct {
	OrderID       string `json:"orderId"`
	CustomerID    string `json:"customerId"`
	Address       string `json:"address"`
	CorrelationID string `json:"correlationId"`
}

// ShipmentResponse represents a shipment response
type ShipmentResponse struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"orderId"`
	CustomerID    string    `json:"customerId"`
	Address       string    `json:"address"`
	TrackingCode  string    `json:"trackingCode"`
	Status        string    `json:"status"`
	CorrelationID string    `json:"correlationId"`
	CreatedAt     time.Time `json:"createdAt"`
}

// NewShippingClient creates a new shipping service client
func NewShippingClient(baseURL string, logger logs.Logger) *ShippingClient {
	return &ShippingClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// CreateShipment sends a shipment creation request
func (c *ShippingClient) CreateShipment(orderID, customerID, address, correlationID string) (*ShipmentResponse, error) {
	req := ShipmentRequest{
		OrderID:       orderID,
		CustomerID:    customerID,
		Address:       address,
		CorrelationID: correlationID,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal shipment request: %w", err)
	}

	c.logger.Info(fmt.Sprintf("[ShippingClient] Creating shipment for order %s (correlation: %s)", orderID, correlationID))

	resp, err := c.client.Post(c.baseURL+"/shipments", "application/json", bytes.NewBuffer(data))
	if err != nil {
		c.logger.Error(fmt.Sprintf("[ShippingClient] Failed to create shipment: %v", err))
		return nil, fmt.Errorf("shipping service error: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.logger.Warning(fmt.Sprintf("[ShippingClient] Shipment creation failed with status %d: %s", resp.StatusCode, string(body)))
		return nil, fmt.Errorf("shipment creation failed: %s", string(body))
	}

	var shipment ShipmentResponse
	if err := json.Unmarshal(body, &shipment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal shipment response: %w", err)
	}

	c.logger.Info(fmt.Sprintf("[ShippingClient] Shipment %s created successfully with tracking %s", shipment.ID, shipment.TrackingCode))
	return &shipment, nil
}

// CancelShipment sends a cancellation request to the shipping service
func (c *ShippingClient) CancelShipment(shipmentID, correlationID string) error {
	c.logger.Info(fmt.Sprintf("[ShippingClient] Cancelling shipment %s (correlation: %s)", shipmentID, correlationID))

	req, _ := http.NewRequest("POST", c.baseURL+"/shipments/"+shipmentID+"/cancel", nil)
	req.Header.Set("X-Correlation-ID", correlationID)

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(fmt.Sprintf("[ShippingClient] Failed to cancel shipment: %v", err))
		return fmt.Errorf("cancellation service error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Warning(fmt.Sprintf("[ShippingClient] Cancellation failed with status %d: %s", resp.StatusCode, string(body)))
		return fmt.Errorf("cancellation failed: %s", string(body))
	}

	c.logger.Info(fmt.Sprintf("[ShippingClient] Shipment %s cancelled successfully", shipmentID))
	return nil
}

// GetShipmentByOrderID retrieves a shipment by order ID
func (c *ShippingClient) GetShipmentByOrderID(orderID string) (*ShipmentResponse, error) {
	resp, err := c.client.Get(c.baseURL + "/shipments/order/" + orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("shipment not found")
	}

	var shipment ShipmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&shipment); err != nil {
		return nil, fmt.Errorf("failed to decode shipment: %w", err)
	}

	return &shipment, nil
}
