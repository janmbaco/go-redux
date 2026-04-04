package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/orchestrator/internal/clients"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/orchestrator/internal/handlers"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/orchestrator/internal/order"
	redux "github.com/janmbaco/go-redux/v2"
)

// OrchestratorModule implements Module for orchestrator dependencies
type OrchestratorModule struct {
	paymentURL   string
	inventoryURL string
	shippingURL  string
}

// NewOrchestratorModule creates a new orchestrator module
func NewOrchestratorModule(paymentURL, inventoryURL, shippingURL string) *OrchestratorModule {
	return &OrchestratorModule{
		paymentURL:   paymentURL,
		inventoryURL: inventoryURL,
		shippingURL:  shippingURL,
	}
}

// RegisterServices registers orchestrator services
func (m *OrchestratorModule) RegisterServices(register dependencyinjection.Register) error {
	// Register HTTP clients
	register.AsSingleton(
		new(*clients.PaymentClient),
		func(logger logs.Logger) (*clients.PaymentClient, error) {
			return clients.NewPaymentClient(m.paymentURL, logger), nil
		},
		map[int]string{0: "logger"},
	)

	register.AsSingleton(
		new(*clients.InventoryClient),
		func(logger logs.Logger) (*clients.InventoryClient, error) {
			return clients.NewInventoryClient(m.inventoryURL, logger), nil
		},
		map[int]string{0: "logger"},
	)

	register.AsSingleton(
		new(*clients.ShippingClient),
		func(logger logs.Logger) (*clients.ShippingClient, error) {
			return clients.NewShippingClient(m.shippingURL, logger), nil
		},
		map[int]string{0: "logger"},
	)

	// Register SagaOrchestrator
	register.AsSingleton(
		new(*order.SagaOrchestrator),
		func(
			store redux.Store[map[string]any],
			paymentClient *clients.PaymentClient,
			inventoryClient *clients.InventoryClient,
			shippingClient *clients.ShippingClient,
			logger logs.Logger,
		) (*order.SagaOrchestrator, error) {
			return order.NewSagaOrchestrator(store, paymentClient, inventoryClient, shippingClient, logger), nil
		},
		map[int]string{0: "store", 1: "paymentClient", 2: "inventoryClient", 3: "shippingClient", 4: "logger"},
	)

	// Register OrderHandler
	register.AsSingleton(
		new(*handlers.OrderHandler),
		func(orchestrator *order.SagaOrchestrator, logger logs.Logger) (*handlers.OrderHandler, error) {
			return handlers.NewOrderHandler(orchestrator, logger), nil
		},
		map[int]string{0: "orchestrator", 1: "logger"},
	)

	return nil
}
