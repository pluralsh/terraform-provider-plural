provider "plural" {
  use_cli = true
}

resource "plural_provider" "gcp_provider" {
  name           = "gcp"
  cloud          = "gcp"
  cloud_settings = {
    gcp = {
      # credentials = "" # Required, can be sourced from PLURAL_GCP_CREDENTIALS
    }
  }
}

resource "plural_cluster" "gcp_workload_cluster" {
  name           = "gcp-workload-cluster"
  handle         = "gcp-workload-cluster"
  cloud          = "gcp"
  provider_id    = plural_provider.gcp_provider.id
  version        = "1.25.11"
  cloud_settings = {
    gcp = {
      region = "" # Required
      network = "" # Required
      project = "" # Required
    }
  }
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
  depends_on = [plural_provider.gcp_provider]
}
