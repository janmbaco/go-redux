package ioc

import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	"github.com/janmbaco/go-redux/src"
)

func init(){
	static.Container.Register().AsSingleton(new(redux.ActionsObjectFactory), redux.NewActionsObjectFactory, nil)
	static.Container.Register().AsSingleton(new(redux.BusinessParamFactory), redux.NewBusinessParamFactory, nil)
	static.Container.Register().AsSingleton(new(redux.BusinesParamBuilder), redux.NewBusinessParamBuilder, nil)
	static.Container.Register().AsSingleton(new(redux.StateManagementFactory), redux.NewStateManagementFactory, nil)
	static.Container.Register().AsSingleton(new(redux.Store), redux.NewStore, nil)
}

