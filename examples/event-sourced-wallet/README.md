# Event-Sourced Wallet

Event Sourcing + CQRS pattern with Redux for distributed systems.

## Overview

This example demonstrates **Event Sourcing** with **CQRS** using Redux pattern:

- **Event Store**: Append-only log of all wallet events
- **Event Replay**: Rebuild state from events
- **Snapshots**: Performance optimization (every 100 events)
- **CQRS**: Separate read model (balance projection)
- **Time Travel**: Query historical state
- **Audit Trail**: Complete transaction history

## Features

### Commands (Write Side)
- CreateWallet
- Deposit
- Withdraw
- Transfer

### Events
- WalletCreated
- MoneyDeposited
- MoneyWithdrawn
- MoneyTransferred
- SnapshotTaken

### Projections (Read Side)
- BalanceProjection: Fast balance queries

## Architecture

```
┌─────────────────┐
│   Commands      │
│  (HTTP API)     │
└────────┬────────┘
         │
         v
┌─────────────────┐      ┌──────────────┐
│  Redux Store    │─────>│ Event Store  │
│  (Aggregates)   │      │ (Append-Only)│
└────────┬────────┘      └──────────────┘
         │                       │
         │                       │ Replay
         v                       v
┌─────────────────┐      ┌──────────────┐
│  Projections    │<─────│ Snapshots    │
│  (Read Models)  │      │ (Every 100)  │
└─────────────────┘      └──────────────┘
```

## Running

### Local Development

```bash
cd examples/event-sourced-wallet
go run cmd/server/main.go
```

Server runs on http://localhost:8080

### Docker

```bash
docker build -t event-sourced-wallet .
docker run -p 8080:8080 event-sourced-wallet
```

## API Examples

### 1. Create Wallet

```bash
curl -X POST http://localhost:8080/wallets \
  -H "Content-Type: application/json" \
  -d '{
    "owner_id": "user-123",
    "currency": "USD",
    "initial_balance": 1000.0
  }'
```

Response:
```json
{
  "id": "wallet-abc",
  "owner_id": "user-123",
  "currency": "USD",
  "balance": 1000.0,
  "version": 1,
  "created_at": "2025-12-05T10:00:00Z"
}
```

### 2. Deposit Money

```bash
curl -X POST http://localhost:8080/wallets/wallet-abc/deposit \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 500.0,
    "description": "Salary"
  }'
```

### 3. Withdraw Money

```bash
curl -X POST http://localhost:8080/wallets/wallet-abc/withdraw \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 200.0,
    "description": "ATM withdrawal"
  }'
```

### 4. Transfer Money

```bash
curl -X POST http://localhost:8080/wallets/wallet-abc/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "to_wallet_id": "wallet-xyz",
    "amount": 100.0,
    "description": "Payment to friend"
  }'
```

### 5. Get Balance (CQRS Projection)

```bash
curl http://localhost:8080/wallets/wallet-abc/balance
```

Response:
```json
{
  "wallet_id": "wallet-abc",
  "balance": 1200.0
}
```

### 6. Get Event History

```bash
curl http://localhost:8080/wallets/wallet-abc/events
```

Response:
```json
{
  "wallet_id": "wallet-abc",
  "event_count": 4,
  "events": [
    {
      "id": "evt-1",
      "wallet_id": "wallet-abc",
      "type": "WALLET_CREATED",
      "data": { "initial_balance": 1000.0 },
      "timestamp": "2025-12-05T10:00:00Z",
      "version": 1
    },
    {
      "id": "evt-2",
      "wallet_id": "wallet-abc",
      "type": "MONEY_DEPOSITED",
      "data": { "amount": 500.0, "description": "Salary" },
      "timestamp": "2025-12-05T10:01:00Z",
      "version": 2
    }
  ]
}
```

### 7. Replay Events

```bash
curl -X POST http://localhost:8080/wallets/wallet-abc/replay
```

Rebuilds wallet state from all events.

## Testbed

Run automated test with 10,000 operations:

```bash
./test.sh
```

Expected results:
- 10,000 events replayed in <100ms
- Snapshots created every 100 events
- Balance projection matches source wallet

