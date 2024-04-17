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

resource "plural_cluster" "byok_workload_cluster" {
  name           = "workload-cluster-tf"
  handle         = "wctf"
  protect        = "false"
  kubeconfig = {
    # Required, can be sourced from environment variables
    # export PLURAL_KUBE_CONFIG_PATH to read from local file
  }
  metadata = jsonencode({
    test1 = "string"
    test2 = false
    test3 = jsonencode({
      abc = false
    })
  })
  # Alternative method
  #  metadata = <<EOF
  #{
  #  "test1": "string",
  #  "test2": false,
  #  "test3": {"abc": false},
  #}
  #EOF
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
}

#data "plural_cluster" "byok_workload_cluster" {
#  handle = "wctf"
#}
