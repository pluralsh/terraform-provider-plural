terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.25"
    }
  }
}

provider "plural" {
  use_cli = true
}

###############################################################################
# (Optional) Input variables – move to variables.tf if you prefer
###############################################################################

variable "aws_access_key_id" {
  type      = string
  sensitive = true
}

variable "aws_secret_access_key" {
  type      = string
  sensitive = true
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

data "plural_user" "john" {
  email = "john@plural.sh"
}

data "plural_user" "john_doe" {
  email = "marcin@plural.sh"
}

###############################################################################
# Resources
###############################################################################

# Group that will get read access to the cloud connection

data "plural_group" "existing-cloud_admins" {
  name        = "existing-cloud-admins"
}

# Cloud-connection resource (AWS example)
resource "plural_cloud_connection" "aws" {
  name           = "your-connection-name"
  cloud_provider = "AWS"

  configuration = {
    aws = {
      access_key_id     = var.aws_access_key_id
      secret_access_key = var.aws_secret_access_key
      region            = var.aws_region
    }
  }

  read_bindings = [
    {
      user_id = data.plural_user.john_doe.id
    },
    {
      group_id = data.plural_group.existing-cloud_admins.id
    }
  ]
}

###############################################################################
# Outputs
###############################################################################

output "cloud_connection_id" {
  value       = plural_cloud_connection.aws.id
  description = "ID of the created cloud connection"
}

output "cloud_connection_details" {
  value       = plural_cloud_connection.aws
  description = "All attributes of the cloud connection"
  sensitive   = true
}

output "group_id" {
  value       = data.plural_group.existing-cloud_admins.id
  description = "ID of the created group"
}
