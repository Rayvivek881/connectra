# PowerShell deployment script for Lambda function
param(
    [string]$StackName = "connectra-api",
    [string]$Region = "us-east-1"
)

Write-Host "Deploying Connectra API to AWS Lambda..." -ForegroundColor Green

# Check prerequisites
if (-not (Get-Command sam -ErrorAction SilentlyContinue)) {
    Write-Host "Error: SAM CLI is not installed" -ForegroundColor Red
    Write-Host "Install it from: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html"
    exit 1
}

if (-not (Get-Command aws -ErrorAction SilentlyContinue)) {
    Write-Host "Error: AWS CLI is not installed" -ForegroundColor Red
    Write-Host "Install it from: https://aws.amazon.com/cli/"
    exit 1
}

# Build Lambda function
Write-Host "Building Lambda function..." -ForegroundColor Cyan
& .\scripts\build.ps1

# Build with SAM
Write-Host "Building with SAM..." -ForegroundColor Cyan
sam build

# Deploy
Write-Host "Deploying to AWS..." -ForegroundColor Cyan
sam deploy --stack-name $StackName --region $Region --guided

Write-Host "Deployment complete!" -ForegroundColor Green
