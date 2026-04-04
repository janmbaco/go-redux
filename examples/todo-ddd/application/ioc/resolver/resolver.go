package resolver

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-redux/v2/examples/todo-ddd/application"
	"github.com/janmbaco/go-redux/v2/handlers"
)

// GetTodoHandler resolves the TodoHandler from the IoC container
func GetTodoHandler(resolver dependencyinjection.Resolver) handlers.ActionHandler[application.AppState] {
	return dependencyinjection.ResolveTenant[handlers.ActionHandler[application.AppState]](resolver, "TodoHandler")
}

// GetFilterHandler resolves the FilterHandler from the IoC container
func GetFilterHandler(resolver dependencyinjection.Resolver) handlers.ActionHandler[application.AppState] {
	return dependencyinjection.ResolveTenant[handlers.ActionHandler[application.AppState]](resolver, "FilterHandler")
}
