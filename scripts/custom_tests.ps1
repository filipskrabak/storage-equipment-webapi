param(
    [string]$BaseUrl = "http://localhost:5000/api"
)

Write-Host "Testing Storage Equipment API at $BaseUrl" -ForegroundColor Green
Write-Host "=================================================" -ForegroundColor Green

# Test variables
$equipmentId = ""
$testResults = @()

function Test-Endpoint {
    param(
        [string]$Name,
        [string]$Method,
        [string]$Url,
        [string]$Body = "",
        [int]$ExpectedStatus = 200
    )
    
    Write-Host "`nTesting: $Name" -ForegroundColor Yellow
    Write-Host "$Method $Url"
    
    try {
        $headers = @{ "Content-Type" = "application/json" }
        
        if ($Body) {
            $response = Invoke-WebRequest -Uri $Url -Method $Method -Body $Body -Headers $headers -UseBasicParsing
        } else {
            $response = Invoke-WebRequest -Uri $Url -Method $Method -Headers $headers -UseBasicParsing
        }
        
        $statusCode = $response.StatusCode
        
        if ($statusCode -eq $ExpectedStatus) {
            Write-Host "✅ PASS - Status: $statusCode" -ForegroundColor Green
            
            # Parse JSON response if content type is JSON
            $responseData = $null
            if ($response.Content) {
                try {
                    $responseData = $response.Content | ConvertFrom-Json
                } catch {
                    $responseData = $response.Content
                }
            }
            
            $script:testResults += @{ Test = $Name; Status = "PASS"; Response = $responseData }
            return $responseData
        } else {
            Write-Host "❌ FAIL - Expected: $ExpectedStatus, Got: $statusCode" -ForegroundColor Red
            $script:testResults += @{ Test = $Name; Status = "FAIL"; Error = "Wrong status code: Expected $ExpectedStatus, Got $statusCode" }
        }
    }
    catch {
        # Extract status code from error if available
        $statusCode = "Unknown"
        if ($_.Exception.Response) {
            $statusCode = [int]$_.Exception.Response.StatusCode
        }
        
        Write-Host "❌ FAIL - Status: $statusCode, Error: $($_.Exception.Message)" -ForegroundColor Red
        $script:testResults += @{ Test = $Name; Status = "FAIL"; Error = "Status: $statusCode, $($_.Exception.Message)" }
    }
}

function Test-EndpointExpectError {
    param(
        [string]$Name,
        [string]$Method,
        [string]$Url,
        [string]$Body = "",
        [int]$ExpectedStatus = 400
    )
    
    Write-Host "`nTesting: $Name" -ForegroundColor Yellow
    Write-Host "$Method $Url"
    
    try {
        $headers = @{ "Content-Type" = "application/json" }
        
        if ($Body) {
            $response = Invoke-WebRequest -Uri $Url -Method $Method -Body $Body -Headers $headers -UseBasicParsing -ErrorAction Stop
        } else {
            $response = Invoke-WebRequest -Uri $Url -Method $Method -Headers $headers -UseBasicParsing -ErrorAction Stop
        }
        
        $statusCode = $response.StatusCode
        Write-Host "❌ FAIL - Expected error but got success with status: $statusCode" -ForegroundColor Red
        $script:testResults += @{ Test = $Name; Status = "FAIL"; Error = "Expected error status $ExpectedStatus but got success with status $statusCode" }
    }
    catch {
        $statusCode = "Unknown"
        if ($_.Exception.Response) {
            $statusCode = [int]$_.Exception.Response.StatusCode
        }
        
        if ($statusCode -eq $ExpectedStatus) {
            Write-Host "✅ PASS - Expected error status: $statusCode" -ForegroundColor Green
            $script:testResults += @{ Test = $Name; Status = "PASS"; Response = "Expected error received with status $statusCode" }
        } else {
            Write-Host "❌ FAIL - Expected: $ExpectedStatus, Got: $statusCode" -ForegroundColor Red
            $script:testResults += @{ Test = $Name; Status = "FAIL"; Error = "Expected error status $ExpectedStatus but got $statusCode" }
        }
    }
}

