# Architecture Guide

## Overview

go-redux v2 is designed to work seamlessly with **Domain-Driven Design (DDD)** and **Clean Architecture** patterns. This guide explains the architectural principles and how to structure your applications.

## Core Principles

### 1. Separation of Concerns

The architecture enforces clear boundaries between layers:

- **Domain Layer**: Business entities and rules
- **Application Layer**: Use cases and orchestration
- **Infrastructure Layer**: External concerns (I/O, persistence, etc.)

### 2. Dependency Rule

Dependencies point inward:
- Infrastructure → Application → Domain
- Domain has **zero dependencies**
- Application depends only on Domain
- Infrastructure depends on both

### 3. Immutable State

State changes create new instances rather than mutating existing ones, ensuring:
- Predictable behavior
- Time-travel debugging capabilities
- Thread-safe operations

## Layer Structure

### Domain Layer

Pure business logic with no framework dependencies.

```
domain/
├── entities/          # Business objects
│   ├── counter.go
│   └── todo.go
├── services/          # Business rules
│   ├── incrementer.go
│   └── todo_service.go
└── ioc/              # Domain dependencies
    └── register.go
```

**Example:**

```go
// domain/incrementer.go
type Incrementer struct{}

func (i *Incrementer) AutoIncrement(counter *Counter) *Counter {
    if counter.Count < 10 {
        counter.Count += 1
    } else {
        counter.Count += 2
    }
    return counter
}
```

**Rules:**
- No imports from application or infrastructure
- Pure functions when possible
- Business logic only

### Application Layer

Orchestrates domain logic and manages state.

```
application/
├── state.go           # Application state definition
├── actions.go         # Action definitions
├── handlers/          # Redux handlers
│   ├── counter_handler.go
│   └── todo_handler.go
└── ioc/              # Application dependencies
    └── register.go
```

**Example:**

```go
// application/counter_handler.go
type CounterHandler struct {
    incrementer *domain.Incrementer
}

var AutoIncrementAction = redux.NewAction[struct{}]("AUTO_INCREMENT")

func NewCounterHandler(inc *domain.Incrementer) handlers.ActionHandler[CounterState] {
    h := &CounterHandler{incrementer: inc}

    return redux.NewActionHandlerBuilder[CounterState]().
        On(AutoIncrementAction, h.handleAutoIncrement).
        Build()
}

func (h *CounterHandler) handleAutoIncrement(state CounterState, _ any) CounterState {
    counter := &domain.Counter{Count: state.Count}
    result := h.incrementer.AutoIncrement(counter)
    state.Count = result.Count
    return state
}
```

**Rules:**
- Depends on domain layer
- Manages state transformations
- Orchestrates domain services
- No I/O operations

### Infrastructure Layer

External concerns and composition.

```
infrastructure/
├── subscribers/      # State change listeners
│   ├── logger_subscriber.go
│   ├── persistence_subscriber.go
│   └── analytics_subscriber.go
└── ioc/             # Infrastructure dependencies
    └── register.go
```

**Example:**

```go
// infrastructure/subscribers/logger_subscriber.go
type LoggerSubscriber struct {
    logger logs.Logger
}

func (s *LoggerSubscriber) OnStateChange(state application.GameState) {
    s.logger.Info("State changed", "players", len(state.Players))
}
```

**Rules:**
- Can depend on all layers
- Handles I/O, persistence, external services
- Composition root lives here

## IoC Container Integration

### Why IoC?

1. **Testability**: Easy to mock dependencies
2. **Flexibility**: Swap implementations
3. **Lifecycle Management**: Singleton, transient, scoped instances

### Registration Pattern

```go
// domain/ioc/register.go
func Register(container *di.Container) {
    di.RegisterSingleton(container, domain.NewIncrementer)
}

// application/ioc/register.go
func Register(container *di.Container) {
    domainIoc.Register(container)
    
    di.RegisterSingleton(container, func(inc *domain.Incrementer) handlers.ActionHandler[application.CounterState] {
        return application.NewCounterHandler(inc)
    })
}
```

### Tenant Pattern

For multiple instances of the same handler:

```go
// Register with tenant
di.RegisterSingletonTenant(container, "todos", func(svc *domain.TodoService) handlers.ActionHandler[AppState] {
    return application.NewTodoHandler(svc)
})

// Resolve with tenant
handler := di.ResolveTenant[handlers.ActionHandler[AppState]](container, "todos")
```

## Modular Architecture

### AddModule Pattern

```go
// Compose multiple handlers
store.AddModule(counterHandler)
store.AddModule(todoHandler)
store.AddModule(filterHandler)
```

Each module:
- Handles specific action types
- Independent of other modules
- Composable at runtime

### State Composition

