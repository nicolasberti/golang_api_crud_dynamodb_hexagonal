#!/bin/bash

set -e

echo "🚀 Deploying Product CRUD Hexagonal..."

# Check if Terraform is initialized
if [ ! -d "terraform/.terraform" ]; then
    echo "📦 Initializing Terraform..."
    cd terraform
    terraform init
    cd ..
fi

# Apply Terraform
echo "🏗️  Applying Terraform infrastructure..."
cd terraform
terraform apply -auto-approve
cd ..

# Get outputs
TABLE_NAME=$(cd terraform && terraform output -raw dynamodb_table_name)
echo "📊 DynamoDB table created: $TABLE_NAME"

# Update .env with table name
if [ -f .env ]; then
    sed -i.bak "s/DYNAMODB_TABLE=.*/DYNAMODB_TABLE=$TABLE_NAME/" .env
    echo "✅ Updated .env with table name"
else
    echo "⚠️  .env file not found. Please copy .env.example to .env"
fi

echo "🎉 Deployment completed!"
echo "📋 Next steps:"
echo "   1. Configure your AWS credentials"
echo "   2. Run: go run cmd/api/main.go"
echo "   3. API will be available at: http://localhost:8080"