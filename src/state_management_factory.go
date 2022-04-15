package redux

import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

type StateManagementFactoryParamter struct{
	InitialState interface{}
	Selector string
	StorePublisher eventsmanager.Publisher
}

type StateManagementFactory interface {
	Create(parameter StateManagementFactoryParamter) StateManagement
}

type stateManagementFactory struct {
	resolver dependencyinjection.Resolver
}

func NewStateManagementFactory(container dependencyinjection.Container) StateManagementFactory {
	container.Register().AsType(new(StateManagement), NewStateManager, map[uint]string{0: _initialState, 1: _selector, 2: _storePublisher})
	return &stateManagementFactory{container.Resolver()}
}

func (s *stateManagementFactory) Create(parameter StateManagementFactoryParamter) StateManagement {
	return s.resolver.Type(new(StateManagement), map[string]interface{}{_initialState: parameter.InitialState, _selector: parameter.Selector, _storePublisher: parameter.StorePublisher}).(*stateManagement)
}
