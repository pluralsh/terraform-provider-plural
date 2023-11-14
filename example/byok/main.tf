provider "plural" {
  console_url = ""
  access_token = ""
}

resource "plural_cluster" "byok_workload_cluster" {
  name = "workload-cluster"
  handle = "workload-cluster"
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
}
