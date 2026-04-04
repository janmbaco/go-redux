package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-redux/v2"
	"github.com/janmbaco/go-redux/v2/examples/workflow-employee-onboarding/internal/employee"
	"github.com/janmbaco/go-redux/v2/examples/workflow-employee-onboarding/internal/workflow"
)

// OnboardingHandler handles HTTP requests
type OnboardingHandler struct {
	store  redux.Store[employee.EmployeeState]
	engine *workflow.WorkflowEngine
	logger logs.Logger
}

// NewOnboardingHandler creates a new handler
func NewOnboardingHandler(
	store redux.Store[employee.EmployeeState],
	engine *workflow.WorkflowEngine,
	logger logs.Logger,
) *OnboardingHandler {
	return &OnboardingHandler{
		store:  store,
		engine: engine,
		logger: logger,
	}
}

// CreateOnboardingRequest represents the create onboarding request
type CreateOnboardingRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Department string `json:"department"`
	Position   string `json:"position"`
	HireDate   string `json:"hire_date"`
}

// CreateOnboarding creates a new employee onboarding
func (h *OnboardingHandler) CreateOnboarding(w http.ResponseWriter, r *http.Request) {
	var req CreateOnboardingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Parse hire date
	hireDate, err := time.Parse("2006-01-02", req.HireDate)
	if err != nil {
		http.Error(w, "invalid hire_date format (use YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// Create employee
	emp := employee.NewEmployee(
		req.FirstName,
		req.LastName,
		req.Email,
		req.Department,
		req.Position,
		hireDate,
	)

	// Start onboarding workflow
	if err := h.engine.StartOnboarding(context.Background(), emp); err != nil {
		h.logger.Error(fmt.Sprintf("Failed to start onboarding: %v", err))
		http.Error(w, "failed to start onboarding", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"employee_id": emp.ID,
		"message":     "onboarding started",
	})
}

// GetEmployee returns employee details
func (h *OnboardingHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	employeeID := r.URL.Query().Get("id")
	if employeeID == "" {
		http.Error(w, "employee_id required", http.StatusBadRequest)
		return
	}

	state := h.store.GetState()
	emp, exists := state.Employees[employeeID]
	if !exists {
		http.Error(w, "employee not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

// GetWorkflowStatus returns workflow execution status
func (h *OnboardingHandler) GetWorkflowStatus(w http.ResponseWriter, r *http.Request) {
	employeeID := r.URL.Query().Get("id")
	if employeeID == "" {
		http.Error(w, "employee_id required", http.StatusBadRequest)
		return
	}

	wf, err := h.engine.GetWorkflowStatus(employeeID)
	if err != nil {
		http.Error(w, "workflow not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wf)
}

// ResumeWorkflow resumes a workflow after crash
func (h *OnboardingHandler) ResumeWorkflow(w http.ResponseWriter, r *http.Request) {
	employeeID := r.URL.Query().Get("id")
	if employeeID == "" {
		http.Error(w, "employee_id required", http.StatusBadRequest)
		return
	}

	if err := h.engine.ResumeWorkflow(context.Background(), employeeID); err != nil {
		h.logger.Error(fmt.Sprintf("Failed to resume workflow: %v", err))
		http.Error(w, "failed to resume workflow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "workflow resumed",
	})
}

// ListEmployees returns all employees
func (h *OnboardingHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	state := h.store.GetState()

	employees := make([]*employee.Employee, 0, len(state.Employees))
	for _, emp := range state.Employees {
		employees = append(employees, emp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":     len(employees),
		"employees": employees,
	})
}

// HealthCheck returns service health
func (h *OnboardingHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"service": "workflow-employee-onboarding",
		"status":  "healthy",
	})
}
