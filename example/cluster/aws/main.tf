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
      # access_key_id = "" # Required, can be sourced from PLURAL_AWS_ACCESS_KEY_ID
      # secret_access_key = "" # Required, can be sourced from PLURAL_AWS_SECRET_ACCESS_KEY
    }
  }
}

data "plural_provider" "aws_provider" {
  cloud = "aws"
}

resource "plural_cluster" "aws_cluster" {
  name = "aws-cluster-tf"
  handle = "awstf"
  version = "1.24"
  provider_id = data.plural_provider.aws_provider.id
  cloud = "aws"
  protect = "false"
  cloud_settings = {
    aws = {
      region = "us-east-1"
    }
  }
  node_pools = {
    pool1 = {
      name = "pool1"
      min_size = 1
      max_size = 5
      instance_type = "t5.large"
    },
    pool2 = {
      name = "pool2"
      min_size = 1
      max_size = 5
      instance_type = "t5.large"
      labels = {
        "key1" = "value1"
        "key2" = "value2"
      },
      taints = [
        {
          key = "test"
          value = "test"
          effect = "NoSchedule"
        }
      ]
    },
    pool3 = {
      name = "pool3"
      min_size = 1
      max_size = 5
      instance_type = "t5.large"
      cloud_settings = {
        aws = {
          launch_template_id = "test"
        }
      }
    }
  }
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
}
