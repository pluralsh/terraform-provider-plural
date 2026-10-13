terraform {
  required_providers {
    plural = {
      source = "pluralsh/plural"
      # Version installed locally with `make install` (latest GitHub tag).
      version = "0.2.40"
    }
  }
}

provider "plural" {
  use_cli = true
}

# Cluster using the deprecated, resource-level kubeconfig.
# The deployment agent is installed into the local kind cluster.
resource "plural_cluster" "kind" {
  name   = "kubeconfig-deprecated-test"
  handle = "kubeconfig-deprecated-test"
  detach = true
  tags = {
    test = "kubeconfig-deprecated"
  }

  kubeconfig = {
    config_path    = pathexpand("~/.kube/config")
    config_context = "kind-kubeconfig-deprecated"
  }
}
