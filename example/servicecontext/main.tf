terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.27"
    }
  }
}

provider "plural" {
  use_cli = true
}

resource "plural_service_context" "service_context" {
  name           = "service-context-test"
  configuration = jsonencode({
    "env" = "prod"
    "test" = "some-value"
    "array" = [1, 2, 3]
    "nested_field" = {
      "test" = "nested-value"
    }
  })
  secrets = {
    "test" = "some-secret-value"
  }
}

data "plural_service_context" "service_context" {
  name = plural_service_context.service_context.name
}
