terraform {
  required_providers {
    plural = {
      source = "pluralsh/plural"
      version = "0.2.28"
    }
  }
}

provider "plural" {
  use_cli = true
}

data "plural_project" "test" {
  name = "test"
}

resource "plural_cluster" "test" {
  name = "test-cluster"
  handle = "test"
  protect = false
  detach = true
  project_id = data.plural_project.test.id
  tags = {
    "managed-by" = "terraform-provider-plural"
  }
}
