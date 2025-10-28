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
  service_id = "b775ef81-b469-4c8c-969d-1d35e97a4ce5"
  warmup     = "30s"
  duration   = "1m"
}

