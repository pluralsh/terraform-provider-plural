terraform {
  required_providers {
    plural = {
      source = "pluralsh/plural"
      version = "0.2.28"
    }
  }
}

provider "plural" {
  use_cli = true
}

resource "plural_cluster" "test" {
  name = "test-cluster"
  handle = "test"
  protect = false
  detach = true
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
}
