provider "plural" {
  use_cli = true
}

data "plural_cluster" "byok_workload_cluster" {
  handle = "floreks-tf-workload-cluster"
}

data "plural_git_repository" "cd-test" {
  url = "https://github.com/pluralsh/plrl-cd-test.git"
}

resource "plural_service_deployment" "cd-test" {
  # Required
  name      = "tf-cd-test"
  namespace = "tf-cd-test"

  cluster = {
    handle = data.plural_cluster.byok_workload_cluster.handle
  }

  repository = {
    id     = data.plural_git_repository.cd-test.id
    ref    = "main"
    folder = "kubernetes"
  }

  # Optional
  version = "0.0.2"
  docs_path = "doc"
  protect   = false

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

  sync_config = {
    namespace_metadata = {
      annotations = {
        "testannotationkey" : "testannotationvalue"
      }
      labels = {
        "testlabelkey" : "testlabelvalue"
      }
    }
  }

  depends_on = [
    data.plural_cluster.byok_workload_cluster,
    data.plural_git_repository.cd-test
  ]
}
