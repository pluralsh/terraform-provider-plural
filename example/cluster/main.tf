terraform {
  required_providers {
    plural = {
      source = "pluralsh/plural"
      version = "0.2.29"
    }
  }
}

provider "plural" {
  use_cli = true
  kubeconfig = {
    # It can be sourced from environment variables instead, i.e.: export PLURAL_KUBE_CONFIG_PATH=$KUBECONFIG
    config_path = pathexpand("~/.kube/config")
  }
}

data "plural_project" "test" {
  name = "test"
}

resource "plural_cluster" "test" {
  name = "test-cluster"
  handle = "test"
  protect = false
  detach = true
  project_id = data.plural_project.test.id
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
}
