terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.27"
    }
  }
}

provider "plural" {
  use_cli = true
}

data "plural_cluster" "cluster" {
  handle = "mgmt"
}

data "plural_git_repository" "repository" {
  url = "https://github.com/zreigz/tf-hello.git"
}

resource "plural_git_repository" "repository" {
  url = "https://github.com/zreigz/tf-hello.git"
}
