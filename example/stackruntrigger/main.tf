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

resource "plural_stack_run_trigger" "trigger" {
  id = "ecc8966a-edfe-4e48-b5ea-f87c3e97d0a3"
}
