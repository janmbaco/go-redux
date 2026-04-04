# Saga Order Fulfillment Example

Distributed transaction orchestration using Redux pattern with automatic compensation for failures.

## Pattern Overview

**Saga Pattern** coordinates long-running distributed transactions by breaking them into smaller, isolated steps with compensating transactions for rollback.

### When to Use
- Order processing across payment, inventory, and shipping
- Multi-step workflows requiring consistency without distributed locks
- Systems where eventual consistency is acceptable
- Need automatic rollback on failures

### vs DTM
DTM requires external saga coordinator service, go-redux uses in-memory Redux store with simpler architecture.

## Architecture

```
POST /orders
     ↓
   Redux Store (OrderState)
     ↓
   Saga Orchestrator
     ↓
   ┌─────────────┬──────────────┬──────────────┐
   │   Payment   │  Inventory   │   Shipping   │
   │   Service   │   Service    │   Service    │
   └─────────────┴──────────────┴──────────────┘
        ↓              ↓              ↓
   30% Random    30% Random    30% Random
    Failures      Failures      Failures
        ↓              ↓              ↓
   If failure → Automatic Compensation
                (Refund, Release, Cancel)
```

### Components

**Order Entity** (`internal/order/order.go`)
- Tracks saga state (PENDING → PROCESSING → COMPLETED/FAILED/ROLLED_BACK)
- Maintains completed steps for compensation
- Correlation ID for distributed tracing

**Saga Orchestrator** (`internal/order/saga.go`)
- Executes 3-step workflow: Payment → Inventory → Shipping
- Automatic compensation on failure (reverse order)
- Logs all operations with correlation IDs

**Redux Store** (`internal/order/store.go`)
- Immutable state management for orders
- 7 action types: CREATE, UPDATE, COMPLETE, FAIL, ROLLBACK, STEP_COMPLETED, STEP_FAILED
- Statistics tracking (pending, processing, completed, failed)

**Mock Services** (30% random failures)
- Payment: Simulates payment gateway
- Inventory: Simulates stock reservation
- Shipping: Simulates carrier integration

## Running

```bash
cd examples/saga-order-fulfillment

# Create a local env file with demo credentials
cp .env.example .env

# Start the full example
docker compose up --build
```

Server starts on `http://localhost:8080`

The repository does not track real credentials. Local database passwords are expected through `examples/saga-order-fulfillment/.env`, which is ignored by Git.

For standalone development, run the orchestrator and each service from their own module directories with the required database environment variables set locally.

## API Endpoints

### Create Order
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "cust-123",
    "items": [
      {"product_id": "prod-1", "name": "Laptop", "quantity": 1, "price": 999.99},
      {"product_id": "prod-2", "name": "Mouse", "quantity": 2, "price": 25.50}
    ]
  }'
```

Response (success):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "customer_id": "cust-123",
  "status": "PENDING",
  "total_amount": 1050.99,
  "correlation_id": "660e8400-e29b-41d4-a716-446655440001"
}
```

Response (failure with compensation):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "ROLLED_BACK",
  "failed_step": "INVENTORY",
  "failure_reason": "inventory reservation failed: out of stock",
  "completed_steps": ["PAYMENT"]
}
```

### Get Order
```bash
curl http://localhost:8080/orders/550e8400-e29b-41d4-a716-446655440000
```

### List All Orders
```bash
curl http://localhost:8080/orders
```

### Get Statistics
```bash
curl http://localhost:8080/stats
```

Response:
```json
{
  "total_orders": 50,
  "pending_count": 5,
  "processing_count": 10,
  "completed_count": 30,
  "failed_count": 5
}
```

## Saga Workflow

### Success Path
```
1. POST /orders
2. Redux: CREATE_ORDER action
3. Saga starts: Status = PROCESSING
4. Step 1: Payment processed → PAY-123456
5. Redux: STEP_COMPLETED (PAYMENT)
6. Step 2: Inventory reserved → INV-789012
7. Redux: STEP_COMPLETED (INVENTORY)
8. Step 3: Shipment created → SHIP-345678
9. Redux: STEP_COMPLETED (SHIPPING)
10. Redux: COMPLETE_ORDER
11. Status = COMPLETED
```

### Failure with Compensation
```
1. POST /orders
2. Redux: CREATE_ORDER action
3. Saga starts: Status = PROCESSING
4. Step 1: Payment processed → PAY-123456
5. Redux: STEP_COMPLETED (PAYMENT)
6. Step 2: Inventory reservation FAILS (30% random)
7. Redux: STEP_FAILED (INVENTORY)
8. Compensation begins:
   - Refund payment PAY-123456
