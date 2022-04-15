package resolver

import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	_ "github.com/janmbaco/go-infrastructure/errors/ioc"
	_ "github.com/janmbaco/go-infrastructure/eventsmanager/ioc"
	_ "github.com/janmbaco/go-infrastructure/logs/ioc"
	_ "github.com/janmbaco/go-redux/src/ioc"
	"github.com/janmbaco/go-redux/src"
)

func GetStore() redux.Store {
 	return  static.Container.Resolver().Type(new(redux.Store), nil).(redux.Store)
}

func GetBusinessParamBuilder() redux.BusinesParamBuilder{
	return  static.Container.Resolver().Type(new(redux.BusinesParamBuilder), nil).(redux.BusinesParamBuilder)
}