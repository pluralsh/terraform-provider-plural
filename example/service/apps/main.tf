terraform {
  required_providers {
    plural = {
      source = "pluralsh/plural"
      version = "0.2.24"
    }
  }
}

provider "plural" {
  use_cli = true
}

data "plural_cluster" "mgmt" {
  handle = "mgmt"
}

data "plural_git_repository" "repository" {
  url = "https://github.com/pluralsh/plrl-cd-test.git"
}

resource "random_string" "random" {
  length  = 5
  upper   = false
  special = false
}

resource "plural_service_deployment" "apps" {
  name = "test-${random_string.random.result}"
  namespace = "test"
  repository = {
    id = data.plural_git_repository.repository.id
    ref    = "main"
    folder = "kubernetes"
  }
  cluster = {
    id = data.plural_cluster.mgmt.id
  }
  templated = false
  # protect = true
  # configuration = {
  #   "host" = "tf-cd-test.gcp.plural.sh"
  #   "tag" = "sha-4d01e86"
  # }
}