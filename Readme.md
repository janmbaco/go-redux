# Go Redux
**Redux principles for Go**

## Table of Contents

- [Introduction](#introduction)
- [Why Use This Project?](#why-use-this-project)
- [Motivation](#motivation)
- [Installation](#installation)
- [Quick Start](#quick-start)
  - [Actions](#actions)
  - [BusinessParam](#businessparam)
  - [Store](#store)
- [Example](#example)
- [Contributing](#contributing)
- [License](#license)

## Introduction

Go Redux is an implementation of [Redux principles](https://redux.js.org/understanding/thinking-in-redux/three-principles) for Go application development. It helps in creating clean, maintainable, and scalable code by separating business logic from infrastructure.

## Why Use This Project?

- **Single Source of Truth**: Manage the entire state of your application in one place.
- **Predictable State Changes**: State transitions are handled by pure functions, making them predictable and easy to debug.
- **Maintainable Code**: Clear separation of concerns promotes maintainability and scalability.

## Motivation

In the open-source ecosystem for Go, there was a lack of Redux pattern implementations adhering to Redux and SOLID principles. Go Redux fills this gap by providing a tool that separates business logic from infrastructure, maintaining a single source of truth.

To achieve this, the store is built with a *business param* containing a *pure function* to perform state transitions (*Reducer*) and utilities like the initial state and selectors. State transitions are triggered by actions through the *Reducer* and managed by the *Store*.

## Installation

### Prerequisites

- Go 1.16+
- Git

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/janmbaco/go-redux.git
   ```
2. Install the package:
   ```bash
   go get github.com/janmbaco/redux
   ```

## Quick Start

### Actions

An *Action* describes *what* is going to be done and *with* what information. Define actions by declaring a struct with attributes of type *redux.Action*.

```go
type CounterActions struct {
    Increment redux.Action
    Decrement redux.Action
}

func main() {
    counterActions := &CounterActions{}
}
```

### BusinessParam

A *BusinessParam* is an object that contains the *ActionsObject*, the *Reducer*, the *InitialState*, and the *Selector*. It is used to define the business logic and state transitions for a specific part of your application.

#### Creating BusinessParam

```go
builder := resolver.GetBusinessParamBuilder()
builder.SetInitialState(0)
builder.SetActions(counterActions)
```

#### Injecting Business Logic

1. **By Associating Each Action with a Function:**

```go
func Increment(state int, payload int) int {
    return state + payload
}

builder.On(counterActions.Increment, Increment)
```

2. **By Setting an Object with Functions Corresponding to Actions:**

```go
type DecrementLogic struct{}

func (r *DecrementLogic) Decrement(state int, payload int) int {
    return state - payload
}

builder.SetActionsLogicByObject(&DecrementLogic{})
```

Set the selector to identify a part of the state array:

```go
builder.SetSelector("counter")
counterParam := builder.GetBusinessParam()
```

### Store

A *Store* contains the global state of the application and is the only place where state transformations occur. Add *BusinessParams* like a reducer.

```go
store := resolver.GetStore()
store.AddReducer(counterParam)
```

Get the global state:

```go
fmt.Printf("current state: '%v'\n", store.GetState())
// output:
// current state: 'map[counter:0]'
```

Dispatch an action:

```go
store.Dispatch(counterActions.Increment.With(1))
fmt.Printf("current state: '%v'\n", store.GetState())
// output:
// current state: 'map[counter:1]'
```

Subscribe to state transitions:

```go
globalSubscription := func() {
    fmt.Printf("globalSubscription - state changed, current state: '%v'\n", store.GetState())
}
store.Subscribe(&globalSubscription)
store.Dispatch(counterActions.Decrement.With(1))
// output:
// globalSubscription - state changed, current state: 'map[counter:0]'
```

Unsubscribe from state transitions:

```go
store.UnSubscribe(&globalSubscription)
```

Add another *Reducer* to the *Store*:

```go
counter2Actions := &CounterActions{}
counter2Param := builder.
    SetInitialState(10).
    SetActions(counter2Actions).
    On(counter2Actions.Increment, Increment).
    SetActionsLogicByObject(&DecrementLogic{}).
    SetSelector("counter2").
    GetBusinessParam()

store.AddReducer(counter2Param)
```

Remove a *Reducer* from the *Store*:

```go
store.RemoveReducer("counter3")
```

Get a part of the state array:

```go
fmt.Printf("current state counter: '%v'\n", store.GetStateOf("counter"))
```

Subscribe to the state changes of a part of the global state array:

```go
counterSubscribe := func(newState interface{}) {
    fmt.Printf("counterSubscribe - state changed, current state: '%v'\n", newState)
}
store.SubscribeTo("counter", &counterSubscribe)
store.Dispatch(counterActions.Decrement.With(1))
```

Unsubscribe from the state changes of a part of the global state array:

```go
store.UnSubscribeFrom("counter", &counterSubscribe)
```

## Example

```go
package main

import (
    "fmt"
    "github.com/janmbaco/go-redux/src"
    "github.com/janmbaco/go-redux/src/ioc/resolver"
    errorsResolver "github.com/janmbaco/go-infrastructure/errors/ioc/resolver"
)

type CounterActions struct {
    Increment redux.Action
    Decrement redux.Action
}

func Increment(state int, payload int) int {
    return state + payload
}

type DecrementLogic struct{}

func (r *DecrementLogic) Decrement(state int, payload int) int {
    return state - payload
}

func main() {
    counterActions := &CounterActions{}
    store := resolver.GetStore()
    builder := resolver.GetBusinessParamBuilder()

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
    errorManager := errorsResolver.GetErrorManager()
    errorManager.On(new(redux.StoreError), func(err error) {
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
    // counterSubscribe - state changed, current state: '1'
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
```

## Contributing

We welcome contributions! Please feel free to submit pull requests.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.
