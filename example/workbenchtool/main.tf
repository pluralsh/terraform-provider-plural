terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.40"
    }
  }
}

provider "plural" {
  use_cli = true
}

variable "cloud_connection_name" {
  type        = string
  description = "Name of an existing AWS cloud connection used to invoke the Lambda functions."
}

###############################################################################
# Cloud connection lookup
###############################################################################

# Look up an existing cloud connection by name. Its credentials are never
# returned; only id, name and cloud_provider are available.
data "plural_cloud_connection" "by_name" {
  name = var.cloud_connection_name
}

# The same connection looked up by id.
data "plural_cloud_connection" "by_id" {
  id = data.plural_cloud_connection.by_name.id
}

###############################################################################
# Workbench tools
###############################################################################

locals {
  input_schema = jsonencode({
    type = "object"
    properties = {
      action   = { type = "string", enum = ["plan", "execute"] }
      volumeId = { type = "string" }
    }
    required = ["action", "volumeId"]
  })
}

# Deletes resources, so every invocation waits for human approval.
resource "plural_workbench_tool" "delete_volume" {
  name                = "tf_delete_volume"
  tool                = "LAMBDA"
  approval            = true
  cloud_connection_id = data.plural_cloud_connection.by_name.id

  configuration = {
    lambda = {
      lambda_arn   = "arn:aws:lambda:us-east-1:123456789012:function:delete-volume"
      description  = "Delete an unattached EBS volume."
      input_schema = local.input_schema
    }
  }
}

# Read-only, so approval is disabled explicitly.
resource "plural_workbench_tool" "describe_volume" {
  name                = "tf_describe_volume"
  tool                = "LAMBDA"
  approval            = false
  cloud_connection_id = data.plural_cloud_connection.by_id.id

  configuration = {
    lambda = {
      lambda_arn   = "arn:aws:lambda:us-east-1:123456789012:function:describe-volume"
      description  = "Describe an EBS volume."
      input_schema = local.input_schema
    }
  }
}

# approval is not set: the tool is created without approval, and later applies
# keep whatever approval is currently set on the tool, e.g. from the Console UI.
resource "plural_workbench_tool" "echo" {
  name                = "tf_echo"
  tool                = "LAMBDA"
  cloud_connection_id = data.plural_cloud_connection.by_name.id

  configuration = {
    lambda = {
      lambda_arn   = "arn:aws:lambda:us-east-1:123456789012:function:echo"
      description  = "Echo the input back."
      input_schema = local.input_schema
    }
  }
}

###############################################################################
# Outputs
###############################################################################

output "cloud_connection" {
  description = "The looked up cloud connection."
  value = {
    id             = data.plural_cloud_connection.by_name.id
    name           = data.plural_cloud_connection.by_id.name
    cloud_provider = data.plural_cloud_connection.by_name.cloud_provider
  }
}

output "approval" {
  description = "Approval setting of each tool."
  value = {
    delete_volume   = plural_workbench_tool.delete_volume.approval
    describe_volume = plural_workbench_tool.describe_volume.approval
    echo            = plural_workbench_tool.echo.approval
  }
}
