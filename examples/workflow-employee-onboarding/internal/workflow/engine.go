package workflow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-redux/v2"
	"github.com/janmbaco/go-redux/v2/examples/workflow-employee-onboarding/internal/employee"
	"github.com/janmbaco/go-redux/v2/examples/workflow-employee-onboarding/internal/tasks"
)

// WorkflowRepository stores workflow state
type WorkflowRepository interface {
	SaveWorkflow(employeeID string, wf *WorkflowDefinition) error
	LoadWorkflow(employeeID string) (*WorkflowDefinition, error)
	DeleteWorkflow(employeeID string) error
	ListWorkflows() ([]string, error)
}

// WorkflowEngine orchestrates the onboarding workflow
type WorkflowEngine struct {
	store      redux.Store[employee.EmployeeState]
	repository WorkflowRepository
	tasks      *tasks.TaskExecutor
	logger     logs.Logger
	mu         sync.Mutex
}

// NewWorkflowEngine creates a new workflow engine
func NewWorkflowEngine(
	store redux.Store[employee.EmployeeState],
	repository WorkflowRepository,
	taskExecutor *tasks.TaskExecutor,
	logger logs.Logger,
) *WorkflowEngine {
	return &WorkflowEngine{
		store:      store,
		repository: repository,
		tasks:      taskExecutor,
		logger:     logger,
	}
}

// StartOnboarding starts the onboarding workflow for an employee
func (we *WorkflowEngine) StartOnboarding(ctx context.Context, emp *employee.Employee) error {
	we.mu.Lock()
	defer we.mu.Unlock()

	we.logger.Info(fmt.Sprintf("Starting onboarding workflow for employee %s (%s %s)",
		emp.ID, emp.FirstName, emp.LastName))

	// Create employee in store
	we.store.Dispatch(employee.CreateEmployeeAction.With(employee.CreateEmployeePayload{Employee: emp}))

	// Create workflow definition
	workflow := NewOnboardingWorkflow()

	// Save workflow state
	if err := we.repository.SaveWorkflow(emp.ID, workflow); err != nil {
		return fmt.Errorf("failed to save workflow: %w", err)
	}

	// Update employee status
	we.store.Dispatch(employee.UpdateEmployeeAction.With(employee.UpdateEmployeePayload{
		EmployeeID: emp.ID,
		Status:     employee.StatusInProgress,
		Step:       string(StepHRSetup),
	}))

	// Execute workflow asynchronously
	go we.executeWorkflow(ctx, emp.ID, workflow)

	return nil
}

