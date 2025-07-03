terraform {
  required_providers {
    plural = {
      source  = "pluralsh/plural"
      version = "0.2.27"
    }
  }
}

provider "plural" {
  use_cli = true
}

data "plural_user" "user" {
  email = "sebastian@plural.sh"
}

resource "plural_shared_secret" "mysecret" {
  name   = "mysecret"
  secret = "password"
  notification_bindings = [
    { user_id = data.plural_user.user.id }
  ]
}

resource "null_resource" "default" {
  provisioner "local-exec" {
    command = "echo name:${plural_shared_secret.mysecret.name}"
  }
}

output "secretoutput" {
  value     = plural_shared_secret.mysecret.secret
  sensitive = true
}
