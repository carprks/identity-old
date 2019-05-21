provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region = "${var.aws_region}"
  version = "~> 1.35"
}

locals {
  aws_vpc_stack_name = "${var.aws_resource_prefix}-vpc-stack"
  aws_ecs_service_stack_name = "${var.aws_service_name}-svc-stack"
  aws_ecr_repository_name = "${var.aws_resource_prefix}"
  aws_ecs_cluster_name = "${var.aws_resource_prefix}-cluster"
  aws_ecs_service_name = "${var.aws_resource_prefix}-${var.aws_service_name}"
  aws_ecs_execution_role_name = "${var.aws_resource_prefix}-ecs-execution-role"
}

resource "aws_ecr_repository" "identity-repository" {
  name = "${local.aws_ecr_repository_name}"
}

resource "aws_cloudformation_stack" "vpc" {
  name = "${local.aws_vpc_stack_name}"
  template_body = "${file("cloudformation-templates/public-vpc.yaml")}"
  capabilities = ["CAPABILITY_NAMED_IAM"]
  parameters {
    ClusterName = "${local.aws_ecs_cluster_name}"
    ExecutionRoleName = "${local.aws_ecs_execution_role_name}"
  }
}

resource "aws_cloudformation_stack" "ecs_service" {
  name = "${local.aws_ecs_service_stack_name}"
  template_body = "${file("cloudformation-templates/public-service.yaml")}"
  depends_on = [
    "aws_cloudformation_stack.vpc",
    "aws_ecr_repository.identity-repository"
  ]
  parameters {
    ContainerMemory = 512
    ContainerPort = 80
    StackName = "${local.aws_vpc_stack_name}"
    ServiceName = "${local.aws_ecs_service_name}"
  }
}