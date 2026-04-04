package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-redux/v2/examples/workflow-employee-onboarding/internal/employee"
	"github.com/janmbaco/go-redux/v2/examples/workflow-employee-onboarding/internal/handlers"
	"github.com/janmbaco/go-redux/v2/examples/workflow-employee-onboarding/internal/persistence"
	"github.com/janmbaco/go-redux/v2/examples/workflow-employee-onboarding/internal/tasks"
	"github.com/janmbaco/go-redux/v2/examples/workflow-employee-onboarding/internal/workflow"
)

func main() {
	// Initialize logger
	logger := logs.NewLogger()
	logger.Info("Starting workflow-employee-onboarding service...")

	// Initialize components
	store := employee.CreateEmployeeStore(logger)
	repository := persistence.NewMemoryWorkflowRepository()
	taskExecutor := tasks.NewTaskExecutor(logger, 0.2) // 20% failure rate for demo
	var repo workflow.WorkflowRepository = repository
	engine := workflow.NewWorkflowEngine(store, repo, taskExecutor, logger)

	// Initialize HTTP handler
	handler := handlers.NewOnboardingHandler(store, engine, logger)

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthCheck)
	mux.HandleFunc("/onboarding", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateOnboarding(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if r.URL.Query().Get("id") != "" {
				handler.GetEmployee(w, r)
			} else {
				handler.ListEmployees(w, r)
			}
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/workflow/status", handler.GetWorkflowStatus)
	mux.HandleFunc("/workflow/resume", handler.ResumeWorkflow)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	// Start server
	go func() {
		logger.Info(fmt.Sprintf("Server listening on port %s", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("Server error: %v", err))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	logger.Info("Server stopped")
}
