# dev.ps1 — start all services in separate PowerShell windows (Windows equivalent of dev.sh)
# Usage: .\dev.ps1
# Close each terminal window manually when done.

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition

Write-Host "[postgres] Starting Docker Compose (Postgres)..." -ForegroundColor Blue
Start-Process powershell -ArgumentList "-NoExit", "-ExecutionPolicy", "Bypass", "-Command", "Set-Location '$ScriptDir\backend'; make docker-run" `
    -WindowStyle Normal

Write-Host "[backend]  Waiting 5s for Postgres, then starting Go backend on :8080..." -ForegroundColor Green
Start-Process powershell -ArgumentList "-NoExit", "-ExecutionPolicy", "Bypass", "-Command", "Set-Location '$ScriptDir\backend'; Start-Sleep 5; make watch" `
    -WindowStyle Normal

Write-Host "[web]      Starting Next.js dev server on :3000..." -ForegroundColor Yellow
Start-Process powershell -ArgumentList "-NoExit", "-ExecutionPolicy", "Bypass", "-Command", "Set-Location '$ScriptDir\web'; pnpm dev" `
    -WindowStyle Normal

Write-Host ""
Write-Host "All three services are starting in separate windows:" -ForegroundColor Cyan
Write-Host "  - Postgres     (blue window)"   -ForegroundColor Blue
Write-Host "  - Go backend   (green window)  → http://localhost:8080" -ForegroundColor Green
Write-Host "  - Next.js web  (yellow window) → http://localhost:3000" -ForegroundColor Yellow
Write-Host ""
Write-Host "Close each window individually when you are done." -ForegroundColor Cyan
