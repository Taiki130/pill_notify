terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "5.33.0"
    }
    sops = {
      source  = "carlpett/sops"
      version = "0.7.2"
    }
  }
}

provider "github" {}
provider "sops" {}

resource "github_actions_secret" "secrets" {
  for_each = nonsensitive(data.sops_file.secrets.data)

  repository      = "pill_notify"
  secret_name     = each.key
  plaintext_value = each.value
}

data "sops_file" "secrets" {
  source_file = "secrets.yaml"
}
