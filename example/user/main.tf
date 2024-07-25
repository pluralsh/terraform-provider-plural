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

data "plural_user" "user" {
  email = "marcin@plural.sh"
}

data "plural_group" "group" {
  name = "team"
}