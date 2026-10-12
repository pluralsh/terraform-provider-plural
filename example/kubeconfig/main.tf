terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.28"
    }
  }
}

# Kubeconfig is configured on the provider. Use provider aliases for multiple clusters.
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