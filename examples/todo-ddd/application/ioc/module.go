package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-redux/v2/examples/todo-ddd/application"
	"github.com/janmbaco/go-redux/v2/handlers"
)

// TodoApplicationModule registers application-layer services
type TodoApplicationModule struct{}

// NewTodoApplicationModule creates a new application module
func NewTodoApplicationModule() *TodoApplicationModule {
	return &TodoApplicationModule{}
}

// RegisterServices registers handlers with their dependencies
func (m *TodoApplicationModule) RegisterServices(register dependencyinjection.Register) error {
	// Register TodoHandler as singleton tenant with dependency injection
	dependencyinjection.RegisterSingletonTenantWithParams[handlers.ActionHandler[application.AppState]](
		register,
		"TodoHandler",
		application.NewTodoHandler,
		nil, // auto-resolve dependencies by type
	)

	// Register FilterHandler as singleton tenant (no dependencies)
	dependencyinjection.RegisterSingletonTenant[handlers.ActionHandler[application.AppState]](
		register,
		"FilterHandler",
		application.NewFilterHandler,
	)

	return nil
}
