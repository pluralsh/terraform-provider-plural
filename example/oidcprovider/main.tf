terraform {
  required_providers {
    plural = {
      source = "pluralsh/plural"
      version = "0.2.27"
    }
  }
}

provider "plural" {
  use_cli = true
}

data "plural_user" "user" {
  email = "marcin@plural.sh"
}

resource "plural_oidc_provider" "provider" {
  name = "tf-test-provider"
  auth_method = "BASIC"
  type = "CONSOLE"
  description = "test provider"
  redirect_uris = ["localhost:8000"]
  bindings = [
    { user_id = data.plural_user.user.id }
  ]
}
