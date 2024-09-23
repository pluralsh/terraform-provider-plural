terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.1"
    }
  }
}

provider "plural" {
  use_cli = true
}

resource "plural_service_context" "service_context" {
  name           = "service-context-tf"
  configuration = {
    "env" = "prod"
    "test" = "some-value"
  }
  secrets = {
    "test" = "some-secret-value"
  }
}

data "plural_service_context" "service_context" {
  name = "service-context-tf"
}
