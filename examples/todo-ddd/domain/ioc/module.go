package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-redux/v2/examples/todo-ddd/domain"
)

// TodoDomainModule registers domain-layer services
type TodoDomainModule struct{}

// NewTodoDomainModule creates a new domain module
func NewTodoDomainModule() *TodoDomainModule {
	return &TodoDomainModule{}
}

// RegisterServices registers domain services
func (m *TodoDomainModule) RegisterServices(register dependencyinjection.Register) error {
	// Register domain service
	dependencyinjection.RegisterSingleton(register, domain.NewTodoService)

	return nil
}
