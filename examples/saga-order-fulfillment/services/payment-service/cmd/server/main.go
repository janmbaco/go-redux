package main

import (
	"fmt"
	"net/http"
	"os"

	configResolver "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc/resolver"
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-infrastructure/v2/persistence"
	persistenceioc "github.com/janmbaco/go-infrastructure/v2/persistence/ioc"
	"github.com/janmbaco/go-infrastructure/v2/server"
	serverIoc "github.com/janmbaco/go-infrastructure/v2/server/ioc"
	serverResolver "github.com/janmbaco/go-infrastructure/v2/server/ioc/resolver"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/services/payment-service/internal/handlers"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/services/payment-service/internal/ioc"
)

// Config represents payment service configuration
type Config struct {
	Port     string `json:"port"`
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"database"`
}

func main() {
	// Default configuration
	defaultConfig := &Config{}
	defaultConfig.Port = ":8081"
	defaultConfig.Database.Host = getEnv("DB_HOST", "payment-db")
	defaultConfig.Database.Port = getEnv("DB_PORT", "5432")
	defaultConfig.Database.User = getEnv("DB_USER", "payment")
	defaultConfig.Database.Password = getEnv("DB_PASSWORD", "")
	defaultConfig.Database.Name = getEnv("DB_NAME", "paymentdb")

	// Build container with all required modules
	container := dependencyinjection.NewBuilder().
		AddModules(serverIoc.ConfigureServerModules()...).
		AddModule(persistenceioc.ConfigureDatabaseModule(
			defaultConfig.Database.Host,
			defaultConfig.Database.Port,
			defaultConfig.Database.User,
			defaultConfig.Database.Password,
			defaultConfig.Database.Name,
			persistence.Postgres,
		)).
		AddModule(ioc.NewPaymentModule()).
		MustBuild()

	resolver := container.Resolver()

	// Get config handler with file monitoring
	configHandler := configResolver.GetFileConfigHandler(
		resolver,
		"payment-service.json",
		defaultConfig,
	)

	// Get listener builder with bootstrapper
	listener, err := serverResolver.GetListenerBuilder(resolver, configHandler).
		SetBootstrapper(func(config interface{}, serverSetter *server.ServerSetter) error {
			cfg := config.(*Config)

			// Resolve handler from DI
			handler := dependencyinjection.Resolve[*handlers.PaymentHandler](resolver)

			// Get logger
			logger := dependencyinjection.Resolve[logs.Logger](resolver)

			// Create HTTP mux
			mux := http.NewServeMux()

			// Register routes
			mux.HandleFunc("POST /payments", handler.ProcessPayment)
			mux.HandleFunc("POST /payments/{id}/refund", handler.RefundPayment)
			mux.HandleFunc("GET /payments/{id}", handler.GetPayment)
			mux.HandleFunc("GET /payments/order/{orderId}", handler.GetPaymentByOrderID)
			mux.HandleFunc("GET /health", handler.HealthCheck)

			// Configure server setter
			serverSetter.Name = "payment-service"
			serverSetter.Addr = cfg.Port
			serverSetter.Handler = mux
			serverSetter.ServerType = server.HTTPServer

			logger.Info(fmt.Sprintf("[Payment] Starting Payment Service on %s", cfg.Port))
			logger.Info("[Payment] Endpoints:")
			logger.Info("[Payment]   POST   /payments              - Process payment")
			logger.Info("[Payment]   POST   /payments/:id/refund   - Refund payment (compensation)")
			logger.Info("[Payment]   GET    /payments/:id          - Get payment by ID")
			logger.Info("[Payment]   GET    /payments/order/:orderId - Get payments by order")
			logger.Info("[Payment]   GET    /health                - Health check")
			logger.Info("[Payment] Note: 30% random failure rate, config hot-reload enabled")

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
