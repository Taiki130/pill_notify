terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "6.0.0"
    }
  }
}

resource "github_actions_secret" "secrets" {
  for_each = var.secrets

  repository      = "pill_notify"
  secret_name     = each.key
  plaintext_value = each.value
}

variable "secrets" {
  type    = map(string)
  default = null
}
