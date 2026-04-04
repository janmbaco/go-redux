package persistence

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/janmbaco/go-redux/v2/examples/workflow-employee-onboarding/internal/workflow"
)

// MemoryWorkflowRepository is an in-memory implementation of workflow.WorkflowRepository
type MemoryWorkflowRepository struct {
	workflows map[string]*workflow.WorkflowDefinition
	mu        sync.RWMutex
}

// NewMemoryWorkflowRepository creates a new in-memory repository
func NewMemoryWorkflowRepository() *MemoryWorkflowRepository {
	return &MemoryWorkflowRepository{
		workflows: make(map[string]*workflow.WorkflowDefinition),
	}
}

// SaveWorkflow saves workflow state
func (r *MemoryWorkflowRepository) SaveWorkflow(employeeID string, wf *workflow.WorkflowDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Deep copy to avoid mutation issues
	data, err := json.Marshal(wf)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow: %w", err)
	}

	var copy workflow.WorkflowDefinition
	if err := json.Unmarshal(data, &copy); err != nil {
		return fmt.Errorf("failed to unmarshal workflow: %w", err)
	}

	r.workflows[employeeID] = &copy
	return nil
}

// LoadWorkflow loads workflow state
func (r *MemoryWorkflowRepository) LoadWorkflow(employeeID string) (*workflow.WorkflowDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	wf, exists := r.workflows[employeeID]
	if !exists {
		return nil, fmt.Errorf("workflow not found for employee %s", employeeID)
	}

	// Deep copy
	data, err := json.Marshal(wf)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow: %w", err)
	}

	var copy workflow.WorkflowDefinition
	if err := json.Unmarshal(data, &copy); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow: %w", err)
	}

	return &copy, nil
}

// DeleteWorkflow deletes workflow state
func (r *MemoryWorkflowRepository) DeleteWorkflow(employeeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.workflows, employeeID)
	return nil
}

// ListWorkflows returns all employee IDs with workflows
func (r *MemoryWorkflowRepository) ListWorkflows() ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.workflows))
	for id := range r.workflows {
		ids = append(ids, id)
	}
	return ids, nil
}
