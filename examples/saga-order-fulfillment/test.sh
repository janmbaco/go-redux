#!/bin/bash

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080"

echo "=========================================="
echo "  Saga Order Fulfillment Testbed"
echo "=========================================="
echo ""

# Health check
echo -e "${YELLOW}1. Health Check${NC}"
curl -s "$BASE_URL/health" | jq '.'
echo ""
sleep 1

# Create orders (some will fail due to 30% failure rate)
echo -e "${YELLOW}2. Creating 10 Orders${NC}"
for i in {1..10}; do
    echo -n "Creating order $i... "
    RESPONSE=$(curl -s -X POST "$BASE_URL/orders" \
        -H "Content-Type: application/json" \
        -d "{
            \"customerId\": \"CUST-$i\",
            \"productId\": \"PROD-$(( $i % 3 + 1 ))\",
            \"quantity\": $(( $i % 5 + 1 )),
            \"amount\": $(( $i * 10 + 50 )),
            \"address\": \"123 Main St, City $i\"
        }")
    ORDER_ID=$(echo $RESPONSE | jq -r '.id')
    echo -e "${GREEN}Order ID: $ORDER_ID${NC}"
    sleep 0.5
done
echo ""

# Wait for sagas to complete
echo -e "${YELLOW}3. Waiting for Sagas to Complete (10 seconds)${NC}"
for i in {10..1}; do
    echo -n "$i... "
    sleep 1
done
echo ""
echo ""

# Get all orders
echo -e "${YELLOW}4. Fetching All Orders${NC}"
curl -s "$BASE_URL/orders" | jq '.'
echo ""

# Get statistics
echo -e "${YELLOW}5. Saga Statistics${NC}"
curl -s "$BASE_URL/stats" | jq '.'
echo ""

# Get specific order details
echo -e "${YELLOW}6. Sample Order Details${NC}"
FIRST_ORDER=$(curl -s "$BASE_URL/orders" | jq -r '.orders[0].id')
if [ "$FIRST_ORDER" != "null" ]; then
    curl -s "$BASE_URL/orders/$FIRST_ORDER" | jq '.'
else
    echo "No orders found"
fi
echo ""

echo "=========================================="
echo "  Testbed Complete!"
echo "=========================================="
echo ""
echo "Expected Results:"
echo "  - ~70% orders completed (COMPLETED status)"
echo "  - ~30% orders rolled back (ROLLED_BACK status)"
echo "  - Orders show paymentId, reservationId, shipmentId when successful"
echo "  - Failed orders show errorMessage"
echo ""
