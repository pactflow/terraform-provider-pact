variable "api_token" {
  type = string
}

variable "build_number" {
  type = string
}

provider "pact" {
  # host = "https://tf-acceptance.pact.dius.com.au"
  host = "https://tf-acceptance.test.pactflow.io"
  access_token = var.api_token
}

### API Tokens / Credentials

# NOTE: you probably don't want to use TF for managing these
resource "pact_token" "read_only" {
  type = "read-only"
  name = "Local dev token"
}

resource "pact_token" "read_write" {
  type = "read-write"
  name = "CI token"
}

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

resource "pact_role" "homer_admin" {
  role = "administrator"
  user = pact_user.homer.uuid
}