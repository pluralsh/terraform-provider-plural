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

resource "random_string" "random" {
  length  = 5
  upper   = false
  special = false
}

resource "plural_infrastructure_stack" "stack" {
  name       = "stack-tf-${random_string.random.result}"
  type       = "TERRAFORM"
  detach     = true
  cluster_id = data.plural_cluster.cluster.id
  repository = {
    id     = data.plural_git_repository.repository.id
    ref    = "main"
    folder = "terraform"
  }
  configuration = {
    version = "1.8.1"
  }
}

resource "plural_custom_stack_run" "run-full" {
  name       = "run-tf-full-${random_string.random.result}"
  documentation       = "test"
  stack_id = plural_infrastructure_stack.stack.id
  commands = [{
    cmd = "ls"
    args = ["-al"]
    dir = "/"
  }]
  configuration = [{
    type = "STRING"
    name = "author"
    default = "john"
    documentation = "author name"
    longform = "author name"
    placeholder = "author name, i.e. john"
    optional = true
    condition = {
      operation = "PREFIX"
      field = "author"
      value = "j"
    }
  }]
}

resource "plural_custom_stack_run" "run-minimal" {
  name       = "run-tf-minimal-${random_string.random.result}"
}
