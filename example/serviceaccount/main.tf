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

resource "plural_service_account" "bot" {
  name = "Automation Bot"
  email = "bot@plural.sh"
}
