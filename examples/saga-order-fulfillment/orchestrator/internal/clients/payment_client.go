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

// PaymentClient handles communication with payment service
type PaymentClient struct {
	baseURL string
	client  *http.Client
	logger  logs.Logger
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	OrderID       string  `json:"orderId"`
	CustomerID    string  `json:"customerId"`
	Amount        float64 `json:"amount"`
	CorrelationID string  `json:"correlationId"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"orderId"`
	CustomerID    string    `json:"customerId"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	CorrelationID string    `json:"correlationId"`
	CreatedAt     time.Time `json:"createdAt"`
}

// NewPaymentClient creates a new payment service client
func NewPaymentClient(baseURL string, logger logs.Logger) *PaymentClient {
	return &PaymentClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// ProcessPayment sends a payment request to the payment service
func (c *PaymentClient) ProcessPayment(orderID, customerID string, amount float64, correlationID string) (*PaymentResponse, error) {
	req := PaymentRequest{
		OrderID:       orderID,
		CustomerID:    customerID,
		Amount:        amount,
		CorrelationID: correlationID,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment request: %w", err)
	}

	c.logger.Info(fmt.Sprintf("[PaymentClient] Processing payment for order %s (correlation: %s)", orderID, correlationID))

	resp, err := c.client.Post(c.baseURL+"/payments", "application/json", bytes.NewBuffer(data))
	if err != nil {
		c.logger.Error(fmt.Sprintf("[PaymentClient] Failed to process payment: %v", err))
		return nil, fmt.Errorf("payment service error: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.logger.Warning(fmt.Sprintf("[PaymentClient] Payment failed with status %d: %s", resp.StatusCode, string(body)))
		return nil, fmt.Errorf("payment failed: %s", string(body))
	}

	var payment PaymentResponse
	if err := json.Unmarshal(body, &payment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment response: %w", err)
	}

	c.logger.Info(fmt.Sprintf("[PaymentClient] Payment %s processed successfully", payment.ID))
	return &payment, nil
}

// RefundPayment sends a refund request to the payment service
func (c *PaymentClient) RefundPayment(paymentID, correlationID string) error {
	c.logger.Info(fmt.Sprintf("[PaymentClient] Refunding payment %s (correlation: %s)", paymentID, correlationID))

	req, _ := http.NewRequest("POST", c.baseURL+"/payments/"+paymentID+"/refund", nil)
	req.Header.Set("X-Correlation-ID", correlationID)

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(fmt.Sprintf("[PaymentClient] Failed to refund payment: %v", err))
		return fmt.Errorf("refund service error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Warning(fmt.Sprintf("[PaymentClient] Refund failed with status %d: %s", resp.StatusCode, string(body)))
		return fmt.Errorf("refund failed: %s", string(body))
	}

	c.logger.Info(fmt.Sprintf("[PaymentClient] Payment %s refunded successfully", paymentID))
	return nil
}

// GetPaymentByOrderID retrieves a payment by order ID
func (c *PaymentClient) GetPaymentByOrderID(orderID string) (*PaymentResponse, error) {
	resp, err := c.client.Get(c.baseURL + "/payments/order/" + orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment not found")
	}

	var payment PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&payment); err != nil {
		return nil, fmt.Errorf("failed to decode payment: %w", err)
	}

	return &payment, nil
}
