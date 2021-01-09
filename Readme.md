# Go Redux

Go Redux is an implementation of [Redux principles](https://redux.js.org/understanding/thinking-in-redux/three-principles) for golang application development.

## Table of Contents

- [Motivation](#motivation)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Example](#example)

## Motivation
On the ecosstem of open code tools for GO I did not find an implementation of the redux pattern that complied with Redux principles and SOLID principles.

To create clean code, I needed a tool that clearly separated business logic from what is purely infrastructure. Where "a single source of truth" is built with the business logic of the application.

To do this, I considered that the store should be made from business objects. These objects would contain the pure function to make changes to the state (reducer) as well as some utilities like the initial state value, the selector and the actions definitions.

Any state change  should be made by dispatching the action through the reducer of the business object that would execute the business logic.

To avoid state mutation, the store state manager always sends the business object a copy of the current state, and changes the state when the copy it receives from the reducer is different from the original.

## Installation


```bash
$ go get github.com/janmbaco/redux
```

## Quick-start

### Actions

*Actions* are objects that could contain information payloads that are sent for a state change in the *store*.
In this case, the actions are interfaces, which are added to a class.

```go
...
type Actions struct {
	Sum redux.Action
	Substraction redux.Action
}
...
func main() {
	actions := &Actions{}
...
```

### BusinessObject

A *BusinessObject* is an object that contains the *ActionsObject*, the *Reducer*, the initial state and a *Selector*.

A "ActionsObject" is an object that contains the actions that can be performed with a certain business logic.

A *Reducer* is a [pure function](https://en.wikipedia.org/wiki/Pure_function) with `(state interface{}, action redux.Action) => state interface{}` signature.
It describes how an action transform  the state into the next state.

The initial state is the first value of the state. The shape of the state is up to you: it can be a primitive, an array or an object.

A *Selector* is a string that identifies a part of the global state array that BusinessObject Reducer can transform.

To create a business object you can use a builder (BusinessObjectBuilder), to which the initial state and actions must be injected.
```go
...
	builder:= redux.NewBusinessObjectBuilder(0, actions)
...
```
 From there, the different business logics associated with the actions can be established in two ways:

 1. By associating each action to a function, that gets the state and payload and returns the next state:
```go
...
func Sum(state int, payload int) int {
	return state + payload
}
...
	builder.On(actions.Sum, Sum)
...
```
 2. By setting an object that contains functions with the same name as actions and each function gets the state and payload and returns the next state
```go
...
type SubstractionLogic struct{
}

func (r *SubstractionLogic ) Substraction(state int, payload int) int{
	return state - payload
}
...
	builder.SetActionsLogicByObject(&SubstractionLogic{})
...
```

You can set a selector to identify the part of the global state array, and thus use it in the store.
```go
...
	builder.SetSelector("counter")
...
```
After the logic is set, you can get the *BusinessObject*
```go
...
	counterBO := builder.GetBusinessObject()
...
```

### Store

A *Store* is the object that contains the *Single Source of Truth*, that is, it holds the global state of the application, furthermore, it is the only place where a state can transform into the next state.
It is buit by inyecting the *BusinessObjects*
```go
...
	store := redux.NewStore(counterBO)
...
```
To get the global state, you just have to call the GetState() function
```go
...
	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: '0'
...
```
To run an action, you just need to call Dispatch function with the action, you can add the required payload with the With function of the action.
```go
...
	store.Dispatch(actions.Sum.With(1))
	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: '1'
...
```
To subscribe to state changes, you need a pointer to a function, so you will have to assign the function to subscribe  to a pointer and pass the pointer to the *Subscribe(fn *func())*
```go
...
	globalSubscription := func(){
		fmt.Printf("globalSubscription - state changed, current state: '%v'\n", store.GetState())
	}
	store.Subscribe(&globalSubscription)
	store.Dispatch(actions.Substraction.With(1))
	// output:
	// globalSubscription - state changed, current state: '0'
...
```
To unsubscribe to state changes, you need a pointer to a function that is subscribed and pass the pointer to the *UnSubscribe(fn *func())*
```go
...
	store.UnSubscribe(&globalSubscription)
...
```
To  add more *BusinessObject* to the *Store* you can use AddBusinessObject function
```go
...
	actions2 := &Actions{}
	counter2BO := redux.NewBusinessObjectBuilder(10, actions2).
		On(actions2.Sum, Sum).
		SetActionsLogicByObject(&SubstractionLogic{}).
		SetSelector("counter2").
		GetBusinessObject()

	store.AddBusinessObject(counter2BO)

	actions3 := &Actions{}
	counter3BO := redux.NewBusinessObjectBuilder(100, actions3).
		On(actions3.Sum, Sum).
		SetActionsLogicByObject(&SubstractionLogic{}).
		SetSelector("counter3").
		GetBusinessObject()

	store.AddBusinessObject(counter3BO)

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: '[0 10 100]'

	store.Dispatch(actions.Sum.With(1))
	store.Dispatch(actions2.Sum.With(1))
	store.Dispatch(actions3.Sum.With(1))

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: '[1 11 101]'
...
```
To  remove a *BusinessObject* to the *Store* you can use RemoveBusinessObject function
```go
...
	store.RemoveBusinessObject(counter3BO)
	func() {
	   defer func(){
	      if re := recover(); re != nil {
	         fmt.Printf(re.(string)+"\n")
	      }
	   }()
	   store.Dispatch(actions3.Sum.With(1))
	}()
	// output:
	//There is not any Reducers that execute this action!

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: '[1 11 101]'
...
```
When you remove a *BusinessObject* from store, only the logic and actions related to it are deleted, the state array part remains in its last state, because only a reducer can change it.

To get a part of the state array, you can use the selectors with the function *GetStateOf(selector)*.

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
You can subscribe to the state changes of a part of the global state array with the *SubscribeTo(selector, \*Subscribefn)*, you need a pointer to a function, so you will have to assign the function to subscribe  to a pointer and pass the pointer to the *SubscribeTo(selector, \*func(interface{}))*
```go
...
	counterSubscribe := func(newState interface{}){
	   fmt.Printf("counterSubscribe - state changed, current state: '%v'\n", newState)
	}
	store.SubscribeTo("counter", &counterSubscribe)
	store.Dispatch(actions.Substraction.With(1))
	store.Dispatch(actions2.Substraction.With(11))
	// output:
	//counterSubscribe - state changed, current state: '0'

	counter2Subscribe := func(newState interface{}){
	   fmt.Printf("counter2Subscribe - state changed, current state: '%v'\n", newState)
	}
	store.SubscribeTo("counter2", &counter2Subscribe)
	store.Dispatch(actions.Sum.With(1))
	store.Dispatch(actions2.Sum.With(10))
	// output:
	//ounterSubscribe - state changed, current state: '1'
	//counter2Subscribe - state changed, current state: '10'
...
```
To unsubscribe to the state changes of a part of the global state array, you need a pointer to a function that is subscribed and pass the pointer to the *UnSubscribeFrom(selector, \*func(interface{}))*
```go
...
	store.UnSubscribeFrom("counter", &counterSubscribe)
	store.Dispatch(actions.Sum.With(5))
	store.Dispatch(actions2.Substraction.With(8))
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
	"github.com/janmbaco/go-redux/src")

type Actions struct {
	Sum redux.Action
	Substraction redux.Action
}

func Sum(state int, payload int) int {
	return state + payload
}

type SubstractionLogic struct{
}

func (r *SubstractionLogic ) Substraction(state int, payload int) int{
	return state - payload
}

func main() {

	actions := &Actions{}

	builder:= redux.NewBusinessObjectBuilder(0, actions)
	builder.On(actions.Sum, Sum)
	builder.SetActionsLogicByObject(&SubstractionLogic{})
	builder.SetSelector("counter")

	counterBO := builder.GetBusinessObject()

	store := redux.NewStore(counterBO)

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: '0'

	store.Dispatch(actions.Sum.With(1))
	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: '1'

	globalSubscription := func(){
	fmt.Printf("globalSubscription - state changed, current state: '%v'\n", store.GetState())
	}
	store.Subscribe(&globalSubscription)
	store.Dispatch(actions.Substraction.With(1))
	// output:
	// globalSubscription - state changed, current state: '0'

	store.UnSubscribe(&globalSubscription)

	actions2 := &Actions{}
	counter2BO := redux.NewBusinessObjectBuilder(10, actions2).
		On(actions2.Sum, Sum).
		SetActionsLogicByObject(&SubstractionLogic{}).
		SetSelector("counter2").
		GetBusinessObject()

	store.AddBusinessObject(counter2BO)

	actions3 := &Actions{}
	counter3BO := redux.NewBusinessObjectBuilder(100, actions3).
		On(actions3.Sum, Sum).
		SetActionsLogicByObject(&SubstractionLogic{}).
		SetSelector("counter3").
		GetBusinessObject()

	store.AddBusinessObject(counter3BO)

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: '0'

	store.Dispatch(actions.Sum.With(1))
	store.Dispatch(actions2.Sum.With(1))
	store.Dispatch(actions3.Sum.With(1))

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: '[1 11 101]'

	store.RemoveBusinessObject(counter3BO)
	func() {
		defer func(){
		   if re := recover(); re != nil {
		      fmt.Printf(re.(string)+"\n")
		   }
		}()
		store.Dispatch(actions3.Sum.With(1))
	}()
	// output:
	//There is not any Reducers that execute this action!
	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: '[1 11 101]'
	fmt.Printf("current state counter: '%v'\n", store.GetStateOf("counter"))
	fmt.Printf("current state counter2: '%v'\n", store.GetStateOf("counter2"))
	fmt.Printf("current state counter3: '%v'\n", store.GetStateOf("counter3"))
	// output:
	//current state counter: '1' //current state counter2: '11' //current state counter3: '101'

	counterSubscribe := func(newState interface{}){
		fmt.Printf("counterSubscribe - state changed, current state: '%v'\n", newState)
	}
	store.SubscribeTo("counter", &counterSubscribe)
	store.Dispatch(actions.Substraction.With(1))
	store.Dispatch(actions2.Substraction.With(11))
	// output:
	//counterSubscribe - state changed, current state: '0'

	counter2Subscribe := func(newState interface{}){
		fmt.Printf("counter2Subscribe - state changed, current state: '%v'\n", newState)
	}
	store.SubscribeTo("counter2", &counter2Subscribe)
	store.Dispatch(actions.Sum.With(1))
	store.Dispatch(actions2.Sum.With(10))
	// output:
	//ounterSubscribe - state changed, current state: '1' //counter2Subscribe - state changed, current state: '10'

	store.UnSubscribeFrom("counter", &counterSubscribe)
	store.Dispatch(actions.Sum.With(5))
	store.Dispatch(actions2.Substraction.With(8))
	// output:
	//counter2Subscribe - state changed, current state: '2'

	fmt.Printf("current state counter: '%v'\n", store.GetStateOf("counter"))
	// output:
	//current state counter: '6'
}
```

