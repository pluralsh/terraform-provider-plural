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

resource "plural_workbench" "sample" {
  name        = "sample-workbench"
  description = "Sample workbench with two HTTP tools."
  project_id  = data.plural_project.default.id
  tool_ids    = [plural_workbench_tool.echo.id, plural_workbench_tool.status.id]
}
