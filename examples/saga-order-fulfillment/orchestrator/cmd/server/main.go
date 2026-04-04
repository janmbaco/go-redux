package main

import (
	"fmt"
	"net/http"
	"os"

	configResolver "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc/resolver"
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/logs"
	logsioc "github.com/janmbaco/go-infrastructure/v2/logs/ioc"
	"github.com/janmbaco/go-infrastructure/v2/server"
	serverIoc "github.com/janmbaco/go-infrastructure/v2/server/ioc"
	serverResolver "github.com/janmbaco/go-infrastructure/v2/server/ioc/resolver"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/orchestrator/internal/handlers"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/orchestrator/internal/ioc"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/orchestrator/internal/order"
	reduxioc "github.com/janmbaco/go-redux/v2/ioc"
	reduxresolver "github.com/janmbaco/go-redux/v2/ioc/resolver"
)

// Config represents orchestrator configuration
type Config struct {
	Port     string `json:"port"`
	Services struct {
		PaymentURL   string `json:"paymentUrl"`
		InventoryURL string `json:"inventoryUrl"`
		ShippingURL  string `json:"shippingUrl"`
	} `json:"services"`
}

func main() {
	// Default configuration
	defaultConfig := &Config{}
	defaultConfig.Port = ":8080"
	defaultConfig.Services.PaymentURL = getEnv("PAYMENT_URL", "http://payment-service:8081")
	defaultConfig.Services.InventoryURL = getEnv("INVENTORY_URL", "http://inventory-service:8082")
	defaultConfig.Services.ShippingURL = getEnv("SHIPPING_URL", "http://shipping-service:8083")

	// Initial Redux state
	initialState := map[string]any{
		"orders": order.NewOrderState(),
	}

	// Build container with all modules
	container := dependencyinjection.NewBuilder().
		AddModule(logsioc.NewLogsModule()).
		AddModules(serverIoc.ConfigureServerModules()...).
		AddModule(reduxioc.NewReduxModule[map[string]any]()).
		AddModule(ioc.NewOrchestratorModule(
			defaultConfig.Services.PaymentURL,
			defaultConfig.Services.InventoryURL,
			defaultConfig.Services.ShippingURL,
		)).
		MustBuild()

	resolver := container.Resolver()

	// Get Redux store and add order handler
	store := reduxresolver.GetStore(resolver, initialState)
	orderHandler := order.NewOrderHandler()
	store.AddModule(orderHandler)

	// Subscribe to state changes for logging
	logger := dependencyinjection.Resolve[logs.Logger](resolver)
	stateLogger := func(state map[string]any) {
		orderState := state["orders"].(order.OrderState)
		logger.Info(fmt.Sprintf("[Redux] State updated: Total=%d, Completed=%d, Failed=%d, RolledBack=%d",
			orderState.TotalOrders, orderState.CompletedOrders, orderState.FailedOrders, orderState.RolledBack))
	}
	store.Subscribe(&stateLogger)

	// Get config handler with file monitoring
	configHandler := configResolver.GetFileConfigHandler(
		resolver,
		"orchestrator.json",
		defaultConfig,
	)

	// Get listener builder with bootstrapper
	listener, err := serverResolver.GetListenerBuilder(resolver, configHandler).
		SetBootstrapper(func(config interface{}, serverSetter *server.ServerSetter) error {
			cfg := config.(*Config)

			// Resolve HTTP handler from DI
			httpHandler := dependencyinjection.Resolve[*handlers.OrderHandler](resolver)

			// Create HTTP mux
			mux := http.NewServeMux()

			// Register routes
			mux.HandleFunc("POST /orders", httpHandler.CreateOrder)
			mux.HandleFunc("GET /orders/{id}", httpHandler.GetOrder)
			mux.HandleFunc("GET /orders", httpHandler.GetAllOrders)
			mux.HandleFunc("GET /stats", httpHandler.GetStatistics)
			mux.HandleFunc("GET /health", httpHandler.HealthCheck)

			// Configure server setter
			serverSetter.Name = "orchestrator"
			serverSetter.Addr = cfg.Port
			serverSetter.Handler = mux
			serverSetter.ServerType = server.HTTPServer

			logger.Info("===========================================")
			logger.Info("   Saga Order Fulfillment Orchestrator")
			logger.Info("===========================================")
			logger.Info(fmt.Sprintf("[Orchestrator] Starting on %s", cfg.Port))
			logger.Info("[Orchestrator] Microservices:")
			logger.Info(fmt.Sprintf("[Orchestrator]   Payment:   %s", cfg.Services.PaymentURL))
			logger.Info(fmt.Sprintf("[Orchestrator]   Inventory: %s", cfg.Services.InventoryURL))
			logger.Info(fmt.Sprintf("[Orchestrator]   Shipping:  %s", cfg.Services.ShippingURL))
			logger.Info("[Orchestrator] Endpoints:")
			logger.Info("[Orchestrator]   POST   /orders       - Create order (start saga)")
			logger.Info("[Orchestrator]   GET    /orders/:id   - Get order by ID")
			logger.Info("[Orchestrator]   GET    /orders       - Get all orders")
			logger.Info("[Orchestrator]   GET    /stats        - Get saga statistics")
			logger.Info("[Orchestrator]   GET    /health       - Health check")
			logger.Info("[Orchestrator] Features:")
			logger.Info("[Orchestrator]   ✓ Redux state management")
			logger.Info("[Orchestrator]   ✓ Saga pattern with compensation")
			logger.Info("[Orchestrator]   ✓ Config hot-reload (edit orchestrator.json)")
			logger.Info("[Orchestrator]   ✓ 30% random failure rate per service")
			logger.Info("===========================================")

			return nil
		}).
		GetListener()

	if err != nil {
		panic(err)
	}

	// Start listener (blocks, handles SIGINT/SIGTERM, config hot-reload)
	<-listener.Start()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
