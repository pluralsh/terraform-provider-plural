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

resource "plural_provider" "azure_provider" {
  name = "azure"
  cloud = "azure"
  cloud_settings = {
    azure = {
      subscription_id = "" # Provide before use
      tenant_id = "" # Provide before use
      client_id = "" # Provide before use
      client_secret = "" # Provide before use
    }
  }
}

data "plural_provider" "azure_provider" {
  cloud = "aws"
}

resource "plural_cluster" "azure_cluster" {
  name = "azure-cluster-tf"
  handle = "aztf"
  version = "1.25.11"
  provider_id = data.plural_provider.azure_provider.id
  cloud = "azure"
  protect = "false"
  cloud_settings = {
    azure = {
      resource_group = "azure-cluster-tf"
      network = "azure-cluster-tf"
      subscription_id = "" # Provide before use
      location = "eastus"
    }
  }
#  node_pools = []
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
  bindings = {
    read = []
    write = []
  }
}
