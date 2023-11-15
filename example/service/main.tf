provider "plural" {
  use_cli = true
}

resource "plural_cluster" "byok_workload_cluster" {
  name   = "tf-workload-cluster"
  handle = "tf-workload-cluster"
  cloud  = "byok"
  tags   = {
    "managed-by" = "terraform-provider-plural"
  }
}

data "plural_git_repository" "cd-test" {
  url = "https://github.com/pluralsh/plrl-cd-test.git"
}

resource "plural_service_deployment" "cd-test" {
  name          = "tf-cd-test"
  namespace     = "tf-cd-test"
  configuration = [
    {
      name : "host"
      value : "tf-cd-test.gcp.plural.sh"
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
    id     = data.plural_git_repository.cd-test.id
    ref    = "main"
    folder = "kubernetes"
  }

  depends_on = [
    plural_cluster.byok_workload_cluster,
    data.plural_git_repository.cd-test
  ]
}
