package redux

import "github.com/janmbaco/go-infrastructure/dependencyinjection"

type BusinessParamFactoryParamter struct {
	InitialState  interface{}
	Reducer       Reducer
	ActionsObject ActionsObject
	Selector      string
}

type BusinessParamFactory interface {
	Create(parameter BusinessParamFactoryParamter) BusinessParam
}

type businessParamFactory struct {
	resolver dependencyinjection.Resolver
}

func NewBusinessParamFactory(container dependencyinjection.Container) BusinessParamFactory {
	container.Register().AsType(new(BusinessParam), NewBusinessParam, map[uint]string{0: _initialState, 1: _reducer, 2: _actionsObject, 3: _selector})
	return &businessParamFactory{container.Resolver()};
}

func (b *businessParamFactory) Create(parameter BusinessParamFactoryParamter) BusinessParam{
	return b.resolver.Type(new(BusinessParam), map[string]interface{}{
		_initialState:  parameter.InitialState,
		_reducer:       parameter.Reducer,
		_selector:      parameter.Selector,
		_actionsObject: parameter.ActionsObject,
	}).(BusinessParam)
}