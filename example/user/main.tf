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

resource "plural_user" "user" {
  name = "Marcin Maciaszczyk"
  email = "marcin@plural.sh"
}

# data "plural_group" "group" {
#   name = "team"
# }
#
# resource "plural_group" "test" {
#   name = "test"
#   description = "test group"
# }
#
# resource "plural_rbac" "rbac" {
#   service_id = "624bff88-05e3-45f6-bc3b-44708594e28e"
#   bindings = {
#     read  = [{
#       user_id = data.plural_user.user.id
#     }]
#     write = [{
#       user_id = data.plural_user.user.id
#     }]
#   }
# }