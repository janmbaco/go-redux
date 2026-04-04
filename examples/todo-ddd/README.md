# Todo DDD Example

This example shows how to use `go-redux` in a small application with:

- typed actions
- multiple handlers in the same store
- domain services separated from application handlers
- a simple filtered view over shared state

## What It Demonstrates

- a domain package with business rules for todos
- an application layer that maps actions to state transitions
- composition through two handlers for todo management and filter management
- integration with `go-infrastructure/v2/dependencyinjection`

## Structure

- `application/`: DTOs, actions, handlers, IoC wiring
- `domain/`: entities and business rules
- `infrastructure/`: example-specific wiring helpers
- `main.go`: composition root and demo execution

## Run

```bash
cd examples/todo-ddd
go run .
```

## Test

```bash
cd examples/todo-ddd
go test ./...
```

## Notes

This example is intentionally small. Its purpose is to show a clean baseline for handler composition and state slicing before moving to more operational examples such as workflows, event sourcing, or sagas.
