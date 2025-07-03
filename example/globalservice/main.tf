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

resource "plural_global_service" "guestbook" {
  name      = "guestbook"
  service_id = "624bff88-05e3-45f6-bc3b-44708594e28e"
  distro = "AKS"
}
