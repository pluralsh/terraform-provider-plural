terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.37"
    }
  }
}

provider "plural" {
  use_cli = true
}

locals {
  # Avoid collisions with leftover tools from partial applies.
  name_prefix = "tf-example-"
}

data "plural_project" "default" {
  name = "default"
}

data "plural_git_repository" "hello" {
  url = "https://github.com/zreigz/tf-hello.git"
}

resource "plural_cloud_connection" "workbench" {
  name           = "${local.name_prefix}workbench-aws"
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
  name = "${local.name_prefix}workbench-observability"
  type = "NEWRELIC"
}

resource "plural_workbench_tool" "echo" {
  name       = "${local.name_prefix}echo"
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
  name       = "${local.name_prefix}status"
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
  name       = "${local.name_prefix}elastic"
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
  name       = "${local.name_prefix}datadog"
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
  name       = "${local.name_prefix}prometheus"
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
  name       = "${local.name_prefix}loki"
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
  name       = "${local.name_prefix}tempo"
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
  name       = "${local.name_prefix}jaeger"
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
  name       = "${local.name_prefix}splunk"
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
  name       = "${local.name_prefix}dynatrace"
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
  name                = "${local.name_prefix}cloudwatch_default"
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
  name       = "${local.name_prefix}azure"
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
  name       = "${local.name_prefix}linear"
  tool       = "LINEAR"
  project_id = data.plural_project.default.id
  configuration = {
    linear = {
      access_token = "replace-with-access-token"
    }
  }
}

resource "plural_workbench_tool" "atlassian" {
  name       = "${local.name_prefix}atlassian"
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
  name       = "${local.name_prefix}sentry"
  tool       = "SENTRY"
  project_id = data.plural_project.default.id
  configuration = {
    sentry = {
      url          = "https://sentry.io"
      access_token = "replace-with-access-token"
    }
  }
}

resource "plural_workbench_tool" "exa" {
  name       = "${local.name_prefix}exa"
  tool       = "EXA"
  project_id = data.plural_project.default.id
  configuration = {
    exa = {
      api_key = "replace-with-api-key"
    }
  }
}

resource "plural_workbench_tool" "github" {
  name       = "${local.name_prefix}github"
  tool       = "GITHUB"
  project_id = data.plural_project.default.id
  configuration = {
    github = {
      access_token = "replace-with-access-token"
      toolset      = "default"
    }
  }
}

resource "plural_workbench_tool" "slack" {
  name       = "${local.name_prefix}slack"
  tool       = "SLACK"
  project_id = data.plural_project.default.id
  configuration = {
    slack = {
      bot_token = "replace-with-bot-token"
    }
  }
}

resource "plural_workbench_tool" "teams" {
  name       = "${local.name_prefix}teams"
  tool       = "TEAMS"
  project_id = data.plural_project.default.id
  configuration = {
    teams = {
      client_id     = "replace-with-client-id"
      client_secret = "replace-with-client-secret"
      tenant_id     = "replace-with-tenant-id"
    }
  }
}

resource "plural_workbench_tool" "gitlab" {
  name       = "${local.name_prefix}gitlab"
  tool       = "GITLAB"
  project_id = data.plural_project.default.id
  configuration = {
    gitlab = {
      token = "replace-with-token"
    }
  }
}

resource "plural_workbench_tool" "bitbucket" {
  name       = "${local.name_prefix}bitbucket"
  tool       = "BITBUCKET"
  project_id = data.plural_project.default.id
  configuration = {
    bitbucket = {
      token = "replace-with-token"
    }
  }
}

resource "plural_workbench_tool" "bitbucket_datacenter" {
  name       = "${local.name_prefix}bitbucket_datacenter"
  tool       = "BITBUCKET_DATACENTER"
  project_id = data.plural_project.default.id
  configuration = {
    bitbucket_datacenter = {
      url   = "https://bitbucket.example.com"
      token = "replace-with-token"
    }
  }
}

resource "plural_workbench_tool" "azure_devops" {
  name       = "${local.name_prefix}azure_devops"
  tool       = "AZURE_DEVOPS"
  project_id = data.plural_project.default.id
  configuration = {
    azure_devops = {
      token = "replace-with-token"
    }
  }
}

