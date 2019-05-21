variable "aws_access_key" {}
variable "aws_secret_key" {}
variable "aws_account_id" {}
variable "aws_region" {
  description = "e.g. eu-west-2"
}
variable "aws_resource_prefix" {
  description = "e.g. carprks-dev"
}
variable "aws_service_name" {
  description = "e.g. identity"
}