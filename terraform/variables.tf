variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "product-crud"
}

variable "table_name" {
  description = "DynamoDB table name"
  type        = string
  default     = "products"
}