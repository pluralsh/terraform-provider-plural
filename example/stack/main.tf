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
  name = "tf-stack-13"
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
    # "test.yml": "value: 123"
  }
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
    raw = jsonencode({
      containers = jsonencode([{
        name  = "pi"
        image = "perl:5.34.0"
        command = jsonencode(["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"])
      }])
      restartPolicy = "Never"
    })
    # ...
  }
  bindings = {
     read = []
     write = []
  }
}