## Event Sourcing Benefits

### 1. Complete Audit Trail
Every transaction is recorded as an immutable event. Perfect for:
- Financial systems
- Regulatory compliance
- Debugging production issues

### 2. Time Travel
Query wallet state at any point in history:
```go
events := eventStore.GetEventsUntil(walletID, timestamp)
wallet := ReplayEvents(events)
fmt.Println(wallet.Balance) // Balance at that timestamp
```

### 3. Event Replay
Rebuild state from scratch:
- Fix bugs in production data
- Test with real event streams
- Create new projections from historical data

### 4. CQRS Performance
Separate read/write models:
- Write: Redux store (normalized, commands)
- Read: Projections (denormalized, queries)

### 5. Snapshots
Optimize replay performance:
- Store state every N events
- Replay from snapshot + remaining events
- O(N) → O(N/100) for 100-event snapshots

## vs Traditional CRUD

| Feature | Event Sourcing | CRUD |
|---------|---------------|------|
| History | ✅ Full audit trail | ❌ Lost on update |
| Time Travel | ✅ Query past state | ❌ Only current |
| Debugging | ✅ Replay events | ❌ Logs only |
| Performance | ⚠️ Need snapshots | ✅ Direct queries |
| Complexity | ⚠️ Higher | ✅ Lower |

## vs Eventuate Framework

| Feature | go-redux + Event Store | Eventuate |
|---------|----------------------|-----------|
| Language | Go | Java |
| Dependencies | Embedded | Kafka + CDC |
| Learning Curve | Low (Redux pattern) | High (framework) |
| Setup | `go run main.go` | Kafka + Debezium + DB |
| Use Case | Embedded systems | Enterprise |

## When to Use

✅ **Use Event Sourcing when:**
- You need complete audit trail (finance, healthcare)
- Historical analysis is critical
- Regulatory compliance requires immutable logs
- Debugging production with real data
- Multiple projections from same events

❌ **Don't use when:**
- Simple CRUD is sufficient
- No audit requirements
- Performance is critical (unless snapshotting)
- Team unfamiliar with event sourcing

## Performance

Tested with 10,000 operations:
- Event append: ~1ms/event
- Replay without snapshot: ~100ms
- Replay with snapshot (every 100): ~10ms
- Projection query: <1ms

## Implementation Details

### Wallet Aggregate
```go
type Wallet struct {
    ID      string
    Balance float64
    Version int64 // For optimistic locking
}

// Command: generates event
func (w *Wallet) Deposit(amount float64) (*Event, error) {
    if amount <= 0 {
        return nil, ErrInvalidAmount
    }
    
    w.Balance += amount
    w.Version++
    
    return &Event{
        Type: MoneyDeposited,
        Data: map[string]interface{}{"amount": amount},
        Version: w.Version,
    }, nil
}

// Event replay: rebuilds state
func (w *Wallet) ApplyEvent(event *Event) {
    switch event.Type {
    case MoneyDeposited:
        w.Balance += event.Data["amount"].(float64)
    case MoneyWithdrawn:
        w.Balance -= event.Data["amount"].(float64)
    }
    w.Version = event.Version
}
```

### Redux Integration
```go
// Pure reducer
func depositMoneyReducer(state WalletState, cmd DepositCommand) WalletState {
    wallet := state.Wallets[cmd.WalletID]
    event, err := wallet.Deposit(cmd.Amount, cmd.Description)
    if err != nil {
        return state // Invalid command
    }
    
    state.Events = append(state.Events, event)
    return state
}

// Action handler
builder := handlers.NewActionHandlerBuilder[WalletState]()
builder.On(DepositMoneyAction, depositMoneyReducer)
```

### Snapshotting
```go
type Snapshotter struct {
    threshold int64 // Every N events
}

func (s *Snapshotter) CheckAndCreateSnapshot(w *Wallet, eventCount int64) {
    if eventCount % s.threshold == 0 {
        snapshot := &Snapshot{
            Balance: w.Balance,
            Version: w.Version,
        }
        s.store.SaveSnapshot(w.ID, snapshot)
    }
}
```

## License

MIT
