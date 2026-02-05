# PowerShell test script for validator
Write-Host "=== Testing Validator in PowerShell ===" -ForegroundColor Cyan
Write-Host ""

Write-Host "Test 1: Code 'promo' (should return true)" -ForegroundColor Yellow
$result = & .\validator.exe promo data\campaign_codes.txt data\membership_codes.txt
Write-Host "Result: $result" -ForegroundColor Green
Write-Host "Exit Code: $LASTEXITCODE" -ForegroundColor Gray
Write-Host ""

Write-Host "Test 2: Code 'xyzki' (should return false)" -ForegroundColor Yellow  
$result = & .\validator.exe xyzki data\campaign_codes.txt data\membership_codes.txt
Write-Host "Result: $result" -ForegroundColor Green
Write-Host "Exit Code: $LASTEXITCODE" -ForegroundColor Gray
Write-Host ""

Write-Host "Test 3: Invalid code 'PROMO' (should error - uppercase)" -ForegroundColor Yellow
$result = & .\validator.exe PROMO data\campaign_codes.txt data\membership_codes.txt 2>&1
Write-Host "Result: $result" -ForegroundColor Green
Write-Host "Exit Code: $LASTEXITCODE" -ForegroundColor Gray
Write-Host ""

Write-Host "=== Tests Complete ===" -ForegroundColor Cyan
