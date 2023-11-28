terraform {
  required_providers {
    plural = {
      source = "pluralsh/plural"
      version = "0.0.1"
    }
  }
}

provider "plural" {
  use_cli = true
}

resource "plural_provider" "aws_provider" {
  name = "aws"
  cloud = "aws"
  cloud_settings = {
    aws = {
      access_key_id = "" # Provide before use
      secret_access_key = "" # Provide before use
    }
  }
}

resource "plural_cluster" "aws_cluster" {
  name = "aws-cluster-tf"
  handle = "awstf"
  version = "1.24"
  provider_id = plural_provider.aws_provider.id
  cloud = "aws"
  protect = "false"
  cloud_settings = {
    aws = {
      region = "us-east-1"
    }
  }
  node_pools = []
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
  bindings = {
    read = []
    write = []
  }
}
