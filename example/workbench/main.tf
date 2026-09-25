terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.39"
    }
  }
}

provider "plural" {
  use_cli = true
}

locals {
  # Avoid collisions with leftover tools from partial applies.
  name_prefix = "tf_"
}

data "plural_project" "default" {
  name = "default"
}

data "plural_git_repository" "hello" {
  url = "https://github.com/zreigz/tf-hello.git"
}

resource "plural_cloud_connection" "workbench" {
  name           = "${local.name_prefix}workbench_aws"
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
  name = "${local.name_prefix}workbench_observability"
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

data "plural_service_deployment" "console" {
  cluster = "mgmt"
  name    = "console"
}

resource "plural_monitor" "error_logs" {
  name            = "${local.name_prefix}error_logs"
  service_id      = data.plural_service_deployment.console.id
  description     = "Fires when the console logs too many errors or fatal messages."
  severity        = "MEDIUM"
  type            = "LOG"
  evaluation_cron = "*/15 * * * *"

  query = {
    log = {
      query       = "error OR fatal"
      bucket_size = "10m"
      duration    = "2h"
      operator    = "OR"
      facets = [
        { key = "namespace", value = "plrl-console" },
      ]
    }
  }

  threshold = {
    aggregate = "MAX"
    value     = 20
  }
}

resource "plural_monitor" "cpu_usage" {
  name            = "${local.name_prefix}cpu_usage"
  service_id      = data.plural_service_deployment.console.id
  workbench_id    = plural_workbench.full.id
  description     = "Fires when the console uses more than 2 CPU cores on average."
  prompt          = "Investigate the high console CPU usage and suggest a fix."
  severity        = "HIGH"
  type            = "METRICS"
  evaluation_cron = "*/10 * * * *"

  query = {
    metrics = {
      query    = "sum(rate(container_cpu_usage_seconds_total{namespace=\"plrl-console\", container!=\"\"}[5m]))"
      step     = "1m"
      duration = "1h"
    }
  }

  threshold = {
    aggregate = "AVG"
    value     = 2
  }

  modes = {
    budget = {
      cost   = 10
      tokens = 200000
    }
    kubernetes = {
      update             = true
      exclude_namespaces = ["kube-system", "kube-public"]
    }
  }
}

resource "plural_dashboard" "overview" {
  workbench_id = plural_workbench.full.id
  name         = "${local.name_prefix}console_overview"
  description  = "CPU, memory and restarts of the console on the mgmt cluster."

  # Inputs are referenced in graph markdown and datasource inputs as ${name},
  # which has to be escaped as $${name} in Terraform strings.
  inputs = [
    {
      name        = "namespace"
      label       = "Namespace"
      description = "Namespace the console is running in."
      type        = "TEXT"
      default     = "plrl-console"
    },
    {
      name        = "window"
      label       = "Rate window"
      description = "Window used to calculate CPU usage rate."
      type        = "SELECT"
      default     = "15m"
      options     = ["5m", "15m", "1h"]
    },
    {
      name        = "pod"
      label       = "Pod"
      description = "Filter graphs by console pod."
      type        = "SELECT"
      default     = ".*"
      datasource = {
        type = "LABELS"
        tool = "plrl_metric_label_search"
        input = jsonencode({
          metric = "container_cpu_usage_seconds_total"
          label  = "pod"
          query  = "console"
        })
      }
    },
  ]

  graphs = [
    {
      identifier = "overview"
      title      = "Overview"
      type       = "SECTION"
      options    = jsonencode({ collapsed = true })
      layout     = { x = 0, y = 0, w = 12, h = 1 }
    },
    {
      identifier  = "cpu"
      title       = "CPU usage"
      description = "CPU cores used by console pods."
      type        = "TIMESERIES"
      section_id  = "overview"
      layout      = { x = 0, y = 1, w = 12, h = 4 }
      datasource = {
        type  = "METRICS"
        tool  = "plrl_metrics"
        input = jsonencode({ query = "sum by (pod) (rate(container_cpu_usage_seconds_total{namespace=\"$${namespace}\", pod=~\"$${pod}\", container!=\"\"}[$${window}]))" })
      }
    },
    {
      identifier  = "memory"
      title       = "Memory usage"
      description = "Working set memory of console pods."
      type        = "TIMESERIES"
      section_id  = "overview"
      layout      = { x = 0, y = 5, w = 8, h = 4 }
      datasource = {
        type  = "METRICS"
        tool  = "plrl_metrics"
        input = jsonencode({ query = "sum by (pod) (container_memory_working_set_bytes{namespace=\"$${namespace}\", pod=~\"$${pod}\", container!=\"\"})" })
      }
    },
    {
      identifier  = "restarts"
      title       = "Restarts"
      description = "Container restarts of console pods in the selected window."
      type        = "STAT"
      layout      = { x = 8, y = 5, w = 4, h = 4 }
      datasource = {
        type  = "METRICS"
        tool  = "plrl_metrics"
        input = jsonencode({ query = "sum(increase(kube_pod_container_status_restarts_total{namespace=\"$${namespace}\", pod=~\"$${pod}\"}[$${window}]))" })
      }
    },
  ]
}
