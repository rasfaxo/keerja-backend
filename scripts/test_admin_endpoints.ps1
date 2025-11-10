# Configuration
$ApiUrl = "http://localhost:8080/api/v1"
$AdminEmail = "admin@keerja.com"  # Replace with your admin user
$AdminPassword = "admin123"       # Replace with actual password

# Function to get auth token
function Get-AuthToken {
    $loginBody = @{
        email = $AdminEmail
        password = $AdminPassword
    } | ConvertTo-Json
    
    $response = Invoke-RestMethod -Method Post -Uri "$ApiUrl/auth/login" `
        -ContentType "application/json" -Body $loginBody
    return $response.data.token
}

# Function to test endpoints
function Test-Endpoint {
    param (
        [string]$Method,
        [string]$Endpoint,
        [string]$Payload,
        [string]$Token,
        [string]$Description
    )

    Write-Host "`nTesting $Description" -ForegroundColor Green
    Write-Host "Method: $Method"
    Write-Host "Endpoint: $Endpoint"
    
    if ($Payload) {
        Write-Host "Payload: $Payload"
    }

    $headers = @{
        "Authorization" = "Bearer $Token"
        "Content-Type" = "application/json"
    }

    try {
        if ($Payload) {
            $response = Invoke-RestMethod -Method $Method -Uri "$ApiUrl$Endpoint" `
                -Headers $headers -Body $Payload
        } else {
            $response = Invoke-RestMethod -Method $Method -Uri "$ApiUrl$Endpoint" `
                -Headers $headers
        }
        Write-Host "Response:"
        $response | ConvertTo-Json -Depth 10
    }
    catch {
        Write-Host "Error: $_" -ForegroundColor Red
        Write-Host $_.ErrorDetails.Message
    }
    
    Write-Host "Test completed`n" -ForegroundColor Green
}

# Get auth token
Write-Host "Getting authentication token..."
$token = Get-AuthToken

if (-not $token) {
    Write-Host "Failed to get authentication token" -ForegroundColor Red
    exit 1
}

Write-Host "Successfully got authentication token" -ForegroundColor Green

# Test Province endpoints
Write-Host "`nTesting Province Endpoints" -ForegroundColor Green
Test-Endpoint -Method "POST" -Endpoint "/admin/master/provinces" `
    -Payload '{"code":"TEST","name":"Test Province"}' `
    -Token $token -Description "Create Province"

Test-Endpoint -Method "GET" -Endpoint "/admin/master/provinces" `
    -Token $token -Description "List Provinces"

# Test City endpoints
Write-Host "`nTesting City Endpoints" -ForegroundColor Green
Test-Endpoint -Method "POST" -Endpoint "/admin/master/cities" `
    -Payload '{"name":"Test City","provinceId":1}' `
    -Token $token -Description "Create City"

Test-Endpoint -Method "GET" -Endpoint "/admin/master/cities?province_id=1" `
    -Token $token -Description "List Cities"

# Test District endpoints
Write-Host "`nTesting District Endpoints" -ForegroundColor Green
Test-Endpoint -Method "POST" -Endpoint "/admin/master/districts" `
    -Payload '{"name":"Test District","cityId":1}' `
    -Token $token -Description "Create District"

Test-Endpoint -Method "GET" -Endpoint "/admin/master/districts?city_id=1" `
    -Token $token -Description "List Districts"

# Test Industry endpoints
Write-Host "`nTesting Industry Endpoints" -ForegroundColor Green
Test-Endpoint -Method "POST" -Endpoint "/admin/master/industries" `
    -Payload '{"name":"Test Industry"}' `
    -Token $token -Description "Create Industry"

Test-Endpoint -Method "GET" -Endpoint "/admin/master/industries" `
    -Token $token -Description "List Industries"

# Test Job Type endpoints
Write-Host "`nTesting Job Type Endpoints" -ForegroundColor Green
Test-Endpoint -Method "POST" -Endpoint "/admin/master/job-types" `
    -Payload '{"code":"TEST","name":"Test Job Type","order":1}' `
    -Token $token -Description "Create Job Type"

Test-Endpoint -Method "GET" -Endpoint "/admin/master/job-types" `
    -Token $token -Description "List Job Types"

# Test Company Size endpoints
Write-Host "`nTesting Company Size Endpoints" -ForegroundColor Green
Test-Endpoint -Method "POST" -Endpoint "/admin/meta/company-sizes" `
    -Payload '{"label":"Test Size","minEmployees":1,"maxEmployees":10}' `
    -Token $token -Description "Create Company Size"

Test-Endpoint -Method "GET" -Endpoint "/admin/meta/company-sizes" `
    -Token $token -Description "List Company Sizes"

Write-Host "`nAll tests completed" -ForegroundColor Green