# PowerShell build script for Lambda function
Write-Host "Building Lambda function for Connectra API..." -ForegroundColor Green

# Build for Linux/amd64 (Lambda runtime)
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o bootstrap -ldflags="-s -w" ./cmd/lambda/main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build complete! Binary: bootstrap" -ForegroundColor Green
    $size = (Get-Item bootstrap).Length / 1MB
    Write-Host "Binary size: $([math]::Round($size, 2)) MB" -ForegroundColor Cyan
} else {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}
