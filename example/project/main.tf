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

data "plural_project" "default" {
  name = "default"
}

data "plural_cluster" "mgmt" {
  handle = "mgmt"
}

data "plural_git_repository" "tf-hello" {
  url = "https://github.com/zreigz/tf-hello.git"
}

resource "random_string" "random" {
  length  = 5
  upper   = false
  special = false
}

# resource "plural_project" "test" {
#   name       = "test-${random_string.random.result}"
#   description = "test project created by terraform"
# }
#
# resource "plural_cluster" "byok" {
#   name       = "byok-${random_string.random.result}"
#   project_id = data.plural_project.default.id
#   kubeconfig = {
#     # Required, can be sourced from environment variables
#     # export PLURAL_KUBE_CONFIG_PATH to read from local file
#   }
# }

resource "plural_infrastructure_stack" "tf-hello" {
  name       = "tf-hello-${random_string.random.result}"
  type       = "TERRAFORM"
  cluster_id = data.plural_cluster.mgmt.id
  repository = {
    id     = data.plural_git_repository.tf-hello.id
    ref    = "main"
    folder = "terraform"
  }
  configuration = {
    image = "ghcr.io/pluralsh/harness"
    version = "sha-e9b2089-terraform-1.8"
  }
}