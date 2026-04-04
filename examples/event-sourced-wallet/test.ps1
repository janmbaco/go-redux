Write-Host "=========================================="
Write-Host "  Event-Sourced Wallet Testbed"
Write-Host "=========================================="
Write-Host ""

Write-Host "1. Health Check"
Invoke-RestMethod -Uri "http://localhost:8080/health" | ConvertTo-Json
Write-Host ""

Write-Host "2. Creating Wallets"
$wallet1Response = Invoke-RestMethod -Uri "http://localhost:8080/wallets" -Method Post `
  -ContentType "application/json" `
  -Body '{"owner_id":"user-1","currency":"USD","initial_balance":1000}'
$wallet1 = $wallet1Response.id
Write-Host "Wallet 1 ID: $wallet1"
Write-Host ""

$wallet2Response = Invoke-RestMethod -Uri "http://localhost:8080/wallets" -Method Post `
  -ContentType "application/json" `
  -Body '{"owner_id":"user-2","currency":"USD","initial_balance":500}'
$wallet2 = $wallet2Response.id
Write-Host "Wallet 2 ID: $wallet2"
Write-Host ""

Write-Host "3. Deposit to Wallet 1"
Invoke-RestMethod -Uri "http://localhost:8080/wallets/$wallet1/deposit" -Method Post `
  -ContentType "application/json" `
  -Body '{"amount":500,"description":"Salary"}' | ConvertTo-Json
Write-Host ""

Write-Host "4. Withdraw from Wallet 1"
Invoke-RestMethod -Uri "http://localhost:8080/wallets/$wallet1/withdraw" -Method Post `
  -ContentType "application/json" `
  -Body '{"amount":200,"description":"ATM"}' | ConvertTo-Json
Write-Host ""

Write-Host "5. Transfer from Wallet 1 to Wallet 2"
Invoke-RestMethod -Uri "http://localhost:8080/wallets/$wallet1/transfer" -Method Post `
  -ContentType "application/json" `
  -Body "{`"to_wallet_id`":`"$wallet2`",`"amount`":300,`"description`":`"Payment`"}" | ConvertTo-Json
Write-Host ""

Write-Host "6. Balance Check (CQRS Projection)"
Write-Host "Wallet 1 Balance:"
Invoke-RestMethod -Uri "http://localhost:8080/wallets/$wallet1/balance" | ConvertTo-Json
Write-Host ""
Write-Host "Wallet 2 Balance:"
Invoke-RestMethod -Uri "http://localhost:8080/wallets/$wallet2/balance" | ConvertTo-Json
Write-Host ""

Write-Host "7. Event History"
$events = Invoke-RestMethod -Uri "http://localhost:8080/wallets/$wallet1/events"
Write-Host "Event Count: $($events.event_count)"
Write-Host ""

Write-Host "8. Creating 100 deposits to trigger snapshot..."
for ($i=1; $i -le 100; $i++) {
  Invoke-RestMethod -Uri "http://localhost:8080/wallets/$wallet1/deposit" -Method Post `
    -ContentType "application/json" `
    -Body "{`"amount`":10,`"description`":`"Test deposit $i`"}" | Out-Null
  
  if ($i % 20 -eq 0) {
    Write-Host "  Deposited $i times..."
  }
}
Write-Host "  All 100 deposits completed"
Write-Host ""

Write-Host "9. Event Count After 100 Operations"
$events = Invoke-RestMethod -Uri "http://localhost:8080/wallets/$wallet1/events"
Write-Host "Total Events: $($events.event_count)"
Write-Host ""

Write-Host "10. Replaying All Events"
$startTime = Get-Date
Invoke-RestMethod -Uri "http://localhost:8080/wallets/$wallet1/replay" -Method Post | ConvertTo-Json
$endTime = Get-Date
$replayTime = ($endTime - $startTime).TotalMilliseconds
Write-Host "Replay completed in ${replayTime}ms"
Write-Host ""

Write-Host "11. Final Balance"
Invoke-RestMethod -Uri "http://localhost:8080/wallets/$wallet1/balance" | ConvertTo-Json
Write-Host ""

Write-Host "=========================================="
Write-Host "  Testbed Complete!"
Write-Host "=========================================="
Write-Host ""
Write-Host "Expected Results:"
Write-Host "  - Wallet 1: Initial 1000 + 500 - 200 - 300 + (100 * 10) = 2000"
Write-Host "  - Wallet 2: Initial 500 + 300 = 800"
Write-Host "  - Event replay completed in <100ms"
Write-Host "  - Snapshot created at event 100"

