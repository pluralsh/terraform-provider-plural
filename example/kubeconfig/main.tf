terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.25"
    }
  }
}

#####################################
######### Deprecated method #########
#####################################
# provider "plural" {
#   alias   = "deprecated"
#   use_cli = true
# }
#
# resource "plural_cluster" "deprecated" {
#   provider = plural.deprecated
#   name     = "byok"
#   protect  = "false"
#   detach   = true
#   kubeconfig = {
#     config_path = pathexpand("~/.kube/config") # This can no longer be sourced from environment variables.
#   }
# }

#####################################
############ New method #############
#####################################
provider "plural" {
  alias   = "new"
  use_cli = true
  kubeconfig = {
    # Can be sourced from environment variables, export PLURAL_KUBE_CONFIG_PATH to read from local file:
    # export PLURAL_KUBE_CONFIG_PATH=$KUBECONFIG
  }
}

resource "plural_cluster" "new" {
  provider = plural.new
  name     = "byok"
  protect  = "false"
  detach   = true

  metadata = jsonencode({
    test1 = "test"
    test2 = false
    test3 = jsonencode({
      abc = false
    })
  })
}