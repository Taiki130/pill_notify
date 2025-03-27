terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "6.0.0"
    }
    sops = {
      source  = "carlpett/sops"
      version = "1.2.0"
    }
    spacelift = {
      source  = "spacelift-io/spacelift"
      version = "1.20.4"
    }
  }
}

provider "github" {}
provider "sops" {}
provider "spacelift" {
  api_key_endpoint = "https://taiki130.app.spacelift.io"
  api_key_id       = data.sops_file.tf_secrets.data["SPACELIFT_API_KEY_ID"]
  api_key_secret   = data.sops_file.tf_secrets.data["SPACELIFT_API_KEY_SECRET"]
}

resource "spacelift_module" "github_actions_secret" {
  name               = "github_actions_secret"
  terraform_provider = "github"
  administrative     = false
  branch             = "main"
  description        = "create github_actions_secret"
  repository         = "pill_notify"
  project_root       = "tofu/github"
  space_id           = "root"
}

module "github_actions_secret" {
  source  = "spacelift.io/taiki130/github_actions_secret/github"
  version = "1.1.2"

  secrets = nonsensitive(data.sops_file.secrets.data)
}

data "sops_file" "secrets" {
  source_file = "secrets.yaml"
}

data "sops_file" "tf_secrets" {
  source_file = "tf_secrets.yaml"
}