9. Redux: ROLLBACK_ORDER
10. Status = ROLLED_BACK
```

## Key Features

### 1. Automatic Compensation
If any step fails, orchestrator automatically executes compensating transactions in reverse order:
- Shipping failure → Cancel shipment, release inventory, refund payment
- Inventory failure → Release inventory, refund payment
- Payment failure → No compensation needed

### 2. Correlation ID Tracking
Every order has unique correlation ID for distributed tracing across services and logs.

### 3. Redux Immutability
All state transitions immutable, enabling time-travel debugging and audit trails.

### 4. Random Failures (Demo)
Each service has 30% failure rate to demonstrate compensation logic:
```go
if s.rng.Float64() < 0.30 {
    return "", fmt.Errorf("payment failed: insufficient funds")
}
```

### 5. go-infrastructure/v2 Integration
- **server**: HTTP server with clean routing
- **logs**: Structured logging with correlation IDs
- **dependencyinjection**: IoC container for service wiring
- **errors**: Standardized error handling

## Testing Saga Behavior

Run multiple orders to observe compensation:

```bash
# Create 10 orders (expect ~3 failures with compensation)
for i in {1..10}; do
  curl -X POST http://localhost:8080/orders \
    -H "Content-Type: application/json" \
    -d "{\"customer_id\":\"cust-$i\",\"items\":[{\"product_id\":\"prod-1\",\"name\":\"Item\",\"quantity\":1,\"price\":100}]}"
  echo ""
done

# Check statistics
curl http://localhost:8080/stats
```

## Logs Example

```
[INFO] [Saga:660e8400] Starting saga for order 550e8400
[INFO] [Saga:660e8400] Executing payment step
[INFO] [Payment] Processing payment for order 550e8400 (amount: 1050.99)
[INFO] [Payment] Payment successful: PAY-1705234567890
[INFO] [Saga:660e8400] Payment completed: PAY-1705234567890
[INFO] [Saga:660e8400] Executing inventory step
[INFO] [Inventory] Reserving 3 items for order 550e8400
[ERROR] [Inventory] Reservation failed for order 550e8400: out of stock
[ERROR] [Saga:660e8400] Inventory step failed: inventory reservation failed: out of stock
[ERROR] [Saga:660e8400] Starting compensation for failed step: INVENTORY
[INFO] [Payment] Refunding payment: PAY-1705234567890
[INFO] [Payment] Refund completed: PAY-1705234567890
[INFO] [Saga:660e8400] Payment refunded: PAY-1705234567890
[INFO] [Saga:660e8400] Compensation completed
```

## Comparison with Alternatives

| Feature | go-redux Saga | DTM | Temporal |
|---------|---------------|-----|----------|
| Setup Complexity | Low (in-memory) | Medium (external service) | High (cluster) |
| Compensation | Automatic | Manual | Automatic |
| Type Safety | Go generics | Interface-based | Interface-based |
| Observability | Redux DevTools + logs | DTM dashboard | Temporal UI |
| Learning Curve | Low (Redux) | Medium | High |
| Best For | Greenfield/simple | Polyglot | Complex workflows |

## Production Hardening

The example is complete as an executable demonstration of saga orchestration with compensation. The following items are not required for the demo to work, but would be the natural next steps for a production-oriented implementation:

1. **Persist orchestrator state**: the microservices already persist their own records, but the orchestrator state still lives in memory.
2. **Add retries and backoff**: transient network or service failures should not fail immediately without retry policy.
3. **Add idempotency guarantees**: repeated commands should be safe across retries and partial failures.
4. **Add metrics and tracing**: success rate, compensation rate, latency per step, and correlation-aware traces.
5. **Add saga-focused tests**: success, compensation, restart, and duplicate-delivery scenarios.

## Related Examples

- **event-sourced-wallet**: Event sourcing with CQRS and replay
- **workflow-employee-onboarding**: Multi-step workflow with manual approvals
- **game-server**: Real-time coordination with WebSocket

## Resources

- [Saga Pattern (Microsoft)](https://learn.microsoft.com/en-us/azure/architecture/reference-architectures/saga/saga)
- [go-infrastructure/v2 Docs](https://github.com/janmbaco/go-infrastructure)
- [Redux Pattern Guide](https://redux.js.org/understanding/thinking-in-redux/three-principles)