```go
type AppState struct {
    Todos  TodoState
    Filter FilterState
}
```

Each handler modifies its slice of state.

## Subscriber Pattern

### Multiple Subscribers

```go
// Broadcast subscriber
broadcastSub := subscribers.NewBroadcastSubscriber()
broadcastCallback := broadcastSub.OnStateChange
store.Subscribe(&broadcastCallback)

// Persistence subscriber
persistenceSub := subscribers.NewPersistenceSubscriber()
persistenceCallback := persistenceSub.OnStateChange
store.Subscribe(&persistenceCallback)
```

### Subscriber Types

1. **Logging**: Track state changes
2. **Persistence**: Save to database
3. **Analytics**: Metrics and monitoring
4. **Broadcasting**: WebSocket/SSE updates

## Best Practices

### 1. Keep Domain Pure

```go
// ✅ Good - Pure domain logic
func (s *TodoService) AddTodo(todos []*Todo, nextID int, text string) ([]*Todo, int) {
    newTodo := NewTodo(nextID, text)
    return append(todos, newTodo), nextID + 1
}

// ❌ Bad - Side effects in domain
func (s *TodoService) AddTodo(db *sql.DB, text string) error {
    // Database access doesn't belong in domain
}
```

### 2. State Immutability

```go
// ✅ Good - Returns new state
func (h *Handler) handleAdd(state State, item Item) State {
    newItems := append([]Item{}, state.Items...)
    newItems = append(newItems, item)
    state.Items = newItems
    return state
}

// ❌ Bad - Mutates state
func (h *Handler) handleAdd(state State, item Item) State {
    state.Items = append(state.Items, item)
    return state
}
```

### 3. Type-Safe Actions

```go
// ✅ Good - Type-safe action
var AddTodoAction = redux.ActionBuilder[string]().
    WithType("ADD_TODO").
    Build()

store.Dispatch(AddTodoAction.With("Buy milk"))

// ❌ Bad - Untyped action
store.Dispatch(map[string]any{"type": "ADD_TODO", "text": "Buy milk"})
```

### 4. Handler Composition

```go
// ✅ Good - Single responsibility
builder := redux.NewActionHandlerBuilder[State]().
    On(AddAction, handleAdd).
    On(RemoveAction, handleRemove).
    Build()

// ❌ Bad - God handler
builder := redux.NewActionHandlerBuilder[State]().
    On(AddAction, handleEverything).
    On(RemoveAction, handleEverything).
    On(UpdateAction, handleEverything).
    Build()
```

## Testing Strategy

### Domain Layer

Test business logic in isolation:

```go
func TestIncrementer_AutoIncrement(t *testing.T) {
    service := domain.NewIncrementer()
    counter := &domain.Counter{Count: 5}
    
    result := service.AutoIncrement(counter)
    
    if result.Count != 6 {
        t.Errorf("Expected 6, got %d", result.Count)
    }
}
```

### Application Layer

Test with mock domain services:

```go
func TestCounterHandler(t *testing.T) {
    mockInc := &MockIncrementer{}
    handler := application.NewCounterHandler(mockInc)
    
    // Test reducer logic
}
```

### Integration Tests

Test full stack with real dependencies:

```go
func TestFullFlow(t *testing.T) {
    container := di.NewContainer()
    appIoc.Register(container)
    
    store := redux.NewStore(initialState, logger)
    handler := di.Resolve[handlers.ActionHandler[State]](container)
    store.AddModule(handler)
    
    // Test dispatches and state changes
}
```

## Performance Considerations

### 1. Avoid Unnecessary Copies

```go
// ✅ Good - Only copy when needed
if !changed {
    return state
}
newState := copyState(state)
return newState

// ❌ Bad - Always copies
newState := copyState(state)
return newState
```

### 2. Batch Operations

```go
// ✅ Good - Single action
store.Dispatch(BulkAddAction.With(items))

// ❌ Bad - Multiple dispatches
for _, item := range items {
    store.Dispatch(AddAction.With(item))
}
```

### 3. Subscriber Efficiency

```go
// ✅ Good - Lightweight subscriber
func (s *Subscriber) OnStateChange(state State) {
    s.queue <- state  // Queue for async processing
}

// ❌ Bad - Heavy work in subscriber
func (s *Subscriber) OnStateChange(state State) {
    s.db.Save(state)  // Blocks state updates
}
```

## Examples in Codebase

- **counter-ddd**: Basic DDD structure
- **todo-ddd**: Multiple handlers with tenant pattern
- **game-server**: Complex domain with multiple services and subscribers

## Further Reading

- [DDD by Eric Evans](https://www.domainlanguage.com/ddd/)
- [Clean Architecture by Robert Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Redux Patterns](https://redux.js.org/style-guide/)
