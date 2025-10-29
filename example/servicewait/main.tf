terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.28"
    }
  }
}

provider "plural" {
  use_cli = true
}

resource "plural_service_wait" "test" {
  cluster = "mgmt"
  service = "console"
  warmup     = "10s"
  duration   = "1m"
}

