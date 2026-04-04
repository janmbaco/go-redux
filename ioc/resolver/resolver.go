package resolver

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	redux "github.com/janmbaco/go-redux/v2"
	"github.com/janmbaco/go-redux/v2/handlers"
)

// GetStore resolves a Store from the IoC container with initialState
func GetStore[S any](resolver dependencyinjection.Resolver, initialState S) redux.Store[S] {
	return dependencyinjection.ResolveWithParams[redux.Store[S]](
		resolver,
		map[string]any{"initialState": initialState},
	)
}

// GetHandler resolves an ActionHandler from the IoC container
func GetHandler[S any](resolver dependencyinjection.Resolver) handlers.ActionHandler[S] {
	return dependencyinjection.Resolve[handlers.ActionHandler[S]](resolver)
}
