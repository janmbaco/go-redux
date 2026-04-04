Write-Host "=========================================="
Write-Host "  Workflow Employee Onboarding Testbed"
Write-Host "=========================================="
Write-Host ""

$BaseUrl = "http://localhost:8080"

Write-Host "1. Health Check"
Invoke-RestMethod -Uri "$BaseUrl/health" | ConvertTo-Json
Write-Host ""

Write-Host "2. Creating Employee Onboarding (Success Case)"
$response1 = Invoke-RestMethod -Uri "$BaseUrl/onboarding" -Method Post `
  -ContentType "application/json" `
  -Body '{
    "first_name": "Alice",
    "last_name": "Johnson",
    "email": "alice.johnson@company.com",
    "department": "Engineering",
    "position": "Senior Developer",
    "hire_date": "2025-01-15"
  }'
$employee1 = $response1.employee_id
Write-Host "Employee 1 ID: $employee1"
Write-Host ""

Write-Host "3. Wait for workflow to complete (parallel execution)..."
Start-Sleep -Seconds 3
Write-Host ""

Write-Host "4. Get Employee Status"
Invoke-RestMethod -Uri "$BaseUrl/employees?id=$employee1" | ConvertTo-Json
Write-Host ""

Write-Host "5. Get Workflow Status"
Invoke-RestMethod -Uri "$BaseUrl/workflow/status?id=$employee1" | ConvertTo-Json -Depth 10
Write-Host ""

Write-Host "6. Creating Multiple Employees (Test Parallel Workflows)"
Write-Host "Creating 5 employees simultaneously..."

$jobs = @()
for ($i=1; $i -le 5; $i++) {
  $jobs += Start-Job -ScriptBlock {
    param($url, $num)
    Invoke-RestMethod -Uri "$url/onboarding" -Method Post `
      -ContentType "application/json" `
      -Body "{
        `"first_name`": `"Employee`",
        `"last_name`": `"$num`",
        `"email`": `"employee$num@company.com`",
        `"department`": `"Operations`",
        `"position`": `"Analyst`",
        `"hire_date`": `"2025-01-20`"
      }"
  } -ArgumentList $BaseUrl, $i
}

$jobs | Wait-Job | Out-Null
Write-Host "All onboarding workflows started"
Write-Host ""

Write-Host "7. Wait for parallel workflows..."
Start-Sleep -Seconds 4
Write-Host ""

Write-Host "8. List All Employees"
Invoke-RestMethod -Uri "$BaseUrl/employees" | ConvertTo-Json -Depth 5
Write-Host ""

Write-Host "9. Testing Workflow Failure & Rollback"
Write-Host "Note: With 20% failure rate, some steps may fail and trigger rollback"
$responseFail = Invoke-RestMethod -Uri "$BaseUrl/onboarding" -Method Post `
  -ContentType "application/json" `
  -Body '{
    "first_name": "Bob",
    "last_name": "Smith",
    "email": "bob.smith@company.com",
    "department": "Sales",
    "position": "Account Manager",
    "hire_date": "2025-01-18"
  }'
$employeeFail = $responseFail.employee_id
Write-Host "Employee (potential failure) ID: $employeeFail"
Start-Sleep -Seconds 3
Write-Host ""

Write-Host "10. Check Failed Employee Status"
Invoke-RestMethod -Uri "$BaseUrl/employees?id=$employeeFail" | ConvertTo-Json
Write-Host ""

Write-Host "11. Check Workflow Status (may show FAILED or ROLLED_BACK)"
Invoke-RestMethod -Uri "$BaseUrl/workflow/status?id=$employeeFail" | ConvertTo-Json -Depth 10
Write-Host ""

Write-Host "=========================================="
Write-Host "  Testbed Complete!"
Write-Host "=========================================="
Write-Host ""
Write-Host "Expected Results:"
Write-Host "  - Employee 1 (Alice): Status COMPLETED"
Write-Host "  - Steps executed: HR_SETUP → IT_PROVISIONING/TRAINING (parallel) → BADGE"
Write-Host "  - 5 parallel employees: Mix of COMPLETED/IN_PROGRESS"
Write-Host "  - Some workflows may FAIL and trigger automatic ROLLBACK"
Write-Host "  - Rollback reverses steps in order: BADGE → TRAINING/IT → HR_SETUP"
Write-Host ""
Write-Host "Key Features Demonstrated:"
Write-Host "  1. Multi-step workflow with dependencies"
Write-Host "  2. Parallel execution of independent steps"
Write-Host "  3. Automatic rollback on failure"
Write-Host "  4. Workflow state persistence"
Write-Host "  5. Multiple concurrent workflows"
