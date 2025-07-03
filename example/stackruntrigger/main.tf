terraform {
  required_providers {
    plural = {
      source = "pluralsh/plural"
      version = "0.2.27"
    }
  }
}

provider "plural" {
  use_cli = true
}

data "plural_infrastructure_stack" "test" {
  name = "test-job"
}

resource "plural_stack_run_trigger" "trigger" {
  id = data.plural_infrastructure_stack.test.id
  retrigger_key = "test"
}
