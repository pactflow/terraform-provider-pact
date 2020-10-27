# Required as of Terraform version 0.0.13
terraform {
  required_version = ">= 0.13"
  required_providers {
    pact = {
      source  = "github.com/pactflow/pact"
      versions = ["0.0.1"]
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