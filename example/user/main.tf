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

# data "plural_config" "config" {}
#
# data "plural_user" "user" {
#   email = "marcin@plural.sh"
# }

resource "plural_user" "spiderman" {
  name = "Peter Parker"
  email = "spiderman@plural.sh"
}

# data "plural_group" "avengers" {
#   name = "avengers"
#   global = "false"
# }
#
resource "plural_group" "avengers" {
  name = "avengers"
  description = "avengers group"
  global = "false"
}

resource "plural_group_member" "spiderman" {
  user_id = plural_user.spiderman.id
  group_id = plural_group.avengers.id
}

# resource "plural_group_member" "duplicate" {
#   user_id = plural_user.spiderman.id
#   group_id = plural_group.avengers.id
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