# go-redux

`go-redux` is a Go library for coordinating application state through typed actions, composable handlers, and a central store.

The project is aimed at server-side and backend use cases where state transitions need to remain explicit and testable: orchestration, workflows, event-driven processes, and long-lived coordination logic. It is not intended as a general-purpose framework, nor as a replacement for conventional CRUD services.

## Status

- Module path: `github.com/janmbaco/go-redux/v2`
- Go version: `1.24+`
- Integration point for dependency injection and infrastructure concerns: `github.com/janmbaco/go-infrastructure/v2`

## Installation

```bash
go get github.com/janmbaco/go-redux/v2
```

## Design Goals

- Typed actions with explicit payloads
- Handler composition without central switch statements
- Clear state transitions with a single dispatch entry point
- Practical integration with Go services and DI containers
- Testable coordination logic without introducing a framework runtime

## When It Fits

`go-redux` is a good fit when the value of the model is in the sequence of state transitions rather than in simple data persistence.

Typical cases:

- workflow orchestration
- saga-style coordination
- event-sourced or replayable processes
- game or session state coordination
- in-memory domain coordination with explicit projections or subscribers

It is usually not the right tool for:

- simple CRUD applications
- passive data access layers
- frontend/UI state management

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/janmbaco/go-infrastructure/v2/logs"
	redux "github.com/janmbaco/go-redux/v2"
)

type CounterState struct {
	Count int
}

var IncrementAction = redux.NewAction[int]("counter/increment")
var ResetAction = redux.NewAction[struct{}]("counter/reset")

func main() {
	logger := logs.NewLogger()
	store := redux.NewStore(CounterState{Count: 0}, logger)
	defer store.Close()

	handler := redux.NewActionHandlerBuilder[CounterState]().
		On(IncrementAction, func(state CounterState, amount int) CounterState {
			state.Count += amount
			return state
		}).
		On(ResetAction, func(state CounterState) CounterState {
			state.Count = 0
			return state
		}).
		Build()

	if err := store.AddModule(handler); err != nil {
		panic(err)
	}

	printState := func(state CounterState) {
		fmt.Printf("count=%d\n", state.Count)
	}

	if err := store.Subscribe(&printState); err != nil {
		panic(err)
	}

	if err := store.Dispatch(IncrementAction.With(5)); err != nil {
		panic(err)
	}

	if err := store.Dispatch(ResetAction.With(struct{}{})); err != nil {
		panic(err)
	}
}
```

## Core Concepts

### Actions

Actions define the intent of a state transition and carry a typed payload.

```go
var CreateOrderAction = redux.NewAction[CreateOrderCommand]("order/create")
var CancelOrderAction = redux.NewAction[string]("order/cancel")
```

Dispatch always happens through an action instance:

```go
_ = store.Dispatch(CreateOrderAction.With(CreateOrderCommand{CustomerID: "c-1"}))
_ = store.Dispatch(CancelOrderAction.With("order-123"))
```

### Action Handlers

Handlers bind actions to state transitions.

```go
handler := redux.NewActionHandlerBuilder[OrderState]().
	On(CreateOrderAction, func(state OrderState, cmd CreateOrderCommand) OrderState {
		state.Orders = append(state.Orders, Order{CustomerID: cmd.CustomerID})
		return state
	}).
	Build()
```

The builder supports reducers with:

- `func(S) S`
- `func(S) (S, error)`
- `func(S, P) S`
- `func(S, P) (S, error)`

### Store

The store owns the current state, receives actions through `Dispatch`, and notifies subscribers when the state changes.

The root package exposes:

- `NewStore`
- `NewAction`
- `NewActionHandlerBuilder`

Lower-level packages remain available when you want more explicit imports:

- `actions`
- `handlers`
- `events`
- `ioc`

### Modules

Complex applications can be composed from multiple handlers added to the same store.

For map-based state, a handler may declare:

- an initial state
- a selector

This allows a handler to operate on a dedicated slice inside the shared state map.

## Examples

Each top-level example under [`examples/`](/home/baco/workspace/go-redux/examples) is documented independently.

- [`examples/todo-ddd/README.md`](/home/baco/workspace/go-redux/examples/todo-ddd/README.md): small DDD-oriented example with multiple handlers and filtered projections
- [`examples/game-server/README.md`](/home/baco/workspace/go-redux/examples/game-server/README.md): game state coordination with subscribers and domain services
- [`examples/event-sourced-wallet/README.md`](/home/baco/workspace/go-redux/examples/event-sourced-wallet/README.md): event sourcing and projections
- [`examples/workflow-employee-onboarding/README.md`](/home/baco/workspace/go-redux/examples/workflow-employee-onboarding/README.md): workflow orchestration and rollback
- [`examples/saga-order-fulfillment/README.md`](/home/baco/workspace/go-redux/examples/saga-order-fulfillment/README.md): saga-style distributed coordination

Several examples are independent Go modules with their own `go.mod`. The repository CI validates the root module and each example module separately.

## Testing

Root module:

```bash
go test ./...
go vet ./...
```

Example modules:

```bash
cd examples/todo-ddd && go test ./...
cd examples/game-server && go test ./...
cd examples/event-sourced-wallet && go test ./...
cd examples/workflow-employee-onboarding && go test ./...
cd examples/saga-order-fulfillment/orchestrator && go test ./...
```

Repository-level testing rules live in:

- [`.github/instructions/testing.instructions.md`](/home/baco/workspace/go-redux/.github/instructions/testing.instructions.md)

## Integration with go-infrastructure

The library is designed to work comfortably with `go-infrastructure/v2`, especially for:

- dependency injection
- logging
- event publication
- server bootstrap in example applications

The IoC helpers in [`ioc/`](/home/baco/workspace/go-redux/ioc) are intentionally small. They are meant to support composition, not to hide the library model.

## Versioning

This repository is on the `v2` module path. Releases for this line should use `v2.x.y` tags.

## Documentation

- [Architecture guide](/home/baco/workspace/go-redux/docs/architecture.md)
- [Plan and notes](/home/baco/workspace/go-redux/docs/plan.md)

## License

[MIT](/home/baco/workspace/go-redux/LICENSE)
