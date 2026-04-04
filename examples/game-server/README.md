# Game Server Example

This example shows how to model real-time game state coordination with `go-redux`.

It combines:

- typed actions for game events
- domain services for movement and combat rules
- multiple handlers in the same store
- multiple subscribers reacting to state changes

## What It Demonstrates

- player lifecycle: join, move, attack, heal, leave
- separation between domain rules and application handlers
- subscriber fan-out for broadcast, persistence, and analytics concerns
- store composition with dependency injection

## Structure

- `application/`: actions, payloads, handlers, DTOs, IoC wiring
- `domain/`: player model and game rules
- `infrastructure/subscribers/`: reaction layer for state changes
- `main.go`: composition root and executable demo

## Run

```bash
cd examples/game-server
go run .
```

## Test

```bash
cd examples/game-server
go test ./...
```

## Notes

This example is still a console demo, not a production game server. Its value is in showing how `go-redux` can coordinate a mutable runtime state while keeping business rules and side effects separated.
