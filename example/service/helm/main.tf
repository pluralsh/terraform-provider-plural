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

data "plural_cluster" "cluster" {
  handle = "mgmt"
}

resource "plural_service_deployment" "cd-test" {
  # Required
  name      = "tf-cd-helm-test"
  namespace = "tf-cd-helm-test"

  cluster = {
    handle = data.plural_cluster.cluster.handle
  }

  helm = {
    chart = "grafana"
    version = "8.x.x"
    url = "https://grafana.github.io/helm-charts"
  }

  # Optional
  version = "0.0.2"
  docs_path = "doc"
  protect   = false

  configuration = {
    "host" = "tf-cd-test.gcp.plural.sh",
    "tag" = "sha-4d01e86"
  }

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
}
