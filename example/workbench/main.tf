terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.32"
    }
  }
}

provider "plural" {
  use_cli = true
}

data "plural_project" "default" {
  name = "default"
}

data "plural_git_repository" "hello" {
  url = "https://github.com/zreigz/tf-hello.git"
}

resource "plural_workbench_tool" "echo" {
  name       = "echo"
  tool       = "HTTP"
  project_id = data.plural_project.default.id
  configuration = {
    http = {
      url          = "https://httpbin.org/post"
      method       = "POST"
      headers      = { "Content-Type" = "application/json" }
      body         = "{\"message\": \"{{input.message}}\"}"
      input_schema = jsonencode({
        type = "object"
        properties = {
          message = {
            type = "string"
          }
        }
        required = ["message"]
      })
    }
  }
}

resource "plural_workbench_tool" "status" {
  name       = "status"
  tool       = "HTTP"
  project_id = data.plural_project.default.id
  configuration = {
    http = {
      url          = "https://httpbin.org/anything/{{input.id}}"
      method       = "GET"
      input_schema = jsonencode({
        type = "object"
        properties = {
          id = {
            type = "string"
            description = "Optional id to append to the request path."
          }
        }
      })
    }
  }
}

resource "plural_workbench_tool" "minimal" {
  name       = "minimal"
  tool       = "HTTP"
}

resource "plural_workbench" "full" {
  name        = "full"
  description = "Sample workbench with two HTTP tools and other optional fields set."
  system_prompt = "You are a helpful assistant."
  project_id  = data.plural_project.default.id
  repository_id = data.plural_git_repository.hello.id
  agent_runtime = "mgmt/gemini"
  configuration = {
    coding = {
      mode = "WRITE"
      repositories = ["https://github.com/pluralsh/echo-skill"]
    }
    infrastructure = {
      stacks = true
      services = true
      kubernetes = true
    }
  }
  # skills = {
  #   ref = {
  #     ref = "main"
  #     folder = "examples/echo-skill"
  #     files = ["skill.py"]
  #   }
  # }
  tool_ids    = [plural_workbench_tool.echo.id, plural_workbench_tool.status.id]
}

resource "plural_workbench" "minimal" {
  name        = "minimal"
}
