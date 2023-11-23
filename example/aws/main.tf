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
      access_key_id = ""
      secret_access_key = ""
    }
  }
}

resource "plural_cluster" "aws_cluster" {
  name = "workload-cluster-tf"
  handle = "wctf"
  version = "1.23"
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

data "plural_cluster" "byok_workload_cluster" {
  handle = "wctf"
}

