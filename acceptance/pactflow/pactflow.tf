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
  host = "http://localhost:9292"
  access_token = var.api_token
}

variable "api_token" {
  type = string
}

variable "build_number" {
  type = string
}


### API Tokens / Credentials

# NOTE: you probably don't want to use TF for managing these
# resource "pact_token" "read_only" {
#   type = "read-only"
#   name = "Local dev token"
# }

# resource "pact_token" "read_write" {
#   type = "read-write"
#   name = "CI token"
# }

### Secrets

resource "pact_secret" "jenkins_token" {
  name = "JenkinsTriggerToken${var.build_number}"
  description = "API token to trigger Jenkins builds"
  value = "super secret thing"
}

### Applications

resource "pact_pacticipant" "example" {
  name = "pactflow-example-consumer"
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
  name = "Simpsons"
}

resource "pact_team" "Futurama" {
  name = "Futurama"
  pacticipants = [
    pact_pacticipant.GraphQLAPI.name,
    pact_pacticipant.example.name
  ]
}

resource "pact_team" "Cartoons" {
  name = "Cartoons"
}

### Users

resource "pact_user" "homer" {
  name = "Homer Simpson"
  email = "rholshausen@dius.com.au"
  active = true
  roles = [
    "c1878b8e-d09e-11ea-8fde-af02c4677eb7",
    "9fa50562-a42b-4771-aa8e-4bb3d623ae60",
    "e9282e22-416b-11ea-a16e-57ee1bb61d18",
    "cf75d7c2-416b-11ea-af5e-53c3b1a4efd8" # Admin
  ]
}

resource "pact_user" "bender_system_user" {
  name = "Bender Rodríguez"
  email = "bskurrie@dius.com.au"
  type = "system"
  active = true
}

### Assign users to Teams

resource "pact_team_assignment" "TeamFuturama" {
  team = pact_team.Futurama.uuid
  users = [
    pact_user.bender_system_user.uuid,
    pact_user.homer.uuid
  ]
}

resource "pact_team_assignment" "TeamSimpsons" {
  team = pact_team.Simpsons.uuid
  users = [
    pact_user.homer.uuid
  ]
}

resource "pact_team_assignment" "TeamCartoons" {
  team = pact_team.Cartoons.uuid
  users = [
    pact_user.bender_system_user.uuid,
    pact_user.homer.uuid
  ]
}

### Webhooks

resource "pact_webhook" "ui_changed" {
  description = "Trigger an API build when the UI changes"
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
  name = "specialrole"
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

# NOTE: legacy resource has changed name from previous versions
# resource "pact_role_v1" "homer_admin" {
#   role = "administrator"
#   user = pact_user.homer.uuid
# }