# Test 1: Get all equipment (should work even with empty database)
$allEquipment = Test-Endpoint -Name "Get All Equipment" -Method "GET" -Url "$BaseUrl/equipment"

# Test 2: Create new equipment
$newEquipment = @{
    name = "Test MRI Scanner"
    serialNumber = "TEST-MRI-001"
    manufacturer = "Test Manufacturer"
    model = "Test Model X1"
    installationDate = "2023-01-15"
    location = "Test Room 101"
    serviceInterval = 90
    lastService = "2023-03-15"
    lifeExpectancy = 10
    status = "operational"
    notes = "Test equipment for API testing"
} | ConvertTo-Json -Depth 3

$createdEquipment = Test-Endpoint -Name "Create Equipment" -Method "POST" -Url "$BaseUrl/equipment" -Body $newEquipment -ExpectedStatus 201

if ($createdEquipment -and $createdEquipment.id) {
    $script:equipmentId = $createdEquipment.id
    Write-Host "Created equipment with ID: $equipmentId" -ForegroundColor Cyan
    
    # Test 3: Get equipment by ID
    Test-Endpoint -Name "Get Equipment by ID" -Method "GET" -Url "$BaseUrl/equipment/$equipmentId"
    
    # Test 4: Update equipment
    $updateData = @{
        name = "Updated MRI Scanner"
        location = "Updated Room 102"
        status = "in_repair"
        notes = "Updated via API test"
    } | ConvertTo-Json -Depth 3
    
    Test-Endpoint -Name "Update Equipment" -Method "PUT" -Url "$BaseUrl/equipment/$equipmentId" -Body $updateData
    
    # Test 5: Get updated equipment to verify changes
    $updatedEquipment = Test-Endpoint -Name "Verify Update" -Method "GET" -Url "$BaseUrl/equipment/$equipmentId"
    
    if ($updatedEquipment -and $updatedEquipment.name -eq "Updated MRI Scanner") {
        Write-Host "✅ Update verification passed" -ForegroundColor Green
        $script:testResults += @{ Test = "Update Verification"; Status = "PASS"; Response = "Data updated correctly" }
    } else {
        Write-Host "❌ Update verification failed" -ForegroundColor Red
        $script:testResults += @{ Test = "Update Verification"; Status = "FAIL"; Error = "Updated data does not match expected values" }
    }
}

# Test 6: Error scenarios
Test-EndpointExpectError -Name "Get Non-existent Equipment" -Method "GET" -Url "$BaseUrl/equipment/non-existent-id" -ExpectedStatus 404

Test-EndpointExpectError -Name "Create Equipment with Missing Fields" -Method "POST" -Url "$BaseUrl/equipment" -Body '{"name":"Incomplete"}' -ExpectedStatus 400

Test-EndpointExpectError -Name "Create Equipment with Invalid JSON" -Method "POST" -Url "$BaseUrl/equipment" -Body 'invalid json' -ExpectedStatus 400

Test-EndpointExpectError -Name "Update Non-existent Equipment" -Method "PUT" -Url "$BaseUrl/equipment/non-existent-id" -Body '{"name":"Test"}' -ExpectedStatus 404

Test-EndpointExpectError -Name "Delete Non-existent Equipment" -Method "DELETE" -Url "$BaseUrl/equipment/non-existent-id" -ExpectedStatus 404

# Test 7: Delete the created equipment
if ($script:equipmentId) {
    # First verify it exists
    Test-Endpoint -Name "Verify Equipment Exists Before Delete" -Method "GET" -Url "$BaseUrl/equipment/$equipmentId"
    
    # Delete it
    Test-Endpoint -Name "Delete Equipment" -Method "DELETE" -Url "$BaseUrl/equipment/$equipmentId" -ExpectedStatus 204
    
    # Verify it's gone
    Test-EndpointExpectError -Name "Verify Equipment Deleted" -Method "GET" -Url "$BaseUrl/equipment/$equipmentId" -ExpectedStatus 404
}