resource "plural_workbench_tool" "pagerduty" {
  name       = "${local.name_prefix}pagerduty"
  tool       = "PAGERDUTY"
  project_id = data.plural_project.default.id
  configuration = {
    pagerduty = {
      api_token = "replace-with-api-token"
    }
  }
}

resource "plural_workbench_tool" "opensearch" {
  name       = "${local.name_prefix}opensearch"
  tool       = "OPENSEARCH"
  project_id = data.plural_project.default.id
  configuration = {
    opensearch = {
      host             = "https://search-example.us-east-1.es.amazonaws.com"
      index            = "logs-*"
      aws_region       = "us-east-1"
      use_pod_identity = true
    }
  }
}

resource "plural_workbench_tool" "lambda" {
  name                = "${local.name_prefix}lambda"
  tool                = "LAMBDA"
  project_id          = data.plural_project.default.id
  cloud_connection_id = plural_cloud_connection.workbench.id
  configuration = {
    lambda = {
      lambda_arn  = "arn:aws:lambda:us-east-1:123456789012:function:example"
      description = "Invoke the example Lambda function."
      input_schema = jsonencode({
        type = "object"
        properties = {
          payload = { type = "string" }
        }
      })
    }
  }
}

resource "plural_workbench_tool" "cloud_run" {
  name                = "${local.name_prefix}cloud_run"
  tool                = "CLOUD_RUN"
  project_id          = data.plural_project.default.id
  cloud_connection_id = plural_cloud_connection.workbench.id
  configuration = {
    cloud_run = {
      identifier  = "projects/example/locations/us-central1/services/example"
      description = "Invoke the example Cloud Run service."
      input_schema = jsonencode({
        type = "object"
        properties = {
          payload = { type = "string" }
        }
      })
    }
  }
}

resource "plural_workbench_tool" "azure_function" {
  name                = "${local.name_prefix}azure_function"
  tool                = "AZURE_FUNCTION"
  project_id          = data.plural_project.default.id
  cloud_connection_id = plural_cloud_connection.workbench.id
  configuration = {
    azure_function = {
      identifier  = "example-function"
      description = "Invoke the example Azure Function."
      input_schema = jsonencode({
        type = "object"
        properties = {
          payload = { type = "string" }
        }
      })
    }
  }
}

resource "plural_workbench_tool" "docker" {
  name       = "${local.name_prefix}docker"
  tool       = "DOCKER"
  project_id = data.plural_project.default.id
  configuration = {
    docker = {
      url      = "https://registry-1.docker.io"
      provider = "BASIC"
      auth = {
        basic = {
          username = "replace-with-username"
          password = "replace-with-password"
        }
      }
    }
  }
}

resource "plural_workbench_tool" "mcp" {
  name          = "${local.name_prefix}mcp"
  tool          = "MCP"
  project_id    = data.plural_project.default.id
  mcp_server_id = "294e2211-e379-40f8-88a4-086f00cd0a31"
}

resource "plural_workbench_tool" "cloud" {
  name                = "${local.name_prefix}cloud"
  tool                = "CLOUD"
  project_id          = data.plural_project.default.id
  cloud_connection_id = plural_cloud_connection.workbench.id
}

resource "plural_workbench_tool" "minimal" {
  name = "${local.name_prefix}minimal"
  tool = "HTTP"
}

resource "plural_workbench" "full" {
  name          = "${local.name_prefix}full"
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
    plural_workbench_tool.exa.id,
    plural_workbench_tool.github.id,
    plural_workbench_tool.slack.id,
    plural_workbench_tool.teams.id,
    plural_workbench_tool.gitlab.id,
    plural_workbench_tool.bitbucket.id,
    plural_workbench_tool.bitbucket_datacenter.id,
    plural_workbench_tool.azure_devops.id,
    plural_workbench_tool.pagerduty.id,
    plural_workbench_tool.opensearch.id,
    plural_workbench_tool.lambda.id,
    plural_workbench_tool.cloud_run.id,
    plural_workbench_tool.azure_function.id,
    plural_workbench_tool.docker.id,
    plural_workbench_tool.mcp.id,
    plural_workbench_tool.cloud.id,
  ]
}

resource "plural_workbench" "minimal" {
  name = "${local.name_prefix}minimal"
}

resource "plural_workbench_webhook" "alerts" {
  workbench_id = plural_workbench.full.id
  name         = "${local.name_prefix}alerts"
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
