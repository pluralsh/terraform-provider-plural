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

data "plural_cluster" "cluster" {
  handle = "mgmt"
}

data "plural_git_repository" "repository" {
  url = "https://github.com/zreigz/tf-hello.git"
}

resource "plural_infrastructure_stack" "stack" {
  name = "tf-stack-2"
  type = "TERRAFORM"
  # approval = false
  cluster_id = data.plural_cluster.cluster.id
  repository = {
    id = data.plural_git_repository.repository.id
    ref = "main"
    folder = "terraform"
  }
  configuration = {
    # image = ""
    version = "1.5.7"
  }
  # files = {}
  environment = [
    {
      name = "USERNAME"
      value = "joe"
    },
    {
      name = "PASSWORD"
      value = "test"
      secret = true
    }
  ]
  job_spec = {
    namespace = "default"
    raw = ""
    # ...
  }
  bindings = {
    read = []
    write = []
  }
}
