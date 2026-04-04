package employee

import (
	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-redux/v2"
	"github.com/janmbaco/go-redux/v2/actions"
	"github.com/janmbaco/go-redux/v2/handlers"
)

// EmployeeState represents the Redux state
type EmployeeState struct {
	Employees map[string]*Employee
}

// Actions
var (
	CreateEmployeeAction     = actions.NewAction[CreateEmployeePayload]("CREATE_EMPLOYEE")
	UpdateEmployeeAction     = actions.NewAction[UpdateEmployeePayload]("UPDATE_EMPLOYEE")
	CompleteOnboardingAction = actions.NewAction[CompleteOnboardingPayload]("COMPLETE_ONBOARDING")
	FailOnboardingAction     = actions.NewAction[FailOnboardingPayload]("FAIL_ONBOARDING")
	RollbackOnboardingAction = actions.NewAction[RollbackOnboardingPayload]("ROLLBACK_ONBOARDING")
)

// CreateEmployeePayload for creating employee
type CreateEmployeePayload struct {
	Employee *Employee
}

// UpdateEmployeePayload for updating employee status
type UpdateEmployeePayload struct {
	EmployeeID string
	Status     Status
	Step       string
}

// CompleteOnboardingPayload for completing onboarding
type CompleteOnboardingPayload struct {
	EmployeeID string
}

// FailOnboardingPayload for failing onboarding
type FailOnboardingPayload struct {
	EmployeeID string
	Step       string
	Error      string
}

// RollbackOnboardingPayload for rolling back onboarding
type RollbackOnboardingPayload struct {
	EmployeeID string
}

// Reducers
func createEmployeeReducer(state EmployeeState, payload CreateEmployeePayload) EmployeeState {
	state.Employees[payload.Employee.ID] = payload.Employee
	return state
}

func updateEmployeeReducer(state EmployeeState, payload UpdateEmployeePayload) EmployeeState {
	if emp, exists := state.Employees[payload.EmployeeID]; exists {
		emp.UpdateStatus(payload.Status, payload.Step)
	}
	return state
}

func completeOnboardingReducer(state EmployeeState, payload CompleteOnboardingPayload) EmployeeState {
	if emp, exists := state.Employees[payload.EmployeeID]; exists {
		emp.MarkCompleted()
	}
	return state
}

func failOnboardingReducer(state EmployeeState, payload FailOnboardingPayload) EmployeeState {
	if emp, exists := state.Employees[payload.EmployeeID]; exists {
		emp.MarkFailed(payload.Step)
	}
	return state
}

func rollbackOnboardingReducer(state EmployeeState, payload RollbackOnboardingPayload) EmployeeState {
	if emp, exists := state.Employees[payload.EmployeeID]; exists {
		emp.MarkRolledBack()
	}
	return state
}

// NewEmployeeActionHandler creates action handler
func NewEmployeeActionHandler() handlers.ActionHandler[EmployeeState] {
	builder := handlers.NewActionHandlerBuilder[EmployeeState]()

	builder.On(CreateEmployeeAction, func(state EmployeeState, payload CreateEmployeePayload) EmployeeState {
		return createEmployeeReducer(state, payload)
	})

	builder.On(UpdateEmployeeAction, func(state EmployeeState, payload UpdateEmployeePayload) EmployeeState {
		return updateEmployeeReducer(state, payload)
	})

	builder.On(CompleteOnboardingAction, func(state EmployeeState, payload CompleteOnboardingPayload) EmployeeState {
		return completeOnboardingReducer(state, payload)
	})

	builder.On(FailOnboardingAction, func(state EmployeeState, payload FailOnboardingPayload) EmployeeState {
		return failOnboardingReducer(state, payload)
	})

	builder.On(RollbackOnboardingAction, func(state EmployeeState, payload RollbackOnboardingPayload) EmployeeState {
		return rollbackOnboardingReducer(state, payload)
	})

	return builder.Build()
}

// CreateEmployeeStore creates the Redux store
func CreateEmployeeStore(logger logs.Logger) redux.Store[EmployeeState] {
	initialState := EmployeeState{Employees: make(map[string]*Employee)}
	store := redux.NewStore(initialState, logger)
	store.AddModule(NewEmployeeActionHandler())
	return store
}
