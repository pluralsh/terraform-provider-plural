terraform {
  required_providers {
    plural = {
      source = "pluralsh/plural"
      version = "0.0.1"
    }
  }
}

provider "plural" {
  use_cli = true
}

resource "plural_cluster" "byok_workload_cluster" {
  name = "workload-cluster-tf"
  handle = "wctf"
  cloud = "byok"
  protect = "false"
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
}

resource "plural_git_repository" "cd-test" {
  url = "https://github.com/pluralsh/plrl-cd-test.git"
}

resource "plural_service_deployment" "cd-test" {
  name          = "cd-test"
  namespace     = "cd-test"
  configuration = [
    {
      name : "host"
      value : "cd-test.gcp.plural.sh"
    },
    {
      name : "tag"
      value : "sha-4d01e86"
    }
  ]

  cluster = {
    handle = plural_cluster.byok_workload_cluster.handle
  }

  repository = {
    id     = plural_git_repository.cd-test.id
    ref    = "main"
    folder = "kubernetes"
  }

  depends_on = [
    plural_git_repository.cd-test
  ]
}
