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

resource "plural_pr_automation_trigger" "trigger" {
  pr_automation_id = "1cc7483e-78dd-4470-9ae4-6eb2c8cc1785"
  pr_automation_branch = "prautomation"
  context = {
    version: "v0.0.0"
  }
}
