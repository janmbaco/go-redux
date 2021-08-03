package main

import (
	"fmt"
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/logs"
	"github.com/janmbaco/go-redux/src"
)

type CounterActions struct {
	Increment redux.Action
	Decrement redux.Action
}

func Increment(state int, payload int) int {
	return state + payload
}

type DecrementLogic struct {
}

func (r *DecrementLogic) Decrement(state int, payload int) int {
	return state - payload
}

func main() {

	container := dependencyinjection.NewContainer()
	facade(container.Register())

	counterActions := &CounterActions{}

	store := container.Resolver().Type(new(redux.Store), nil).(redux.Store)

	builder := container.Resolver().Type(new(redux.BusinesParamBuilder), nil).(redux.BusinesParamBuilder)
	builder.SetInitialState(0)
	builder.SetActions(counterActions)
	builder.On(counterActions.Increment, Increment)
	builder.SetActionsLogicByObject(&DecrementLogic{})
	builder.SetSelector("counter")

	counterParam := builder.GetBusinessParam()

	store.AddReducer(counterParam)

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: 'map[counter:0]'

	store.Dispatch(counterActions.Increment.With(1))
	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: 'map[counter:1]'

	globalSubscription := func() {
		fmt.Printf("globalSubscription - state changed, current state: '%v'\n", store.GetState())
	}
	store.Subscribe(&globalSubscription)
	store.Dispatch(counterActions.Decrement.With(1))
	// output:
	// globalSubscription - state changed, current state: 'map[counter:0]'

	store.Unsubscribe(&globalSubscription)

	counter2Actions := &CounterActions{}
	counter2Param := builder.
		SetInitialState(10).
		SetActions(counter2Actions).
		On(counter2Actions.Increment, Increment).
		SetActionsLogicByObject(&DecrementLogic{}).
		SetSelector("counter2").
		GetBusinessParam()

	store.AddReducer(counter2Param)

	counter3Actions := &CounterActions{}
	counter3Param := builder.
		SetInitialState(100).
		SetActions(counter3Actions).
		On(counter3Actions.Increment, Increment).
		SetActionsLogicByObject(&DecrementLogic{}).
		SetSelector("counter3").
		GetBusinessParam()

	store.AddReducer(counter3Param)

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: 'map[counter:0 counter2:10 counter3:100]'

	store.Dispatch(counterActions.Increment.With(1))
	store.Dispatch(counter2Actions.Increment.With(1))
	store.Dispatch(counter3Actions.Increment.With(1))

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: 'map[counter:1 counter2:11 counter3:101]'

	store.RemoveReducer("counter3")
	errorManager := container.Resolver().Type(new(errors.ErrorManager), nil).(errors.ErrorManager)
	errorManager.On(&redux.StoreError{}, func(err error) {
		// deactivation of the store's own errors so that only an error message appears
		// see redux.StoreError
		fmt.Printf(err.Error() + "\n")
	})
	store.Dispatch(counter3Actions.Increment.With(1))
	// output:
	// There are not any Reducers that execute this action!

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: 'map[counter:1 counter2:11 counter3:101]'

	fmt.Printf("current state counter: '%v'\n", store.GetStateOf("counter"))
	fmt.Printf("current state counter2: '%v'\n", store.GetStateOf("counter2"))
	fmt.Printf("current state counter3: '%v'\n", store.GetStateOf("counter3"))
	// output:
	// current state counter: '1'
	// current state counter2: '11'
	// current state counter3: '101'

	counterSubscribe := func(newState interface{}) {
		fmt.Printf("counterSubscribe - state changed, current state: '%v'\n", newState)
	}
	store.SubscribeTo("counter", &counterSubscribe)
	store.Dispatch(counterActions.Decrement.With(1))
	store.Dispatch(counter2Actions.Decrement.With(11))
	// output:
	// counterSubscribe - state changed, current state: '0'

	counter2Subscribe := func(newState interface{}) {
		fmt.Printf("counter2Subscribe - state changed, current state: '%v'\n", newState)
	}
	store.SubscribeTo("counter2", &counter2Subscribe)
	store.Dispatch(counterActions.Increment.With(1))
	store.Dispatch(counter2Actions.Increment.With(10))
	// output:
	// ounterSubscribe - state changed, current state: '1'
	// counter2Subscribe - state changed, current state: '10'

	store.UnsubscribeFrom("counter", &counterSubscribe)
	store.Dispatch(counterActions.Increment.With(5))
	store.Dispatch(counter2Actions.Decrement.With(8))
	// output:
	// counter2Subscribe - state changed, current state: '2'

	fmt.Printf("current state counter: '%v'\n", store.GetStateOf("counter"))
	// output:
	// current state counter: '6'
}

func facade(register dependencyinjection.Register) {
	register.AsSingleton(new(logs.Logger), logs.NewLogger, nil)
	register.Bind(new(logs.ErrorLogger), new(logs.Logger))
	register.AsSingleton(new(errors.ErrorCatcher), errors.NewErrorCatcher, nil)
	register.AsSingleton(new(errors.ErrorManager), errors.NewErrorManager, nil)
	register.Bind(new(errors.ErrorCallbacks), new(errors.ErrorManager))
	register.AsSingleton(new(errors.ErrorThrower), errors.NewErrorThrower, nil)
	register.AsType(new(redux.ActionsObject), redux.NewActionsObject, map[uint]string{1: "actions"})
	register.AsType(new(redux.BusinessParam), redux.NewBusinessParam, map[uint]string{0: "initialState", 1: "reducer", 2: "actionsObject", 3: "selector"})
	register.AsSingleton(new(redux.BusinesParamBuilder), redux.NewBusinessParamBuilder, nil)
	register.AsSingleton(new(redux.Store), redux.NewStore, nil)
}
