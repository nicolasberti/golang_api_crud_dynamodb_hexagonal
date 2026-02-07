output "dynamodb_table_name" {
  description = "DynamoDB table name"
  value       = aws_dynamodb_table.products.name
}

output "dynamodb_table_arn" {
  description = "DynamoDB table ARN"
  value       = aws_dynamodb_table.products.arn
}

output "iam_role_arn" {
  description = "IAM role ARN for Lambda"
  value       = aws_iam_role.lambda_role.arn
}