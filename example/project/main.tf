terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.1"
    }
  }
}

provider "plural" {
  use_cli = true
}

data "plural_project" "default" {
  name       = "default"
  # id should work as well
}

data "plural_cluster" "cluster" {
  handle = "mgmt"
}



