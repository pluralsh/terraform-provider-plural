terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.25"
    }
  }
}

provider "plural" {
  use_cli = true
}

resource "plural_scm_webhook" "test" {
  type = "GITHUB"
  owner = "pluralsh"
  hmac = "test"
}

resource "plural_scm_webhook" "duplicate" {
  type = "GITHUB"
  owner = "pluralsh"
  hmac = "test"
}
