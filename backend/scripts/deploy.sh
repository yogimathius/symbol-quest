#!/bin/bash

# Symbol Quest Backend Deployment Script

set -e

echo "🚀 Deploying Symbol Quest API to Fly.io..."

# Check if fly CLI is installed
if ! command -v flyctl &> /dev/null; then
    echo "❌ flyctl CLI not found. Please install it first:"
    echo "curl -L https://fly.io/install.sh | sh"
    exit 1
fi

# Check if already logged in
if ! flyctl auth whoami &> /dev/null; then
    echo "❌ Please login to Fly.io first: flyctl auth login"
    exit 1
fi

# Build and test locally first
echo "🔨 Building application..."
go build -o symbol-quest cmd/main.go

echo "🧪 Running quick test..."
go test ./...

# Deploy to Fly.io
echo "🚁 Deploying to Fly.io..."

# Create app if it doesn't exist
if ! flyctl apps list | grep -q "symbol-quest-api"; then
    echo "📦 Creating new Fly.io app..."
    flyctl launch --name symbol-quest-api --region iad --no-deploy
fi

# Create or attach PostgreSQL database
if ! flyctl postgres list | grep -q "symbol-quest-db"; then
    echo "🗄️ Creating PostgreSQL database..."
    flyctl postgres create --name symbol-quest-db --region iad
    flyctl postgres attach --postgres-app symbol-quest-db
fi

# Set secrets (you'll need to add these manually)
echo "🔐 Setting up secrets..."
echo "Please set the following secrets using flyctl secrets set:"
echo "flyctl secrets set JWT_SECRET=\"your-secure-jwt-secret\""
echo "flyctl secrets set OPENAI_API_KEY=\"your-openai-api-key\""
echo "flyctl secrets set STRIPE_SECRET_KEY=\"your-stripe-secret-key\""
echo "flyctl secrets set STRIPE_WEBHOOK_SECRET=\"your-stripe-webhook-secret\""

# Deploy
flyctl deploy

echo "✅ Deployment complete!"
echo "🌐 Your API is available at: https://symbol-quest-api.fly.dev"
echo "🏥 Health check: https://symbol-quest-api.fly.dev/health"