provider "plural" {
  use_cli = true
}

resource "plural_cluster" "byok_workload_cluster" {
  name = "workload-cluster"
  handle = "workload-cluster"
  cloud = "byok"
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
}
