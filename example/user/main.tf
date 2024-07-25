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

data "plural_config" "config" {}

data "plural_user" "user" {
  email = "marcin@plural.sh"
}

data "plural_group" "group" {
  name = "team"
}

# resource "plural_group" "test" {
#   name = "test"
#   description = "test group"
# }
#
# resource "plural_group" "empty" {
#   name = "empty"
# }
