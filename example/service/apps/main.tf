terraform {
  required_providers {
    plural = {
      source = "pluralsh/plural"
      version = "0.2.1"
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

resource "plural_service_deployment" "apps" {
  name = "apps-copy"
  namespace = "infra-copy"
  repository = {
    id = data.plural_git_repository.repository.id
    ref    = "main"
    folder = "kubernetes"
  }
  cluster = {
    id = data.plural_cluster.mgmt.id
  }

  protect = true
  templated = true
}