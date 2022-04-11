# Required as of Terraform version 0.0.13
terraform {
  required_version = ">= 0.13"
  required_providers {
    pact = {
      source  = "github.com/pactflow/pact"
      version = "0.0.1"
    }
  }
}

### Provider configuration

provider "pact" {
  host = "https://tf-acceptance.pactflow.io"
  access_token = var.api_token
  # host = "http://localhost:9292"
}

variable "api_token" {
  type = string
}

variable "build_number" {
  type = string
}

### Secrets

resource "pact_secret" "jenkins_token" {
  name = "JenkinsTriggerToken${var.build_number}"
  description = "API token to trigger Jenkins builds"
  value = "super secret thing"
  team = pact_team.Simpsons.uuid
}

### Applications

resource "pact_pacticipant" "example" {
  name = "pactflow-example-consumer${var.build_number}"
}

resource "pact_pacticipant" "AdminUI" {
  name = "AdminUI${var.build_number}"
  repository_url = "github.com/foo/admin"
}

resource "pact_pacticipant" "GraphQLAPI" {
  name = "GraphQLAPI${var.build_number}"
  repository_url = "github.com/foo/api"
}

### Teams

resource "pact_team" "Simpsons" {
  name = "Simpsons${var.build_number}"
  users = [
    pact_user.homer.uuid
  ]
}

resource "pact_team" "Futurama" {
  name = "Futurama${var.build_number}"
  pacticipants = [
    pact_pacticipant.GraphQLAPI.name,
    pact_pacticipant.example.name
  ]
  users = [
    pact_user.bender_system_user.uuid
  ]
  administrators = [
    pact_user.homer.uuid
  ]
}

resource "pact_team" "Cartoons" {
  name = "Cartoons${var.build_number}"
  users = [
    pact_user.bender_system_user.uuid,
    pact_user.homer.uuid
  ]
}

### Users

resource "pact_user" "homer" {
  name = "Homer Simpson"
  email = "matt+tfacceptance1${var.build_number}@pactflow.io"
  active = true
  roles = [
    "c1878b8e-d09e-11ea-8fde-af02c4677eb7",
    "9fa50562-a42b-4771-aa8e-4bb3d623ae60",
    "e9282e22-416b-11ea-a16e-57ee1bb61d18",
    pact_role.special_role.uuid,
    "cf75d7c2-416b-11ea-af5e-53c3b1a4efd8" # Admin
  ]
}

resource "pact_user" "bender_system_user" {
  name = "Bender Rodríguez"
  email = "matt+tfacceptance2${var.build_number}@pactflow.io"
  type = "system"
  active = true
}

### Webhooks

resource "pact_webhook" "ui_changed" {
  description = "Trigger an API build when the UI changes ${var.build_number}"
  team = pact_team.Simpsons.uuid
  webhook_provider = {
    name = "GraphQLAPI${var.build_number}"
  }
  webhook_consumer = {
    name = "AdminUI${var.build_number}"
  }
  request {
    url = "https://foo.com/some/endpoint"
    method = "POST"
    username = "test"
    password = "password1"
    headers = {
      "X-Content-Type" = "application/json"
    }
    body = <<EOF
{
  "pact": "$${pactbroker.pactUrl}"
}
EOF
  }

  events = ["contract_content_changed", "contract_published"]
  depends_on = [pact_pacticipant.AdminUI, pact_pacticipant.GraphQLAPI]
}

### Roles and Permissions

resource "pact_role" "special_role" {
  name = "specialrole${var.build_number}"
  scopes = [
    "user:manage:*",
    "team:manage:*",
    "user:invite",
    "system_account:manage:*",
    "system_account:read:*",
    "user:read:*",
    "team:read:*",
    "contract_data:manage:*",
    "contract_data:read:*",
    "contract_data:bulk_delete:*",
    "webhook:manage:*",
    "secret:manage:*",
    "role:manage:*",
    "role:read:*",
    "token:manage:own",
    "read_token:manage:own"
  ]
}

### Authentication

resource "pact_authentication" "authentication" {
  github_organizations = ["DiUS", "pactflow"]
  google_domains = ["dius.com.au", "pactflow.io"]
}

### Environments

resource "pact_environment" "staging" {
  name = "staging${var.build_number}"
  display_name = "Staging Environment"
  production = false
  teams = [pact_team.Simpsons.uuid]
}