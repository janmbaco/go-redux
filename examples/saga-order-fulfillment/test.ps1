# PowerShell Testbed for Saga Order Fulfillment

$baseUrl = "http://localhost:8080"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  Saga Order Fulfillment Testbed" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Health check
Write-Host "1. Health Check" -ForegroundColor Yellow
$response = Invoke-RestMethod -Uri "$baseUrl/health" -Method Get
$response | ConvertTo-Json -Depth 10
Write-Host ""
Start-Sleep -Seconds 1

# Create orders
Write-Host "2. Creating 10 Orders" -ForegroundColor Yellow
for ($i = 1; $i -le 10; $i++) {
    Write-Host "Creating order $i... " -NoNewline
    $body = @{
        customerId = "CUST-$i"
        productId = "PROD-$(($i % 3) + 1)"
        quantity = ($i % 5) + 1
        amount = ($i * 10) + 50
        address = "123 Main St, City $i"
    } | ConvertTo-Json

    try {
        $response = Invoke-RestMethod -Uri "$baseUrl/orders" -Method Post -Body $body -ContentType "application/json"
        Write-Host "Order ID: $($response.id)" -ForegroundColor Green
    } catch {
        Write-Host "Failed: $($_.Exception.Message)" -ForegroundColor Red
    }
    Start-Sleep -Milliseconds 500
}
Write-Host ""

# Wait for sagas
Write-Host "3. Waiting for Sagas to Complete (10 seconds)" -ForegroundColor Yellow
for ($i = 10; $i -ge 1; $i--) {
    Write-Host "$i... " -NoNewline
    Start-Sleep -Seconds 1
}
Write-Host ""
Write-Host ""

# Get all orders
Write-Host "4. Fetching All Orders" -ForegroundColor Yellow
$orders = Invoke-RestMethod -Uri "$baseUrl/orders" -Method Get
$orders | ConvertTo-Json -Depth 10
Write-Host ""

# Get statistics
Write-Host "5. Saga Statistics" -ForegroundColor Yellow
$stats = Invoke-RestMethod -Uri "$baseUrl/stats" -Method Get
$stats | ConvertTo-Json -Depth 10
Write-Host ""

# Get specific order
Write-Host "6. Sample Order Details" -ForegroundColor Yellow
if ($orders.orders.Count -gt 0) {
    $firstOrderId = $orders.orders[0].id
    $order = Invoke-RestMethod -Uri "$baseUrl/orders/$firstOrderId" -Method Get
    $order | ConvertTo-Json -Depth 10
} else {
    Write-Host "No orders found" -ForegroundColor Red
}
Write-Host ""

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  Testbed Complete!" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Expected Results:" -ForegroundColor Yellow
Write-Host "  - ~70% orders completed (COMPLETED status)"
Write-Host "  - ~30% orders rolled back (ROLLED_BACK status)"
Write-Host "  - Orders show paymentId, reservationId, shipmentId when successful"
Write-Host "  - Failed orders show errorMessage"
Write-Host ""
