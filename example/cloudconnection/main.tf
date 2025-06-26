terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.25"
    }
  }
}

provider "plural" {
  use_cli = true
}

resource "plural_cloud_connection" "aws" {
  cloud_provider = "aws"
  name          = "your-connection-name"
}
