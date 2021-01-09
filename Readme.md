# Go Redux

Go Redux is an implementation of [Redux principles](https://redux.js.org/understanding/thinking-in-redux/three-principles) for golang application development.

## Table of Contents

- [Motivation](#motivation)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Example](#example)

## Motivation
On the ecosstem of open code tools for GO I did not find an implementation of the redux pattern that complied with Redux and SOLID principles.

To create clean code, I needed a tool that would clearly separate business logic from what is purely infrastructure. Where "a single source of truth" was built by injecting the application's business logic.

To do this, I considered that the store should be built with *business objects*. These objects would contain the *pure function to perform the state transitions* (*Reducer*) as well as some utilities like the initial state value, the *Selector* and the action definitions.

Any state transition must be done by launching an action through the *business object reducer* that would execute the business logic. And this would be managed from the *single source of truth* (*Store*).

To prevent state mutation, the *store state manager* always provides a copy of the current state. Only when an *action* is launched and the *reducer* executes the *business logic* to get the next state is when the state transitions to the next one, and only in the case that the next state value is different from the previous one.

## Installation


```bash
$ go get github.com/janmbaco/redux
```

## Quick-start

### Actions

The *Actions* describes *What* is going to be done and *With* what information.
Defining actions is as simple as declaring the class, adding attributes of type *redux.Action*. Then you can create the object with that declaration.

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

A *ActionsObject* is an object that contains the actions that can be performed with a certain *business logic*.

A *Reducer* is a [pure function](https://en.wikipedia.org/wiki/Pure_function) with `(state interface{}, action redux.Action) => state interface{}` signature.
It describes *How* an action transform the state into the next state.

A *InitialState* is the first value of the state. The shape of the state is up to you: it can be a primitive, an array or an object.

A *Selector* is a string that identifies a part of the global state array that BusinessObject Reducer can transform.

To create a business object you can use a builder (*BusinessObjectBuilder*), to which the *InitialState* and *Actions* must be injected.
```go
...
	builder:= redux.NewBusinessObjectBuilder(0, actions)
...
```
From here, there are two different ways of injecting business logic:

 1. By associating each action to a function that returns the next state:
```go
...
func Sum(state int, payload int) int {
	return state + payload
}
...
	builder.On(actions.Sum, Sum)
...
```
 2. By setting an object that contains functions with the same name as actions and each function returns the next state
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
In both cases the functions must have at least as an input parameter the previous state, as well as optionally a payload parameter from the action.

Optionally, a selector can be configured to identify a part of the state array.
```go
...
	builder.SetSelector("counter")
...
```
When all actions have been associated with a business logic, the *BusinessObject* can be obtained.
```go
...
	counterBO := builder.GetBusinessObject()
...
```

### Store

A *Store* is the object that contains the *Single Source of Truth*, that is, it holds the global state of the application, furthermore, it is the only place where a state can transform into the next state.
It is built by inyecting the *BusinessObjects*
```go
...
	store := redux.NewStore(counterBO)
...
```
To get the global status, just call the *GetState* function
```go
...
	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: '0'
...
```
To execute an action, you only have to call the *Dispatch* function with the action, you can add the necessary payload with the *With* function of the action.
```go
...
	store.Dispatch(actions.Sum.With(1))
	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: '1'
...
```
To subscribe to state transitions you can use *Subscribe(fn \*func())* function, a pointer to a function is necessary, so a variable must be defined that points to the subscription function and then uses the variable pointer as parameter.
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
To unsubscribe from status changes you can use *UnSubscribe(fn \*func())* function, it is necessary to use the same variable pointer as parameter.
```go
...
	store.UnSubscribe(&globalSubscription)
...
```
To add more *BusinessObject* to the *Store* you can use *AddBusinessObject* function
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
To remove a *BusinessObject* to the *Store* you can use *RemoveBusinessObject* function
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
When you remove a *BusinessObject* from the *Store*, only the business logic and the actions related to it are deleted, the part of the state array that transited will remain in its last state, since deleting a reducer should not change the global state.

To get a part of the state array you can use the function *GetStateOf(selector)*, you need know the selector defined in the *BusinessObject*.

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
To unsubscribe to the state changes of a part of the global state array, you can use *UnSubscribeFrom(selector, \*func(interface{}))*. As with the *UnSubscribe* function, you need the same variable that is used in *SubscribeTo* function.
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
	//current state: '[0 10 100]'

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
	//current state counter: '1'
	//current state counter2: '11'
	//current state counter3: '101'

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

