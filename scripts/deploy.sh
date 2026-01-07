#!/bin/bash
set -e

echo "Deploying Connectra API to AWS Lambda..."

# Check prerequisites
if ! command -v sam &> /dev/null; then
    echo "Error: SAM CLI is not installed"
    echo "Install it from: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html"
    exit 1
fi

if ! command -v aws &> /dev/null; then
    echo "Error: AWS CLI is not installed"
    echo "Install it from: https://aws.amazon.com/cli/"
    exit 1
fi

# Build Lambda function
echo "Building Lambda function..."
./scripts/build.sh

# Build with SAM
echo "Building with SAM..."
sam build

# Deploy
echo "Deploying to AWS..."
sam deploy --guided

echo "Deployment complete!"