# Test 8: Stress test - Create multiple equipment items
Write-Host "`nRunning stress test - Creating 5 equipment items..." -ForegroundColor Yellow
$createdIds = @()

for ($i = 1; $i -le 5; $i++) {
    $stressEquipment = @{
        name = "Stress Test Equipment $i"
        serialNumber = "STRESS-$i"
        manufacturer = "Test Manufacturer"
        location = "Test Room $i"
    } | ConvertTo-Json -Depth 3
    
    $result = Test-Endpoint -Name "Stress Test Create $i" -Method "POST" -Url "$BaseUrl/equipment" -Body $stressEquipment -ExpectedStatus 201
    
    if ($result -and $result.id) {
        $createdIds += $result.id
    }
}

# Get all equipment to see the new items
$finalEquipmentList = Test-Endpoint -Name "Get All Equipment After Stress Test" -Method "GET" -Url "$BaseUrl/equipment"

# Clean up stress test equipment
Write-Host "`nCleaning up stress test equipment..." -ForegroundColor Yellow
foreach ($id in $createdIds) {
    if ($id) {
        Test-Endpoint -Name "Cleanup Equipment $id" -Method "DELETE" -Url "$BaseUrl/equipment/$id" -ExpectedStatus 204
    }
}

# Test 9: Additional validation tests
Write-Host "`nRunning additional validation tests..." -ForegroundColor Yellow

# Test with edge case data
$edgeCaseEquipment = @{
    name = "Equipment with Special Characters !@#$%^&*()"
    serialNumber = "EDGE-CASE-123-456"
    manufacturer = "Test Manufacturer with Spaces"
    model = "Model 2023.1"
    installationDate = "2023-12-31"
    location = "Room 999"
    serviceInterval = 1
    lastService = "2023-12-30"
    lifeExpectancy = 1
    status = "operational"
    notes = "This is a test with edge case data including special characters and boundary values."
} | ConvertTo-Json -Depth 3

$edgeCaseResult = Test-Endpoint -Name "Create Equipment with Edge Case Data" -Method "POST" -Url "$BaseUrl/equipment" -Body $edgeCaseEquipment -ExpectedStatus 201

if ($edgeCaseResult -and $edgeCaseResult.id) {
    # Clean up edge case equipment
    Test-Endpoint -Name "Cleanup Edge Case Equipment" -Method "DELETE" -Url "$BaseUrl/equipment/$($edgeCaseResult.id)" -ExpectedStatus 204
}

# Summary
Write-Host "`n=================================================" -ForegroundColor Green
Write-Host "TEST SUMMARY" -ForegroundColor Green
Write-Host "=================================================" -ForegroundColor Green

$passCount = ($testResults | Where-Object { $_.Status -eq "PASS" }).Count
$failCount = ($testResults | Where-Object { $_.Status -eq "FAIL" }).Count
$totalCount = $testResults.Count

Write-Host "Total Tests: $totalCount" -ForegroundColor White
Write-Host "Passed: $passCount" -ForegroundColor Green
Write-Host "Failed: $failCount" -ForegroundColor Red

if ($failCount -gt 0) {
    Write-Host "`nFailed Tests:" -ForegroundColor Red
    $testResults | Where-Object { $_.Status -eq "FAIL" } | ForEach-Object {
        Write-Host "  - $($_.Test): $($_.Error)" -ForegroundColor Red
    }
}

$successRate = if ($totalCount -gt 0) { [math]::Round(($passCount / $totalCount) * 100, 2) } else { 0 }
Write-Host "`nSuccess Rate: $successRate%" -ForegroundColor $(if ($successRate -eq 100) { "Green" } elseif ($successRate -ge 80) { "Yellow" } else { "Red" })

if ($successRate -eq 100) {
    Write-Host "`n🎉 All tests passed! Your API is working correctly." -ForegroundColor Green
    exit 0
} elseif ($successRate -ge 80) {
    Write-Host "`n⚠️  Most tests passed, but some issues were found. Please review the failures." -ForegroundColor Yellow
    exit 1
} else {
    Write-Host "`n❌ Many tests failed. Please check the API implementation and server status." -ForegroundColor Red
    exit 1
}