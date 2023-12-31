provider "plural" {
  use_cli = true
}

data "plural_cluster" "byok_workload_cluster" {
  handle = "gcp-workload-cluster"
}

resource "plural_service_deployment" "cd-test" {
  # Required
  name      = "tf-cd-helm-test"
  namespace = "tf-cd-helm-test"

  cluster = {
    handle = data.plural_cluster.byok_workload_cluster.handle
  }

  # Requires flux-source-controller addon to be installed and flux repo CRD for podinfo to exist
  helm = {
    chart = "podinfo"
    repository = {
      name = "podinfo"
      namespace = "default"
    }
    version = "6.5.3"
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
  ]
}
