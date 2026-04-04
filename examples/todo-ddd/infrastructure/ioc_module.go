package infrastructure

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-redux/v2/examples/todo-ddd/application"
	"github.com/janmbaco/go-redux/v2/examples/todo-ddd/domain"
	"github.com/janmbaco/go-redux/v2/handlers"
)

// TodoModule registers todo-related services in the IoC container
type TodoModule struct{}

// NewTodoModule creates a new todo IoC module
func NewTodoModule() *TodoModule {
	return &TodoModule{}
}

// RegisterServices registers all todo services
func (m *TodoModule) RegisterServices(register dependencyinjection.Register) error {
	// Register domain service
	dependencyinjection.RegisterSingleton(register, domain.NewTodoService)

	// Register todo handler with dependency injection
	var todoHandler handlers.ActionHandler[application.AppState]
	register.AsSingleton(&todoHandler, func(todoService *domain.TodoService) handlers.ActionHandler[application.AppState] {
		return application.NewTodoHandler(todoService)
	}, nil)

	// Register filter handler (no dependencies)
	var filterHandler handlers.ActionHandler[application.AppState]
	register.AsSingleton(&filterHandler, application.NewFilterHandler, nil)

	return nil
}
