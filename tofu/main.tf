terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "5.42.0"
    }
    sops = {
      source  = "carlpett/sops"
      version = "1.0.0"
    }
  }
}

provider "github" {}
provider "sops" {}

resource "github_actions_secret" "secrets" {
  for_each = var.github_actions_secrets

  repository      = "pill_notify"
  secret_name     = each.key
  plaintext_value = data.sops_file.secrets.data[each.key]
}

data "sops_file" "secrets" {
  source_file = "secrets.yaml"
}

variable "github_actions_secrets" {
  type = set(string)
  default = [
    "LINE_TOKEN",
    "MESSAGE",
    "IMAGE_THUMBNAIL_URL",
    "IMAGE_FULLSIZE_URL",
  ]
}
