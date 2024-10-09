terraform {
  required_providers {
    plural = {
      source = "pluralsh/plural"
      version = "0.2.1"
    }
  }
}

provider "plural" {
  use_cli = true
}

resource "plural_oidc_provider" "provider" {
  name = "marcin"
  type = "PLURAL"
  description = "test provider"
  redirect_uris = ["abc", "xyz"]
}
