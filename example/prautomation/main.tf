terraform {
  required_providers {
    plural = {
      source = "pluralsh/plural"
      version = "0.2.1"
    }
  }
}

provider "plural" {
  use_cli = true
}

data "plural_pr_automation" "automation" {
  name = "pr-test"
}

# resource "plural_pr_automation_trigger" "trigger" {
#   pr_automation_id = data.plural_pr_automation.automation.id
#   pr_automation_branch = "prautomation"
#   context = {
#     version: "v0.0.0"
#   }
# }
