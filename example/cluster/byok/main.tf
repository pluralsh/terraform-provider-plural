terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.0.1"
    }
  }
}

provider "plural" {
  use_cli = true
}

resource "plural_cluster" "byok" {
  name           = "byok"
  protect        = "false"
  kubeconfig = {
    # Required, can be sourced from environment variables
    # export PLURAL_KUBE_CONFIG_PATH to read from local file
  }
  metadata = jsonencode({
    test1 = "test"
    test2 = false
    test3 = jsonencode({
      abc = false
    })
  })
  helm_repo_url = "https://pluralsh.github.io/deployment-operator"
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
}
