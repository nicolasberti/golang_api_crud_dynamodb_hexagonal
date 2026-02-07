@echo off
setlocal

echo 🚀 Deploying Product CRUD Hexagonal...

REM Check if Terraform is initialized
if not exist "terraform\.terraform" (
    echo 📦 Initializing Terraform...
    cd terraform
    terraform init
    cd ..
)

REM Apply Terraform
echo 🏗️  Applying Terraform infrastructure...
cd terraform
terraform apply -auto-approve
cd ..

REM Get outputs
for /f "delims=" %%i in ('cd terraform && terraform output -raw dynamodb_table_name') do set TABLE_NAME=%%i
echo 📊 DynamoDB table created: %TABLE_NAME%

REM Update .env with table name
if exist .env (
    powershell -Command "(Get-Content .env) -replace 'DYNAMODB_TABLE=.*', 'DYNAMODB_TABLE=%TABLE_NAME%' | Set-Content .env"
    echo ✅ Updated .env with table name
) else (
    echo ⚠️  .env file not found. Please copy .env.example to .env
)

echo 🎉 Deployment completed!
echo 📋 Next steps:
echo    1. Configure your AWS credentials
echo    2. Run: go run cmd/api/main.go
echo    3. API will be available at: http://localhost:8080

pause