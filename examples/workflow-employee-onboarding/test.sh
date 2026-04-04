#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=========================================="
echo "  Workflow Employee Onboarding Testbed"
echo "=========================================="
echo ""

echo "1. Health Check"
curl -s "$BASE_URL/health" | jq .
echo ""

echo "2. Creating Employee Onboarding (Success Case)"
EMPLOYEE1=$(curl -s -X POST "$BASE_URL/onboarding" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Alice",
    "last_name": "Johnson",
    "email": "alice.johnson@company.com",
    "department": "Engineering",
    "position": "Senior Developer",
    "hire_date": "2025-01-15"
  }' | jq -r '.employee_id')

echo "Employee 1 ID: $EMPLOYEE1"
echo ""

echo "3. Wait for workflow to complete (parallel execution)..."
sleep 3
echo ""

echo "4. Get Employee Status"
curl -s "$BASE_URL/employees?id=$EMPLOYEE1" | jq .
echo ""

echo "5. Get Workflow Status"
curl -s "$BASE_URL/workflow/status?id=$EMPLOYEE1" | jq .
echo ""

echo "6. Creating Multiple Employees (Test Parallel Workflows)"
echo "Creating 5 employees simultaneously..."

for i in {1..5}; do
  curl -s -X POST "$BASE_URL/onboarding" \
    -H "Content-Type: application/json" \
    -d "{
      \"first_name\": \"Employee\",
      \"last_name\": \"$i\",
      \"email\": \"employee$i@company.com\",
      \"department\": \"Operations\",
      \"position\": \"Analyst\",
      \"hire_date\": \"2025-01-20\"
    }" | jq -r '.employee_id' &
done

wait
echo "All onboarding workflows started"
echo ""

echo "7. Wait for parallel workflows..."
sleep 4
echo ""

echo "8. List All Employees"
curl -s "$BASE_URL/employees" | jq .
echo ""

echo "9. Testing Workflow Failure & Rollback"
echo "Note: With 20% failure rate, some steps may fail and trigger rollback"
EMPLOYEE_FAIL=$(curl -s -X POST "$BASE_URL/onboarding" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Bob",
    "last_name": "Smith",
    "email": "bob.smith@company.com",
    "department": "Sales",
    "position": "Account Manager",
    "hire_date": "2025-01-18"
  }' | jq -r '.employee_id')

echo "Employee (potential failure) ID: $EMPLOYEE_FAIL"
sleep 3
echo ""

echo "10. Check Failed Employee Status"
curl -s "$BASE_URL/employees?id=$EMPLOYEE_FAIL" | jq .
echo ""

echo "11. Check Workflow Status (may show FAILED or ROLLED_BACK)"
curl -s "$BASE_URL/workflow/status?id=$EMPLOYEE_FAIL" | jq .
echo ""

echo "=========================================="
echo "  Testbed Complete!"
echo "=========================================="
echo ""
echo "Expected Results:"
echo "  - Employee 1 (Alice): Status COMPLETED"
echo "  - Steps executed: HR_SETUP → IT_PROVISIONING/TRAINING (parallel) → BADGE"
echo "  - 5 parallel employees: Mix of COMPLETED/IN_PROGRESS"
echo "  - Some workflows may FAIL and trigger automatic ROLLBACK"
echo "  - Rollback reverses steps in order: BADGE → TRAINING/IT → HR_SETUP"
echo ""
echo "Key Features Demonstrated:"
echo "  1. Multi-step workflow with dependencies"
echo "  2. Parallel execution of independent steps"
echo "  3. Automatic rollback on failure"
echo "  4. Workflow state persistence"
echo "  5. Multiple concurrent workflows"
