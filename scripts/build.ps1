# VibeCast Production Build & Test Script
Write-Host "Building VibeCast Enterprise Suite..." -ForegroundColor Cyan

Set-Location -Path $PSScriptRoot\..

Write-Host "Building Core SFU Server (cmd/sfu-server)..." -ForegroundColor Yellow
go build -o bin/sfu-server.exe ./cmd/sfu-server

Write-Host "Building Agent Worker Daemon (cmd/agent-worker)..." -ForegroundColor Yellow
go build -o bin/agent-worker.exe ./cmd/agent-worker

Write-Host "Building Egress Recording Server (cmd/egress-server)..." -ForegroundColor Yellow
go build -o bin/egress-server.exe ./cmd/egress-server

Write-Host "Building Ingress Ingest Server (cmd/ingress-server)..." -ForegroundColor Yellow
go build -o bin/ingress-server.exe ./cmd/ingress-server

Write-Host "Building SIP Telephony Gateway (cmd/sip-server)..." -ForegroundColor Yellow
go build -o bin/sip-server.exe ./cmd/sip-server

Write-Host "VibeCast Enterprise All-Stage Build Completed Successfully!" -ForegroundColor Green



