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

resource "random_string" "random" {
  length = 5
  upper = false
  special = false
}

# TODO: Test deletion.
resource "plural_infrastructure_stack" "stack" {
  name = "stack-tf-${random_string.random.result}"
  type = "TERRAFORM"
  approval = true
  cluster_id = data.plural_cluster.cluster.id
  repository = {
    id = data.plural_git_repository.repository.id
    ref = "main"
    folder = "terraform"
  }
  configuration = {
    image = "hashicorp/terraform:1.8.1"
    version = "1.8.1"
  }
  files = {
    # TODO: Test it.
    # "test.yml": "value: 123"
  }
  environment = [
    # TODO: Test it.
    # {
    #   name = "USERNAME"
    #   value = "joe"
    # },
    # {
    #   name = "PASSWORD"
    #   value = "test"
    #   secret = true
    # }
  ]
  job_spec = {
    namespace = "default"
    labels = {
      test = "123"
    }
    service_account = "default"
    containers = [{
      image = "perl:5.34.0"
      args = ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      # TODO: Test without env and env_from.
      env = {}
      env_from = []
    }]
#     raw = jsonencode({
#       containers = [{
#         name  = "pi"
#         image = "perl:5.34.0"
#         command = ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
#       }]
#       restartPolicy = "Never"
#     })
  }
  bindings = {
     read = []
     write = []
  }
}
