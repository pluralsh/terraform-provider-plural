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

resource "plural_cluster" "byok_workload_cluster" {
  name = "workload-cluster-tf"
  handle = "wctf"
  cloud = "byok"
  protect = "false"
  cloud_settings = {
    byok = {
      # Provide before use
    }
  }
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
}

data "plural_cluster" "byok_workload_cluster" {
  handle = "wctf"
}
