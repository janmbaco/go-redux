# Go Redux

Go Redux is an implementation of [Redux principles](https://redux.js.org/understanding/thinking-in-redux/three-principles) for golang application development.

## Table of Contents

- [Motivation](#motivation)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Example](#example)

## Motivation
On the ecosystem of open code tools for GO I did not find an implementation of the redux pattern that complied with Redux and SOLID principles.

To create clean code, I needed a tool that would clearly separate business logic from what is purely infrastructure where "a single source of truth" was built by injecting the application's business logic.

To do this, I considered that the store should be built with *business param* that would contain the *pure function to perform state transitions* (*Reducer*) as well as some utilities like  initial state value, the *Selector* and action definitions.

Any state transition must be done by launching an action through the *Reducer* that would execute the business logic. And this would be managed from the *single source of truth* (*Store*).

To prevent state mutation, the *store state manager* always provides a copy of the current state. Only when an *action* is launched and the *reducer* executes the *business logic* to get the next state is when the state transitions to the next one, and only in the case that the next state value is different from the previous one.

## Installation


```bash
$ go get github.com/janmbaco/redux
```

## Quick-start

### Actions

An *Action* describes *What* is going to be done and *With* what information.
Defining actions is as simple as declaring the class, adding attributes of type *redux.Action*. Then you can create the object with that declaration.

```go
...
type CounterActions struct {
	Increment redux.Action
	Decrement redux.Action
}
...
func main() {
	counterActions := &CounterActions{}
...
```

### BusinessParam

A *BusinessParam* is an object that contains the *ActionsObject*, the *Reducer*, the *InitialState* and the *Selector*.

An *ActionsObject* is an object that contains the actions that can be performed with a certain *business logic*.

A *Reducer* is a [pure function](https://en.wikipedia.org/wiki/Pure_function) with `(state interface{}, action redux.Action) => state interface{}` signature.
It describes *How* an action transforms the state into the next state.

An *InitialState* is the first value of the state. The shape of the state is up to you: it can be a primitive, an array or an object.

A *Selector* is a string that identifies a part of the global state array that BusinessParam Reducer can transform.

To create the *BusinessParam* you can use a builder (*BusinessParamBuilder*), to which the *InitialState* and *Actions* must be injected.
```go
...
	builder:= redux.NewBusinessParamBuilder(0, actions)
...
```
From here, there are two different ways of injecting business logic:

 1. By associating each action to a function that returns the next state:
```go
...
func Increment(state int, payload int) int {
	return state + payload
}
...
	builder.On(actions.Increment, Increment)
...
```
 2. By setting an object that contains functions with the same name as actions, each function returning the next state
```go
...
type DecrementLogic struct{
}

func (r *DecrementLogic) Decrement(state int, payload int) int{
	return state - payload
}
...
	builder.SetActionsLogicByObject(&SubstractionLogic{})
...
```
In both cases functions must have at least as input parameter the previous state, as well as optionally a payload parameter from the action.

Optionally, a selector can be configured to identify a part of the state array.
```go
...
	builder.SetSelector("counter")
...
```
When all actions have been associated with a business logic, the *BusinessParam* can be obtained.
```go
...
	counterParam := builder.GetBusinessParam()
...
```

### Store

A *Store* is the object that contains the *Single Source of Truth*, that is, it holds the global state of the application, furthermore, it is the only place where a state can transform into the next state.
It is built by injecting the *BusinessParams*
```go
...
	store := redux.NewStore(counterParam)
...
```
To get the global state, just call the *GetState* function
```go
...
	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: 'map[counter:0]'
...
```
To execute an action, you only have to call the *Dispatch* function with it, you can add the necessary payload with the *With* function of the action.
```go
...
	store.Dispatch(actions.Increment.With(1))
	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: 'map[counter:1]'
...
```
To subscribe to state transitions you can use *Subscribe(fn \*func())* function, a pointer to a function is necessary, so a variable must be defined that points to the subscription function and then uses the variable pointer as parameter.
```go
...
	globalSubscription := func(){
		fmt.Printf("globalSubscription - state changed, current state: '%v'\n", store.GetState())
	}
	store.Subscribe(&globalSubscription)
	store.Dispatch(actions.Decrement.With(1))
	// output:
	//globalSubscription - state changed, current state: 'map[counter:0]'
...
```
To unsubscribe from status changes you can use *UnSubscribe(fn \*func())* function, it is necessary to use the same variable pointer as parameter.
```go
...
	store.UnSubscribe(&globalSubscription)
...
```
To add another *Reducer* to the *Store* you can use *AddReducer* function
```go
...
    counter2Actions := &CounterActions{}
    counter2Param := redux.NewBusinessParamBuilder(10, counter2Actions).
        On(counter2Actions.Increment, Increment).
        SetActionsLogicByObject(&DecrementLogic{}).
        SetSelector("counter2").
        GetBusinessParam()

    store.AddReducer(counter2Param)

    counter3Actions := &CounterActions{}
    counter3Param := redux.NewBusinessParamBuilder(100, counter3Actions).
        On(counter3Actions.Increment, Increment).
        SetActionsLogicByObject(&DecrementLogic{}).
        SetSelector("counter3").
        GetBusinessParam()

    store.AddReducer(counter3Param)


	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: 'map[counter:0 counter2:10 counter3:100]'

	store.Dispatch(actions.Increment.With(1))
	store.Dispatch(actions2.Increment.With(1))
	store.Dispatch(actions3.Increment.With(1))

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: 'map[counter:1 counter2:11 counter3:101]'
...
```
To remove a *Reducer* to the *Store* you can use *RemoveReducer* function
```go
...
	store.RemoveReducer("counter3")
	func() {
		defer func() {
			if re := recover(); re != nil {
				fmt.Printf(re.(string) + "\n")
			}
		}()
		store.Dispatch(counter3Actions.Increment.With(1))
	}()
	// output:
	//There are not any Reducers that execute this action!

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: 'map[counter:1 counter2:11 counter3:101]'
...
```
When you remove a *Reducer* from the *Store*, only the business logic and the actions related to it are deleted, the part of the state array that transitioned will remain in its last state, since deleting a reducer should not change the global state.

To get a part of the state array you can use the function *GetStateOf(selector)*, you need know the selector defined in the *BusinessParam*.

```go
...
	fmt.Printf("current state counter: '%v'\n", store.GetStateOf("counter"))
	fmt.Printf("current state counter2: '%v'\n", store.GetStateOf("counter2"))
	fmt.Printf("current state counter3: '%v'\n", store.GetStateOf("counter3"))
	// output:
	//current state counter: '1'
	//current state counter2: '11'
	//current state counter3: '101'
...
```
To subscribe to the state changes of a part of the global state array, you can use *SubscribeTo(selector, \*func(interface{}))*. As with the *Subscribe* function, you need a variable that points to a function.
```go
...
	counterSubscribe := func(newState interface{}) {
		fmt.Printf("counterSubscribe - state changed, current state: '%v'\n", newState)
	}
	store.SubscribeTo("counter", &counterSubscribe)
	store.Dispatch(counterActions.Decrement.With(1))
	store.Dispatch(counter2Actions.Decrement.With(11))
	// output:
	//counterSubscribe - state changed, current state: '0'

	counter2Subscribe := func(newState interface{}) {
		fmt.Printf("counter2Subscribe - state changed, current state: '%v'\n", newState)
	}
	store.SubscribeTo("counter2", &counter2Subscribe)
	store.Dispatch(counterActions.Increment.With(1))
	store.Dispatch(counter2Actions.Increment.With(10))
	// output:
	//ounterSubscribe - state changed, current state: '1'
	//counter2Subscribe - state changed, current state: '10'
...
```
To unsubscribe to the state changes of a part of the global state array, you can use *UnSubscribeFrom(selector, \*func(interface{}))*. As with the *UnSubscribe* function, you need the same variable that is used in *SubscribeTo* function.
```go
...
	store.UnSubscribeFrom("counter", &counterSubscribe)
	store.Dispatch(counterActions.Increment.With(5))
	store.Dispatch(counter2Actions.Decrement.With(8))
	// output:
	//counter2Subscribe - state changed, current state: '2'

	fmt.Printf("current state counter: '%v'\n", store.GetStateOf("counter"))
	// output:
	//current state counter: '6'
...
```

## Example

```go
package main

import (
	"fmt"
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

	counterActions := &CounterActions{}

	builder := redux.NewBusinessParamBuilder(0, counterActions)
	builder.On(counterActions.Increment, Increment)
	builder.SetActionsLogicByObject(&DecrementLogic{})
	builder.SetSelector("counter")

	counterParam := builder.GetBusinessParam()

	store := redux.NewStore(counterParam)

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: 'map[counter:0]'

	store.Dispatch(counterActions.Increment.With(1))
	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: 'map[counter:1]'

	globalSubscription := func() {
		fmt.Printf("globalSubscription - state changed, current state: '%v'\n", store.GetState())
	}
	store.Subscribe(&globalSubscription)
	store.Dispatch(counterActions.Decrement.With(1))
	// output:
	//	//globalSubscription - state changed, current state: 'map[counter:0]'

	store.UnSubscribe(&globalSubscription)

	counter2Actions := &CounterActions{}
	counter2Param := redux.NewBusinessParamBuilder(10, counter2Actions).
		On(counter2Actions.Increment, Increment).
		SetActionsLogicByObject(&DecrementLogic{}).
		SetSelector("counter2").
		GetBusinessParam()

	store.AddReducer(counter2Param)

	counter3Actions := &CounterActions{}
	counter3Param := redux.NewBusinessParamBuilder(100, counter3Actions).
		On(counter3Actions.Increment, Increment).
		SetActionsLogicByObject(&DecrementLogic{}).
		SetSelector("counter3").
		GetBusinessParam()

	store.AddReducer(counter3Param)

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: 'map[counter:0 counter2:10 counter3:100]'

	store.Dispatch(counterActions.Increment.With(1))
	store.Dispatch(counter2Actions.Increment.With(1))
	store.Dispatch(counter3Actions.Increment.With(1))

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: 'map[counter:1 counter2:11 counter3:101]'

	store.RemoveReducer("counter3")
	func() {
		defer func() {
			if re := recover(); re != nil {
				fmt.Printf(re.(string) + "\n")
			}
		}()
		store.Dispatch(counter3Actions.Increment.With(1))
	}()
	// output:
	//There are not any Reducers that execute this action!

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: 'map[counter:1 counter2:11 counter3:101]'

	fmt.Printf("current state counter: '%v'\n", store.GetStateOf("counter"))
	fmt.Printf("current state counter2: '%v'\n", store.GetStateOf("counter2"))
	fmt.Printf("current state counter3: '%v'\n", store.GetStateOf("counter3"))
	// output:
	//current state counter: '1'
	//current state counter2: '11'
	//current state counter3: '101'

	counterSubscribe := func(newState interface{}) {
		fmt.Printf("counterSubscribe - state changed, current state: '%v'\n", newState)
	}
	store.SubscribeTo("counter", &counterSubscribe)
	store.Dispatch(counterActions.Decrement.With(1))
	store.Dispatch(counter2Actions.Decrement.With(11))
	// output:
	//counterSubscribe - state changed, current state: '0'

	counter2Subscribe := func(newState interface{}) {
		fmt.Printf("counter2Subscribe - state changed, current state: '%v'\n", newState)
	}
	store.SubscribeTo("counter2", &counter2Subscribe)
	store.Dispatch(counterActions.Increment.With(1))
	store.Dispatch(counter2Actions.Increment.With(10))
	// output:
	//ounterSubscribe - state changed, current state: '1'
	//counter2Subscribe - state changed, current state: '10'

	store.UnSubscribeFrom("counter", &counterSubscribe)
	store.Dispatch(counterActions.Increment.With(5))
	store.Dispatch(counter2Actions.Decrement.With(8))
	// output:
	//counter2Subscribe - state changed, current state: '2'

	fmt.Printf("current state counter: '%v'\n", store.GetStateOf("counter"))
	// output:
	//current state counter: '6'

}

```

