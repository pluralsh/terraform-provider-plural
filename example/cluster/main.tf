terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.36"
    }
  }
}

provider "plural" {
  use_cli = true
  kubeconfig = {
    # Can be sourced from environment variables instead:
    # export PLURAL_KUBE_CONFIG_PATH=$KUBECONFIG
    config_path = pathexpand("~/.kube/config")
  }
}

data "plural_project" "default" {
  name = "default"
}

resource "plural_cluster" "test" {
  name       = "test-cluster"
  handle     = "test"
  protect    = true
  detach     = true
  project_id = data.plural_project.default.id
  metadata   = jsonencode({ "test-key" = "metadata-update" })

  tags = {
    "managed-by"    = "terraform-provider-plural"
    "update-marker" = "updated-2"
  }
}

output "cluster_id" {
  value = plural_cluster.test.id
}

output "agent_deployed" {
  value = plural_cluster.test.agent_deployed
}
