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

provider "pact" {
  host = "http://localhost"
  basic_auth_username = "pact_broker"
  basic_auth_password = "pact_broker"
}

resource "pact_pacticipant" "AdminUI" {
  name = "AdminUI"
  repository_url = "github.com/foo/admin"
}

resource "pact_pacticipant" "GraphQLAPI" {
  name = "GraphQLAPI"
  repository_url = "github.com/foo/api"
}

data "pact_pacticipants" "all" {
  depends_on = [pact_pacticipant.AdminUI, pact_pacticipant.GraphQLAPI]

  lifecycle {
    postcondition {
      condition     = alltrue([for name in ["AdminUI", "GraphQLAPI"] : contains([for p in self.pacticipants : p.name], name)])
      error_message = "pact_pacticipants data source did not list the pacticipants created above"
    }
  }
}

resource "pact_webhook" "ui_changed" {
  description = "Trigger an API build when the UI changes"
  webhook_provider = {
    name = "GraphQLAPI"
  }
  webhook_consumer = {
    name = "AdminUI"
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

resource "pact_webhook" "nonjson" {
  description = "POST non-JSON"
  request {
    url = "https://foo.com/some/endpoint"
    method = "POST"
    body = "json={\"parameter\": [{\"ame\": \"TAG_AND_PUSH\", \"value\": \"false\"}]}"
  }
  events = ["contract_content_changed"]
}
