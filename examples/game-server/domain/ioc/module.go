package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-redux/v2/examples/game-server/domain"
)

// GameDomainModule registers domain-layer services
type GameDomainModule struct{}

// NewGameDomainModule creates a new domain module
func NewGameDomainModule() *GameDomainModule {
	return &GameDomainModule{}
}

// RegisterServices registers domain services
func (m *GameDomainModule) RegisterServices(register dependencyinjection.Register) error {
	// Register domain services
	dependencyinjection.RegisterSingleton(register, domain.NewMovementService)
	dependencyinjection.RegisterSingleton(register, domain.NewCombatService)

	return nil
}
