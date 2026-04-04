package workflow

import (
	"fmt"
	"time"
)

// StepName represents workflow step names
type StepName string

const (
	StepHRSetup        StepName = "HR_SETUP"
	StepITProvisioning StepName = "IT_PROVISIONING"
	StepTraining       StepName = "TRAINING"
	StepBadge          StepName = "BADGE"
)

// StepStatus represents step execution status
type StepStatus string

const (
	StepStatusPending    StepStatus = "PENDING"
	StepStatusRunning    StepStatus = "RUNNING"
	StepStatusCompleted  StepStatus = "COMPLETED"
	StepStatusFailed     StepStatus = "FAILED"
	StepStatusRolledBack StepStatus = "ROLLED_BACK"
)

// Step represents a workflow step
type Step struct {
	Name         StepName                      `json:"name"`
	Status       StepStatus                    `json:"status"`
	StartedAt    *time.Time                    `json:"started_at,omitempty"`
	CompletedAt  *time.Time                    `json:"completed_at,omitempty"`
	Error        string                        `json:"error,omitempty"`
	Dependencies []StepName                    `json:"dependencies"`
	Rollback     func(employeeID string) error `json:"-"` // Don't serialize function
}

// WorkflowDefinition defines the onboarding workflow
type WorkflowDefinition struct {
	Steps map[StepName]*Step
}

// NewOnboardingWorkflow creates the employee onboarding workflow definition
func NewOnboardingWorkflow() *WorkflowDefinition {
	return &WorkflowDefinition{
		Steps: map[StepName]*Step{
			StepHRSetup: {
				Name:         StepHRSetup,
				Status:       StepStatusPending,
				Dependencies: []StepName{},
				Rollback: func(employeeID string) error {
					// Rollback HR setup (delete employee record, revoke HR access)
					time.Sleep(100 * time.Millisecond) // Simulate work
					return nil
				},
			},
			StepITProvisioning: {
				Name:         StepITProvisioning,
				Status:       StepStatusPending,
				Dependencies: []StepName{StepHRSetup}, // Requires HR setup
				Rollback: func(employeeID string) error {
					// Rollback IT (delete accounts, revoke access)
					time.Sleep(100 * time.Millisecond)
					return nil
				},
			},
			StepTraining: {
				Name:         StepTraining,
				Status:       StepStatusPending,
				Dependencies: []StepName{StepHRSetup}, // Can run in parallel with IT
				Rollback: func(employeeID string) error {
					// Rollback training (cancel sessions, remove from LMS)
					time.Sleep(100 * time.Millisecond)
					return nil
				},
			},
			StepBadge: {
				Name:         StepBadge,
				Status:       StepStatusPending,
				Dependencies: []StepName{StepITProvisioning, StepTraining}, // Requires both IT and Training
				Rollback: func(employeeID string) error {
					// Rollback badge (deactivate, return to security)
					time.Sleep(100 * time.Millisecond)
					return nil
				},
			},
		},
	}
}

// CanExecute checks if a step can be executed based on dependencies
func (wd *WorkflowDefinition) CanExecute(stepName StepName) bool {
	step := wd.Steps[stepName]
	if step.Status != StepStatusPending {
		return false
	}

	// Check if all dependencies are completed
	for _, dep := range step.Dependencies {
		depStep := wd.Steps[dep]
		if depStep.Status != StepStatusCompleted {
			return false
		}
	}

	return true
}

// MarkRunning marks a step as running
func (wd *WorkflowDefinition) MarkRunning(stepName StepName) {
	step := wd.Steps[stepName]
	step.Status = StepStatusRunning
	now := time.Now()
	step.StartedAt = &now
}

// MarkCompleted marks a step as completed
func (wd *WorkflowDefinition) MarkCompleted(stepName StepName) {
	step := wd.Steps[stepName]
	step.Status = StepStatusCompleted
	now := time.Now()
	step.CompletedAt = &now
}

// MarkFailed marks a step as failed
func (wd *WorkflowDefinition) MarkFailed(stepName StepName, err error) {
	step := wd.Steps[stepName]
	step.Status = StepStatusFailed
	step.Error = err.Error()
	now := time.Now()
	step.CompletedAt = &now
}

// MarkRolledBack marks a step as rolled back
func (wd *WorkflowDefinition) MarkRolledBack(stepName StepName) {
	step := wd.Steps[stepName]
	step.Status = StepStatusRolledBack
}

// GetExecutableSteps returns steps that can be executed now
func (wd *WorkflowDefinition) GetExecutableSteps() []StepName {
	var executable []StepName
	for name := range wd.Steps {
		if wd.CanExecute(name) {
			executable = append(executable, name)
		}
	}
	return executable
}

// GetCompletedSteps returns completed steps in reverse order (for rollback)
func (wd *WorkflowDefinition) GetCompletedSteps() []StepName {
	var completed []StepName

	// Order: Badge -> Training/IT -> HR (reverse of execution)
	stepOrder := []StepName{StepBadge, StepTraining, StepITProvisioning, StepHRSetup}

	for _, name := range stepOrder {
		if wd.Steps[name].Status == StepStatusCompleted {
			completed = append(completed, name)
		}
	}

	return completed
}

// IsCompleted checks if all steps are completed
func (wd *WorkflowDefinition) IsCompleted() bool {
	for _, step := range wd.Steps {
		if step.Status != StepStatusCompleted {
			return false
		}
	}
	return true
}

// HasFailed checks if any step has failed
func (wd *WorkflowDefinition) HasFailed() bool {
	for _, step := range wd.Steps {
		if step.Status == StepStatusFailed {
			return true
		}
	}
	return false
}

// GetFailedStep returns the failed step if any
func (wd *WorkflowDefinition) GetFailedStep() (StepName, error) {
	for name, step := range wd.Steps {
		if step.Status == StepStatusFailed {
			return name, fmt.Errorf("step %s failed: %s", name, step.Error)
		}
	}
	return "", nil
}
