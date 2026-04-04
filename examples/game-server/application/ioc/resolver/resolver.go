package resolver

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-redux/v2/examples/game-server/application"
	"github.com/janmbaco/go-redux/v2/handlers"
)

// GetPlayerHandler resolves the PlayerHandler from the IoC container
func GetPlayerHandler(resolver dependencyinjection.Resolver) handlers.ActionHandler[application.GameState] {
	return dependencyinjection.ResolveTenant[handlers.ActionHandler[application.GameState]](resolver, "PlayerHandler")
}

// GetCombatHandler resolves the CombatHandler from the IoC container
func GetCombatHandler(resolver dependencyinjection.Resolver) handlers.ActionHandler[application.GameState] {
	return dependencyinjection.ResolveTenant[handlers.ActionHandler[application.GameState]](resolver, "CombatHandler")
}
