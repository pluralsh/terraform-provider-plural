terraform {
  required_providers {
    plural = {
      source = "pluralsh/plural"
      version = "0.2.27"
    }
  }
}

provider "plural" {
  use_cli = true
}

#resource "plural_provider" "azure_provider" {
#  name = "azure"
#  cloud = "azure"
#  cloud_settings = {
#    azure = {
#      # subscription_id = "" # Required, can be sourced from PLURAL_AZURE_SUBSCRIPTION_ID
#      # tenant_id = "" # Required, can be sourced from PLURAL_AZURE_TENANT_ID
#      # client_id = "" # Required, can be sourced from PLURAL_AZURE_CLIENT_ID
#      # client_secret = "" # Required, can be sourced from PLURAL_AZURE_CLIENT_SECRET
#    }
#  }
#}

data "plural_provider" "azure_provider" {
  cloud = "azure"
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
      subscription_id = "" # Required
      location = "eastus"
    }
  }
  metadata = jsonencode({
    test1 = "test"
    test2 = false
    test3 = jsonencode({
      abc = false
    })
  })
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
}
