terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.1"
    }
  }
}

provider "plural" {
  use_cli = true
}

resource "plural_cloud_connection" "aws" {
  name     = "my-aws-connection"
  provider = "AWS"
}

resource "plural_cloud_connection" "gcp" {
  name     = "my-gcp-connection"
  provider = "GCP"
}

resource "plural_cloud_connection" "azure" {
  name     = "my-azure-connection"
  provider = "AZURE"
}

data "plural_cloud_connection" "existing" {
  name = "existing-connection"
  # Or by ID if you know it
  # id = "existing-id"
}

output "aws_read_bindings" {
  description = "The read bindings for the AWS connection"
  value       = plural_cloud_connection.aws.read_bindings
}

output "connection_ids" {
  description = "The IDs of all created connections"
  value = {
    aws   = plural_cloud_connection.aws.id
    gcp   = plural_cloud_connection.gcp.id
    azure = plural_cloud_connection.azure.id
  }
}
