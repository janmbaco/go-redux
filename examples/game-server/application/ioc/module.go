package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-redux/v2/examples/game-server/application"
	"github.com/janmbaco/go-redux/v2/handlers"
)

// GameApplicationModule registers application-layer services
type GameApplicationModule struct{}

// NewGameApplicationModule creates a new application module
func NewGameApplicationModule() *GameApplicationModule {
	return &GameApplicationModule{}
}

// RegisterServices registers handlers with their dependencies
func (m *GameApplicationModule) RegisterServices(register dependencyinjection.Register) error {
	// Register PlayerHandler as singleton tenant with dependency injection
	dependencyinjection.RegisterSingletonTenantWithParams[handlers.ActionHandler[application.GameState]](
		register,
		"PlayerHandler",
		application.NewPlayerHandler,
		nil, // auto-resolve MovementService
	)

	// Register CombatHandler as singleton tenant with dependency injection
	dependencyinjection.RegisterSingletonTenantWithParams[handlers.ActionHandler[application.GameState]](
		register,
		"CombatHandler",
		application.NewCombatHandler,
		nil, // auto-resolve CombatService
	)

	return nil
}
