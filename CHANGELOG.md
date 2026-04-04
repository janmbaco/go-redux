# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2024

Complete rewrite with breaking changes for Go 1.24+ using generics, DDD, and Clean Architecture.

### Added

- **DDD + Clean Architecture Support**
  - Layered architecture: Domain, Application, Infrastructure
  - Domain layer isolation with pure business logic
  - Application layer with use cases and DTOs
  - Infrastructure layer for Redux store and state management
  
- **Modular Reducer Architecture**
  - `AddModule(moduleName, reducer)` for composing state slices
  - Type-safe module isolation
  - Independent state management per module
  - Supports complex nested state structures
  
- **ActionHandlerBuilder Pattern**
  - `NewActionHandlerBuilder[S]()` for declarative reducer creation
  - `Handle(actionType, handler)` for type-safe action mapping
  - `Build()` returns composed reducer
  - Eliminates boilerplate switch statements
  - Clear separation of action handlers
  
- **Multiple Subscribers Support**
  - `Subscribe(handler *func(S))` with pointer-to-function signature
  - All registered handlers execute on state changes
  - Fixed: closure pointer collision issue
  - Each subscriber maintains unique identity
  
- **Generic types** for type-safe state management
  - `Store[S any]` - Generic store interface
  - `Action[P any]` - Generic action with typed payload
  - `Reducer[S any]` - Generic reducer interface
  - `Middleware[S any]` - Generic middleware interface
  - `Selector[S, R any]` - Generic selector function type
  
- **Functional options pattern**
  - `Option[S]` type for store configuration
  - `WithMiddleware` for adding middleware
  
- **Utility functions**
  - `NewAction[P](type, payload)` - Create typed actions
  - `ReducerFunc[S]` - Function adapter for reducers
  - `CombineReducers[S]` - Compose multiple reducers
  - `CreateSelector[S, R]` - Create memoized selectors

- **Comprehensive Examples**
  - `counter-ddd`: Simple DDD example with auto-increment
  - `todo-ddd`: Todo list with filtering and business rules
  - `game-server`: Complex game state with combat and movement
  - All examples demonstrate full DDD stack with domain tests
  
- **Comprehensive Tests**
  - Domain layer tests for all business logic
  - Application layer validated via examples
  - Infrastructure integration validated via working examples

### Changed

- **Module path**: `github.com/janmbaco/go-redux` → `github.com/janmbaco/go-redux/v2`
- **Go version requirement**: 1.16+ → 1.24+
- **Dependency**: `go-infrastructure` v1.2.5 → v2.1.0
- **Architecture**: Complete redesign following DDD and Clean Architecture
- **Reducer pattern**: Monolithic switch → Modular `AddModule()` + `ActionHandlerBuilder`
- **Subscription pattern**: `Subscribe() <-chan S` → `Subscribe(handler *func(S))`
- **Error handling**: No more `ErrorDefer` pattern, returns errors idiomatically
- **State composition**: Global state → Module-based state slices

### Removed

- **ActionsObject pattern**
  - `ActionsObject` struct
  - `ActionsObjectFactory`
  - Reflection-based action handling
  
- **BusinessParam pattern**
  - `BusinessParam` struct
  - `BusinessParamBuilder`
  - `BusinessParamFactory`
  - `.On()` method for action handlers
  - `.SetActionsLogicByObject()`
  
- **StateManagement factory**
  - `StateManagement` interface
  - `StateManagementFactory`
  
- **Event system**
  - `StoreSubscribeEvent`
  - `StoreSubscribeEventHandler`
  - `SelectorSubscribeEvent`
  - `SelectorSubscribeEventHandler`
  
- **IoC container patterns**
  - Static `ioc/register.go`
  - Static `ioc/resolver/resolver.go`
  - Replaced with proper DI using `go-infrastructure/v2`
  
- **Dependencies**
  - `github.com/jinzhu/copier`
  - `go-infrastructure/errorschecker` (removed in v2)
  - `go-infrastructure/interfacer.ErrorDefer` (removed in v2)
  
- **Old store implementation**
  - Reflection-based dispatch
  - Channel-based subscriptions
  - Panic-based error handling

### Migration Guide

#### Before (v1)

```go
// Define actions
type CounterActions struct {
    Increment redux.Action
    Decrement redux.Action
}
actions := &CounterActions{}

// Get store from IoC
store := resolver.GetStore()
builder := resolver.GetBusinessParamBuilder()

// Build business param
builder.SetInitialState(0)
builder.SetActions(actions)
builder.On(actions.Increment, func(state int, payload int) int {
    return state + payload
})
builder.SetSelector("counter")

// Add to store
store.AddReducer(builder.GetBusinessParam())

// Dispatch
store.Dispatch(actions.Increment.With(1))

// Subscribe
subscription := func() {
    fmt.Println(store.GetState())
}
store.Subscribe(&subscription)
```

#### After (v2) - Modular DDD Approach

