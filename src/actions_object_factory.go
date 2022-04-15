package redux

import "github.com/janmbaco/go-infrastructure/dependencyinjection"

type ActionsObjectFactory interface {
	Create(actions interface{}) ActionsObject
}

type actionsObjectFactory struct {
	resolver dependencyinjection.Resolver
}

func NewActionsObjectFactory(container dependencyinjection.Container) ActionsObjectFactory {
	container.Register().AsType(new(ActionsObject), NewActionsObject, map[uint]string{1: _actions})
	return &actionsObjectFactory{container.Resolver()}
}

func (a *actionsObjectFactory) Create(actions interface{}) ActionsObject {
	return a.resolver.Type(new(ActionsObject), map[string]interface{}{_actions: actions}).(ActionsObject)
}