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
  cloud = "aws"
  protect = "false"
  cloud_settings = {
    aws = {
      region = "us-east-1"
    }
  }
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
}

data "plural_cluster" "byok_workload_cluster" {
  handle = "wctf"
}