// executeWorkflow executes the workflow steps
func (we *WorkflowEngine) executeWorkflow(ctx context.Context, employeeID string, workflow *WorkflowDefinition) {
	we.logger.Info(fmt.Sprintf("Executing workflow for employee %s", employeeID))

	for !workflow.IsCompleted() && !workflow.HasFailed() {
		select {
		case <-ctx.Done():
			we.logger.Info(fmt.Sprintf("Workflow execution cancelled for employee %s", employeeID))
			return
		default:
		}

		// Get executable steps
		executableSteps := workflow.GetExecutableSteps()
		if len(executableSteps) == 0 {
			// No steps can be executed, wait a bit
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Execute steps in parallel
		var wg sync.WaitGroup
		for _, stepName := range executableSteps {
			wg.Add(1)
			go func(step StepName) {
				defer wg.Done()
				we.executeStep(ctx, employeeID, workflow, step)
			}(stepName)
		}
		wg.Wait()

		// Save workflow state after each iteration
		if err := we.repository.SaveWorkflow(employeeID, workflow); err != nil {
			we.logger.Error(fmt.Sprintf("Failed to save workflow state: %v", err))
		}
	}

	// Check final status
	if workflow.IsCompleted() {
		we.logger.Info(fmt.Sprintf("Workflow completed successfully for employee %s", employeeID))
		we.store.Dispatch(employee.CompleteOnboardingAction.With(employee.CompleteOnboardingPayload{EmployeeID: employeeID}))
		we.repository.SaveWorkflow(employeeID, workflow)
	} else if workflow.HasFailed() {
		failedStep, err := workflow.GetFailedStep()
		we.logger.Error(fmt.Sprintf("Workflow failed at step %s for employee %s: %v", failedStep, employeeID, err))

		// Trigger rollback
		we.rollbackWorkflow(ctx, employeeID, workflow)
	}
}

// executeStep executes a single workflow step
func (we *WorkflowEngine) executeStep(ctx context.Context, employeeID string, workflow *WorkflowDefinition, stepName StepName) {
	we.logger.Info(fmt.Sprintf("Executing step %s for employee %s", stepName, employeeID))

	workflow.MarkRunning(stepName)
	we.store.Dispatch(employee.UpdateEmployeeAction.With(employee.UpdateEmployeePayload{
		EmployeeID: employeeID,
		Status:     employee.StatusInProgress,
		Step:       string(stepName),
	}))

	// Execute the task
	var err error
	switch stepName {
	case StepHRSetup:
		err = we.tasks.HRSetup(employeeID)
	case StepITProvisioning:
		err = we.tasks.ITProvisioning(employeeID)
	case StepTraining:
		err = we.tasks.Training(employeeID)
	case StepBadge:
		err = we.tasks.Badge(employeeID)
	}

	if err != nil {
		we.logger.Error(fmt.Sprintf("Step %s failed for employee %s: %v", stepName, employeeID, err))
		workflow.MarkFailed(stepName, err)
		we.store.Dispatch(employee.FailOnboardingAction.With(employee.FailOnboardingPayload{
			EmployeeID: employeeID,
			Step:       string(stepName),
			Error:      err.Error(),
		}))
	} else {
		we.logger.Info(fmt.Sprintf("Step %s completed for employee %s", stepName, employeeID))
		workflow.MarkCompleted(stepName)
	}
}

// rollbackWorkflow rolls back completed steps
func (we *WorkflowEngine) rollbackWorkflow(ctx context.Context, employeeID string, workflow *WorkflowDefinition) {
	we.logger.Info(fmt.Sprintf("Rolling back workflow for employee %s", employeeID))

	completedSteps := workflow.GetCompletedSteps()
	for _, stepName := range completedSteps {
		select {
		case <-ctx.Done():
			we.logger.Info(fmt.Sprintf("Rollback cancelled for employee %s", employeeID))
			return
		default:
		}

		we.logger.Info(fmt.Sprintf("Rolling back step %s for employee %s", stepName, employeeID))

		step := workflow.Steps[stepName]
		if err := step.Rollback(employeeID); err != nil {
			we.logger.Error(fmt.Sprintf("Failed to rollback step %s: %v", stepName, err))
		} else {
			workflow.MarkRolledBack(stepName)
		}
	}

	we.store.Dispatch(employee.RollbackOnboardingAction.With(employee.RollbackOnboardingPayload{EmployeeID: employeeID}))
	we.repository.SaveWorkflow(employeeID, workflow)
	we.logger.Info(fmt.Sprintf("Rollback completed for employee %s", employeeID))
}

// ResumeWorkflow resumes a workflow after a crash
func (we *WorkflowEngine) ResumeWorkflow(ctx context.Context, employeeID string) error {
	we.mu.Lock()
	defer we.mu.Unlock()

	we.logger.Info(fmt.Sprintf("Resuming workflow for employee %s", employeeID))

	// Load workflow state
	workflow, err := we.repository.LoadWorkflow(employeeID)
	if err != nil {
		return fmt.Errorf("failed to load workflow: %w", err)
	}

	// Continue execution
	go we.executeWorkflow(ctx, employeeID, workflow)

	return nil
}

// GetWorkflowStatus returns the current workflow status
func (we *WorkflowEngine) GetWorkflowStatus(employeeID string) (*WorkflowDefinition, error) {
	return we.repository.LoadWorkflow(employeeID)
}
