package ioc

import (
	di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/eventstore"
)

type EventSourceModule struct{}

func NewEventSourceModule() *EventSourceModule {
	return &EventSourceModule{}
}

func (m *EventSourceModule) RegisterServices(register di.Register) error {
	// Register EventStore as singleton
	di.RegisterSingletonTenant(
		register,
		"memory",
		eventstore.NewMemoryEventStore,
	)
	di.RegisterSingletonTenantWithParams[eventstore.EventStore](
		register,
		"sqlite",
		eventstore.NewSQLiteEventStore,
		nil,
	)
	return nil
}
