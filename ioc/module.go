package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	redux "github.com/janmbaco/go-redux/v2"
)

// ReduxModule implements Module for Redux services
type ReduxModule[S any] struct{}

// NewReduxModule creates a new Redux IoC module.
// The initial state is supplied when resolving the store through named parameters.
func NewReduxModule[S any]() *ReduxModule[S] {
	return &ReduxModule[S]{}
}

// RegisterServices registers all Redux services
func (m *ReduxModule[S]) RegisterServices(register dependencyinjection.Register) error {
	// Register Store as singleton with initialState as named parameter
	dependencyinjection.RegisterSingletonWithParams[redux.Store[S]](
		register,
		redux.NewStore[S],
		map[int]string{0: "initialState"},
	)

	return nil
}
