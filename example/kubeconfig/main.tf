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
#     # It can no longer be sourced from environment variables.
#     config_path = pathexpand("~/.kube/config")
#   }
# }

#####################################
############ New method #############
#####################################
provider "plural" {
  alias   = "new"
  use_cli = true
  kubeconfig = {
    # It can be sourced from environment variables instead, i.e.: export PLURAL_KUBE_CONFIG_PATH=$KUBECONFIG
    config_path = pathexpand("~/.kube/config")
  }
}

resource "plural_cluster" "new" {
  provider = plural.new
  name     = "byok"
  protect  = "false"
  detach   = true
}