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

resource "plural_cluster" "azure_cluster" {
  name = "azure-cluster-tf"
  handle = "aztf"
  version = "1.26.3"
  provider_id = plural_provider.azure_provider.id
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
  node_pools = []
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
  bindings = {
    read = []
    write = []
  }
}
