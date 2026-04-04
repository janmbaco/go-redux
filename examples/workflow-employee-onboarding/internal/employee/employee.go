package employee

import (
	"time"

	"github.com/google/uuid"
)

// Employee represents an employee in the onboarding process
type Employee struct {
	ID          string
	FirstName   string
	LastName    string
	Email       string
	Department  string
	Position    string
	HireDate    time.Time
	Status      Status
	CurrentStep string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Status represents the onboarding status
type Status string

const (
	StatusPending    Status = "PENDING"
	StatusInProgress Status = "IN_PROGRESS"
	StatusCompleted  Status = "COMPLETED"
	StatusFailed     Status = "FAILED"
	StatusRolledBack Status = "ROLLED_BACK"
)

// NewEmployee creates a new employee
func NewEmployee(firstName, lastName, email, department, position string, hireDate time.Time) *Employee {
	return &Employee{
		ID:          uuid.New().String(),
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		Department:  department,
		Position:    position,
		HireDate:    hireDate,
		Status:      StatusPending,
		CurrentStep: "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// UpdateStatus updates the employee status and current step
func (e *Employee) UpdateStatus(status Status, step string) {
	e.Status = status
	e.CurrentStep = step
	e.UpdatedAt = time.Now()
}

// MarkCompleted marks the onboarding as completed
func (e *Employee) MarkCompleted() {
	e.Status = StatusCompleted
	e.CurrentStep = "COMPLETED"
	e.UpdatedAt = time.Now()
}

// MarkFailed marks the onboarding as failed
func (e *Employee) MarkFailed(step string) {
	e.Status = StatusFailed
	e.CurrentStep = step
	e.UpdatedAt = time.Now()
}

// MarkRolledBack marks the onboarding as rolled back
func (e *Employee) MarkRolledBack() {
	e.Status = StatusRolledBack
	e.CurrentStep = "ROLLED_BACK"
	e.UpdatedAt = time.Now()
}
