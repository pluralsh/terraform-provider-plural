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
  length  = 5
  upper   = false
  special = false
}

resource "plural_infrastructure_stack" "stack-full" {
  name       = "stack-tf-full-${random_string.random.result}"
  type       = "TERRAFORM"
  approval   = true
  detach     = true
  cluster_id = data.plural_cluster.cluster.id
  repository = {
    id     = data.plural_git_repository.repository.id
    ref    = "main"
    folder = "terraform"
  }
  configuration = {
    image   = "hashicorp/terraform:1.8.1"
    version = "1.8.1"
  }
  files = {
    "test.yml" = "value: 123"
  }
  environment = [
    {
      name  = "USERNAME"
      value = "joe"
    },
    {
      name   = "PASSWORD"
      value  = "test"
      secret = true
    }
  ]
  job_spec = {
    namespace = "default"
    labels = {
      test = "123"
    }
    service_account = "default"
    containers      = [
      {
        image    = "perl:5.34.0"
        args     = ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]
        env = {}
        env_from = []
      }
    ]
  }
  bindings = {
    read  = []
    write = []
  }
}

resource "plural_infrastructure_stack" "stack-raw" {
  name       = "stack-tf-raw-${random_string.random.result}"
  type       = "TERRAFORM"
  approval   = true
  detach     = true
  cluster_id = data.plural_cluster.cluster.id
  repository = {
    id     = data.plural_git_repository.repository.id
    ref    = "main"
    folder = "terraform"
  }
  configuration = {
    image   = "hashicorp/terraform:1.8.1"
    version = "1.8.1"
  }
  files = {
    "test.yml" = "value: 123"
  }
  environment = [
    {
      name  = "USERNAME"
      value = "joe"
    },
    {
      name   = "PASSWORD"
      value  = "test"
      secret = true
    }
  ]
  job_spec = {
    namespace = "default"
    raw       = yamlencode({
      containers = [
        {
          name    = "pi"
          image   = "perl:5.34.0"
          command = ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]
        }
      ]
      restartPolicy = "Never"
    })
  }
  bindings = {
    read  = []
    write = []
  }
}

resource "plural_infrastructure_stack" "stack-empty" {
  name       = "stack-tf-empty-${random_string.random.result}"
  type       = "TERRAFORM"
  approval   = true
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
  files = {}
  environment = []
  job_spec = {
    namespace = "default"
    raw       = yamlencode({ test = true })
  }
  bindings = {}
}

resource "plural_infrastructure_stack" "stack-minimal" {
  name       = "stack-tf-minimal-${random_string.random.result}"
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
