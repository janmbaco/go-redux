package tasks

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/janmbaco/go-infrastructure/v2/logs"
)

// TaskExecutor executes onboarding tasks
type TaskExecutor struct {
	logger      logs.Logger
	failureRate float64 // Probability of random failure (0.0 - 1.0)
}

// NewTaskExecutor creates a new task executor
func NewTaskExecutor(logger logs.Logger, failureRate float64) *TaskExecutor {
	return &TaskExecutor{
		logger:      logger,
		failureRate: failureRate,
	}
}

// HRSetup performs HR setup tasks
func (te *TaskExecutor) HRSetup(employeeID string) error {
	te.logger.Info(fmt.Sprintf("HR Setup: Creating employee record for %s", employeeID))

	// Simulate work
	time.Sleep(500 * time.Millisecond)

	// Random failure simulation
	if te.shouldFail() {
		return fmt.Errorf("HR system unavailable")
	}

	te.logger.Info(fmt.Sprintf("HR Setup: Employee record created, benefits enrolled"))
	return nil
}

// ITProvisioning performs IT provisioning tasks
func (te *TaskExecutor) ITProvisioning(employeeID string) error {
	te.logger.Info(fmt.Sprintf("IT Provisioning: Creating accounts for %s", employeeID))

	// Simulate work
	time.Sleep(800 * time.Millisecond)

	// Random failure simulation
	if te.shouldFail() {
		return fmt.Errorf("Active Directory sync failed")
	}

	te.logger.Info(fmt.Sprintf("IT Provisioning: Email created, VPN configured, laptop assigned"))
	return nil
}

// Training performs training tasks
func (te *TaskExecutor) Training(employeeID string) error {
	te.logger.Info(fmt.Sprintf("Training: Enrolling %s in onboarding courses", employeeID))

	// Simulate work
	time.Sleep(600 * time.Millisecond)

	// Random failure simulation
	if te.shouldFail() {
		return fmt.Errorf("LMS enrollment failed")
	}

	te.logger.Info(fmt.Sprintf("Training: Courses assigned, schedule sent"))
	return nil
}

// Badge performs badge creation tasks
func (te *TaskExecutor) Badge(employeeID string) error {
	te.logger.Info(fmt.Sprintf("Badge: Creating access badge for %s", employeeID))

	// Simulate work
	time.Sleep(400 * time.Millisecond)

	// Random failure simulation
	if te.shouldFail() {
		return fmt.Errorf("badge printer offline")
	}

	te.logger.Info(fmt.Sprintf("Badge: Access badge created, building access granted"))
	return nil
}

// shouldFail determines if the task should fail randomly
func (te *TaskExecutor) shouldFail() bool {
	return rand.Float64() < te.failureRate
}