```go
// Domain layer - Pure business logic
type Counter struct {
    Value int
}

func (c *Counter) Increment(amount int) {
    c.Value += amount
}

// Application layer - Use case
type IncrementUseCase struct {
    store redux.Store[AppState]
}

var IncrementAction = redux.NewAction[int]("counter/increment")

func (uc *IncrementUseCase) Execute(amount int) error {
    return uc.store.Dispatch(IncrementAction.With(amount))
}

// Infrastructure layer - Reducer with ActionHandlerBuilder
builder := redux.NewActionHandlerBuilder[AppState]()
builder.On(IncrementAction, func(state AppState, amount int) AppState {
    state.Counter.Increment(amount)
    return state
})
counterReducer := builder.Build()

// Create store with modular architecture
type AppState struct {
    Counter Counter
}

store := redux.NewStore(AppState{}, logger)
defer store.Close()

store.AddModule(counterReducer)

// Register in IoC container (if using DI)
container.RegisterSingleton(func() redux.Store[AppState] { return store })

// Subscribe with function pointer
handler := func(state AppState) {
    fmt.Printf("Counter: %d\n", state.Counter.Value)
}
store.Subscribe(&handler)

// Dispatch through use case
useCase := &IncrementUseCase{store: store}
useCase.Execute(5)
```

#### Modular State Example

```go
// Compose state from independent modules
type AppState struct {
    Counter CounterState
    Todos   TodoState
    Game    GameState
}

// Create independent reducers
counterReducer := buildCounterReducer()
todoReducer := buildTodoReducer()
gameReducer := buildGameReducer()

// Add modules to store
rootReducer := redux.ReducerFunc[AppState](func(state AppState, action any) (AppState, error) {
    // Module reducers operate on state slices
    var err error
    state.Counter, err = counterReducer.Reduce(state.Counter, action)
    if err != nil { return state, err }
    
    state.Todos, err = todoReducer.Reduce(state.Todos, action)
    if err != nil { return state, err }
    
    state.Game, err = gameReducer.Reduce(state.Game, action)
    return state, err
})

store := redux.NewStore(AppState{}, rootReducer)
```

### Key Improvements

1. **DDD Architecture**: Clear domain, application, infrastructure separation
2. **Modular State**: Independent modules with `AddModule()` composition
3. **Type Safety**: Generics eliminate runtime type assertions
4. **ActionHandlerBuilder**: Declarative, boilerplate-free reducers
5. **Multiple Subscribers**: All handlers execute (fixed closure collision)
6. **Performance**: No reflection in hot path
7. **Idiomatic Go**: Error returns, function pointers, clean interfaces
8. **Testability**: Pure domain logic, isolated tests
9. **Scalability**: Module-based architecture for large applications
10. **Maintainability**: Clean Architecture, SOLID principles

### Architecture Patterns

#### Layer Structure

```
domain/          # Pure business logic, entities, services
  ├── entities      # Domain models
  ├── services      # Business rules
  └── *_test.go     # Unit tests

application/     # Use cases, DTOs, orchestration
  ├── usecases      # Application services
  └── dto           # Data transfer objects

infrastructure/  # Redux store, external adapters
  ├── state         # State structures
  ├── reducers      # Action handlers
  └── container     # IoC registration
```

#### Subscriber Pattern Fix

The v2 subscriber pattern solves closure pointer collision:

```go
// Pointer-to-function ensures unique identity
handler1 := func(s AppState) { fmt.Println("Handler 1") }
handler2 := func(s AppState) { fmt.Println("Handler 2") }
handler3 := func(s AppState) { fmt.Println("Handler 3") }

store.Subscribe(&handler1)  // All three execute
store.Subscribe(&handler2)
store.Subscribe(&handler3)

store.Dispatch(someAction)  // → Handler 1, Handler 2, Handler 3
```

Internal implementation uses `map[uintptr]func(S)` keyed by pointer address.

### Breaking Changes Summary

| v1 Component | v2 Replacement |
|--------------|----------------|
| `ActionsObject` | `Action[P any]` + domain entities |
| `BusinessParam` | `ActionHandlerBuilder[S]` |
| `BusinessParamBuilder.On()` | `builder.On(action, handler)` |
| `Store.Subscribe() <-chan S` | `Store.Subscribe(handler *func(S))` |
| `StateManagement` | Removed (use `Store` directly) |
| IoC static container | Proper DI with `go-infrastructure/v2` |
| `ErrorDefer` | Return `error` |
| Reflection dispatch | Type switch with generics |
| Monolithic reducer | Modular `AddModule()` composition |
| Global state | Layer separation (domain/app/infra) |

### Dependencies

- `github.com/janmbaco/go-infrastructure/v2` v2.1.0
  - Uses `errors.ValidateNotNil` instead of `errorschecker`
  - No `ErrorDefer` interface
  - Proper DI container support
  
### Documentation

- [README.md](README.md) - Quick start, API reference, examples
- [docs/architecture.md](docs/architecture.md) - Complete architecture guide with DDD patterns
  
## [1.x.x] - Legacy

Previous versions using reflection and factory patterns. See git history for details.

---

For full documentation and examples, see [README.md](README.md) and [docs/architecture.md](docs/architecture.md).

