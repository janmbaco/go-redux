package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/logs"
	iocLog "github.com/janmbaco/go-infrastructure/v2/logs/ioc"
	"github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/eventstore"
	"github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/handlers"
	"github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/projections"
	"github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/snapshots"
	"github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/wallet"
)

func main() {
	fmt.Println("🏦 Event-Sourced Wallet Server Starting...")
	conainer, err := di.NewBuilder().
		AddModule(iocLog.NewLogsModule()).
		Build()

	if err != nil {
		panic("Failed to build IoC container: " + err.Error())
	}
	// Create logger
	logger := di.Resolve[logs.Logger](conainer.Resolver())

	// Create components
	eventStore := eventstore.NewMemoryEventStore()
	projection := projections.NewBalanceProjection()
	snapshotStore := snapshots.NewSnapshotStore()
	snapshotter := snapshots.NewSnapshotter(snapshotStore, 100, logger)

	// Create Redux store
	store := wallet.CreateWalletStore(logger)

	// Create handler
	handler := handlers.NewWalletHandler(store, eventStore, projection, snapshotter, logger)

	// Setup HTTP routes
	mux := http.NewServeMux()

	// Wallet operations
	mux.HandleFunc("POST /wallets", handler.CreateWallet)
	mux.HandleFunc("GET /wallets/{id}", handler.GetWallet)
	mux.HandleFunc("POST /wallets/{id}/deposit", handler.Deposit)
	mux.HandleFunc("POST /wallets/{id}/withdraw", handler.Withdraw)
	mux.HandleFunc("POST /wallets/{id}/transfer", handler.Transfer)

	// CQRS read model
	mux.HandleFunc("GET /wallets/{id}/balance", handler.GetBalance)

	// Event sourcing operations
	mux.HandleFunc("GET /wallets/{id}/events", handler.GetEvents)
	mux.HandleFunc("POST /wallets/{id}/replay", handler.Replay)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"service": "event-sourced-wallet",
			"status":  "healthy",
		})
	})

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	logger.Info(fmt.Sprintf("🚀 Server listening on port %s", port))
	logger.Info("📋 Endpoints:")
	logger.Info("  POST   /wallets                  - Create wallet")
	logger.Info("  GET    /wallets/:id              - Get wallet")
	logger.Info("  POST   /wallets/:id/deposit      - Deposit money")
	logger.Info("  POST   /wallets/:id/withdraw     - Withdraw money")
	logger.Info("  POST   /wallets/:id/transfer     - Transfer money")
	logger.Info("  GET    /wallets/:id/balance      - Get balance (projection)")
	logger.Info("  GET    /wallets/:id/events       - Get event history")
	logger.Info("  POST   /wallets/:id/replay       - Replay events")
	logger.Info("  GET    /health                   - Health check")

	// Create server with graceful shutdown
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed: " + err.Error())
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 Shutting down server...")
	logger.Info("✅ Server stopped")
}
