terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.35"
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

resource "plural_cloud_connection" "workbench" {
  name           = "workbench-aws"
  cloud_provider = "AWS"
  configuration = {
    aws = {
      access_key_id     = "replace-with-access-key-id"
      secret_access_key = "replace-with-secret-access-key"
      region            = "us-east-1"
    }
  }
}

resource "plural_observability_webhook" "workbench" {
  name = "workbench-observability"
  type = "NEWRELIC"
}

resource "plural_workbench_tool" "echo" {
  name       = "echo"
  tool       = "HTTP"
  project_id = data.plural_project.default.id
  configuration = {
    http = {
      url     = "https://httpbin.org/post"
      method  = "POST"
      headers = { "Content-Type" = "application/json" }
      body    = "{\"message\": \"{{input.message}}\"}"
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
      url    = "https://httpbin.org/anything/{{input.id}}"
      method = "GET"
      input_schema = jsonencode({
        type = "object"
        properties = {
          id = {
            type        = "string"
            description = "Optional id to append to the request path."
          }
        }
      })
    }
  }
}

resource "plural_workbench_tool" "elastic" {
  name       = "elastic"
  tool       = "ELASTIC"
  project_id = data.plural_project.default.id
  configuration = {
    elastic = {
      url      = "https://my-elastic-instance.es.io:9200"
      username = "replace-with-username"
      password = "replace-with-password"
      index    = "logs-*"
    }
  }
}

resource "plural_workbench_tool" "datadog" {
  name       = "datadog"
  tool       = "DATADOG"
  project_id = data.plural_project.default.id
  configuration = {
    datadog = {
      site    = "datadoghq.com"
      api_key = "replace-with-api-key"
      app_key = "replace-with-app-key"
    }
  }
}

resource "plural_workbench_tool" "prometheus" {
  name       = "prometheus"
  tool       = "PROMETHEUS"
  project_id = data.plural_project.default.id
  configuration = {
    prometheus = {
      url      = "https://prometheus.example.com"
      username = "replace-with-username"
      password = "replace-with-password"
    }
  }
}

resource "plural_workbench_tool" "loki" {
  name       = "loki"
  tool       = "LOKI"
  project_id = data.plural_project.default.id
  configuration = {
    loki = {
      url      = "https://loki.example.com"
      username = "replace-with-username"
      password = "replace-with-password"
    }
  }
}

resource "plural_workbench_tool" "tempo" {
  name       = "tempo"
  tool       = "TEMPO"
  project_id = data.plural_project.default.id
  configuration = {
    tempo = {
      url      = "https://tempo.example.com"
      username = "replace-with-username"
      password = "replace-with-password"
    }
  }
}

resource "plural_workbench_tool" "jaeger" {
  name       = "jaeger"
  tool       = "JAEGER"
  project_id = data.plural_project.default.id
  configuration = {
    jaeger = {
      url      = "https://jaeger.example.com"
      username = "replace-with-username"
      password = "replace-with-password"
    }
  }
}

resource "plural_workbench_tool" "splunk" {
  name       = "splunk"
  tool       = "SPLUNK"
  project_id = data.plural_project.default.id
  configuration = {
    splunk = {
      url   = "https://splunk.example.com"
      token = "replace-with-token"
    }
  }
}

resource "plural_workbench_tool" "dynatrace" {
  name       = "dynatrace"
  tool       = "DYNATRACE"
  project_id = data.plural_project.default.id
  configuration = {
    dynatrace = {
      url            = "https://my-env.live.dynatrace.com"
      platform_token = "replace-with-platform-token"
    }
  }
}

resource "plural_workbench_tool" "cloudwatch" {
  name                = "cloudwatch_default"
  tool                = "CLOUDWATCH"
  project_id          = data.plural_project.default.id
  cloud_connection_id = plural_cloud_connection.workbench.id
  configuration = {
    cloudwatch = {
      region          = "us-east-1"
      log_group_names = ["/aws/eks/default/application"]
    }
  }
}

resource "plural_workbench_tool" "azure" {
  name       = "azure"
  tool       = "AZURE"
  project_id = data.plural_project.default.id
  configuration = {
    azure = {
      subscription_id = "replace-with-subscription-id"
      tenant_id       = "replace-with-tenant-id"
      client_id       = "replace-with-client-id"
      client_secret   = "replace-with-client-secret"
    }
  }
}

resource "plural_workbench_tool" "linear" {
  name       = "linear"
  tool       = "LINEAR"
  project_id = data.plural_project.default.id
  configuration = {
    linear = {
      access_token = "replace-with-access-token"
    }
  }
}

resource "plural_workbench_tool" "atlassian" {
  name       = "atlassian"
  tool       = "ATLASSIAN"
  project_id = data.plural_project.default.id
  configuration = {
    atlassian = {
      email     = "user@example.com"
      api_token = "replace-with-api-token"
    }
  }
}

resource "plural_workbench_tool" "sentry" {
  name       = "sentry"
  tool       = "SENTRY"
  project_id = data.plural_project.default.id
}

resource "plural_workbench_tool" "mcp" {
  name          = "mcp"
  tool          = "MCP"
  project_id    = data.plural_project.default.id
  mcp_server_id = "294e2211-e379-40f8-88a4-086f00cd0a31"
}

resource "plural_workbench_tool" "cloud" {
  name                = "cloud"
  tool                = "CLOUD"
  project_id          = data.plural_project.default.id
  cloud_connection_id = plural_cloud_connection.workbench.id
}

resource "plural_workbench_tool" "minimal" {
  name = "minimal"
  tool = "HTTP"
}

resource "plural_workbench" "full" {
  name          = "full"
  description   = "Sample workbench with all tool types."
  system_prompt = "You are a helpful assistant."
  project_id    = data.plural_project.default.id
  repository_id = data.plural_git_repository.hello.id
  agent_runtime = "mgmt/gemini"
  configuration = {
    coding = {
      mode         = "WRITE"
      repositories = ["https://github.com/pluralsh/echo-skill"]
    }
    infrastructure = {
      stacks     = true
      services   = true
      kubernetes = true
    }
    observability = {
      logs    = true
      metrics = true
    }
  }
  tool_ids = [
    plural_workbench_tool.echo.id,
    plural_workbench_tool.status.id,
    plural_workbench_tool.elastic.id,
    plural_workbench_tool.datadog.id,
    plural_workbench_tool.prometheus.id,
    plural_workbench_tool.loki.id,
    plural_workbench_tool.tempo.id,
    plural_workbench_tool.jaeger.id,
    plural_workbench_tool.splunk.id,
    plural_workbench_tool.dynatrace.id,
    plural_workbench_tool.cloudwatch.id,
    plural_workbench_tool.azure.id,
    plural_workbench_tool.linear.id,
    plural_workbench_tool.atlassian.id,
    plural_workbench_tool.sentry.id,
    plural_workbench_tool.mcp.id,
    plural_workbench_tool.cloud.id,
  ]
}

resource "plural_workbench" "minimal" {
  name = "minimal"
}

resource "plural_workbench_webhook" "alerts" {
  workbench_id = plural_workbench.full.id
  name         = "alerts"
  webhook_id   = plural_observability_webhook.workbench.id
  prompt       = "Investigate this alert and summarize root cause."

  matches = {
    regex            = "severity.*(critical|high)"
    case_insensitive = true
  }
}

resource "plural_workbench_cron" "daily_check" {
  workbench_id = plural_workbench.full.id
  crontab      = "0 9 * * 1-5"
  prompt       = "Run a morning health check and summarize notable issues."
}
