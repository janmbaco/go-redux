#!/bin/bash

echo "=========================================="
echo "  Event-Sourced Wallet Testbed"
echo "=========================================="

BASE_URL="http://localhost:8080"

# Check health
echo ""
echo "1. Health Check"
curl -s $BASE_URL/health | jq

# Create wallet 1
echo ""
echo "2. Creating Wallet 1"
WALLET1=$(curl -s -X POST $BASE_URL/wallets \
  -H "Content-Type: application/json" \
  -d '{"owner_id":"user-1","currency":"USD","initial_balance":1000}' | jq -r '.id')
echo "Wallet 1 ID: $WALLET1"

# Create wallet 2
echo ""
echo "3. Creating Wallet 2"
WALLET2=$(curl -s -X POST $BASE_URL/wallets \
  -H "Content-Type: application/json" \
  -d '{"owner_id":"user-2","currency":"USD","initial_balance":500}' | jq -r '.id')
echo "Wallet 2 ID: $WALLET2"

# Deposit to wallet 1
echo ""
echo "4. Depositing \$500 to Wallet 1"
curl -s -X POST $BASE_URL/wallets/$WALLET1/deposit \
  -H "Content-Type: application/json" \
  -d '{"amount":500,"description":"Salary"}' | jq

# Withdraw from wallet 1
echo ""
echo "5. Withdrawing \$200 from Wallet 1"
curl -s -X POST $BASE_URL/wallets/$WALLET1/withdraw \
  -H "Content-Type: application/json" \
  -d '{"amount":200,"description":"ATM"}' | jq

# Transfer from wallet 1 to wallet 2
echo ""
echo "6. Transferring \$300 from Wallet 1 to Wallet 2"
curl -s -X POST $BASE_URL/wallets/$WALLET1/transfer \
  -H "Content-Type: application/json" \
  -d "{\"to_wallet_id\":\"$WALLET2\",\"amount\":300,\"description\":\"Payment\"}" | jq

# Get balances from projection
echo ""
echo "7. Balance Check (CQRS Projection)"
echo "Wallet 1 Balance:"
curl -s $BASE_URL/wallets/$WALLET1/balance | jq
echo ""
echo "Wallet 2 Balance:"
curl -s $BASE_URL/wallets/$WALLET2/balance | jq

# Get event history
echo ""
echo "8. Event History - Wallet 1"
curl -s $BASE_URL/wallets/$WALLET1/events | jq

# Simulate 100 deposits to trigger snapshot
echo ""
echo "9. Creating 100 deposits to trigger snapshot..."
for i in {1..100}; do
  curl -s -X POST $BASE_URL/wallets/$WALLET1/deposit \
    -H "Content-Type: application/json" \
    -d '{"amount":10,"description":"Test deposit '$i'"}' > /dev/null
  
  if [ $((i % 20)) -eq 0 ]; then
    echo "  Deposited $i times..."
  fi
done
echo "  All 100 deposits completed"

# Check final event count
echo ""
echo "10. Event Count After 100 Operations"
curl -s $BASE_URL/wallets/$WALLET1/events | jq '{wallet_id, event_count}'

# Replay events
echo ""
echo "11. Replaying All Events"
START_TIME=$(date +%s%3N)
curl -s -X POST $BASE_URL/wallets/$WALLET1/replay | jq
END_TIME=$(date +%s%3N)
REPLAY_TIME=$((END_TIME - START_TIME))
echo "Replay completed in ${REPLAY_TIME}ms"

# Final balance
echo ""
echo "12. Final Balance (should match after replay)"
curl -s $BASE_URL/wallets/$WALLET1/balance | jq

echo ""
echo "=========================================="
echo "  Testbed Complete!"
echo "=========================================="
echo ""
echo "Expected Results:"
echo "  - Wallet 1: Initial 1000 + 500 - 200 - 300 + (100 * 10) = 2000"
echo "  - Wallet 2: Initial 500 + 300 = 800"
echo "  - Event replay completed in <100ms"
echo "  - Snapshot created at event 100"
