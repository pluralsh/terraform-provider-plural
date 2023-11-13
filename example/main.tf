provider "pluralcd" {
  console_url = ""
  access_token = ""
}

resource "pluralcd_aws_provider" "aws_provider" {
  access_key_id = ""
  secret_access_key = ""
}

resource "pluralcd_aws_cluster" "aws_workload_cluster" {
  name = "workload-cluster"
  handle = "workload-cluster"
  version = "1.23"
  region = "us-east-1"
  provider = pluralcd_aws_provider.aws_provider.id
  tags = {
    "managed-by": "terraform-provider-pluralcd"
  }
}
