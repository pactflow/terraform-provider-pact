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
  host = "https://tf-acceptance.pact.dius.com.au"
  access_token = var.api_token
}

variable "api_token" {
  type = string
}

variable "build_number" {
  type = string
}


resource "pact_pacticipant" "AdminUI" {
  name = "AdminUI${var.build_number}"
  repository_url = "github.com/foo/admin"
}

resource "pact_pacticipant" "GraphQLAPI" {
  name = "GraphQLAPI${var.build_number}"
  repository_url = "github.com/foo/api"
}

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

resource "pact_secret" "jenkins_token" {
  name = "JenkinsTriggerToken${var.build_number}"
  description = "API token to trigger Jenkins builds"
  value = "super secret thing"
}

resource "pact_token" "read_only" {
  type = "read-only"
  name = "Local dev token"
}

resource "pact_token" "read_write" {
  type = "read-write"
  name = "CI token"
}

resource "pact_user" "somebody" {
  name = "Tester McPerson"
  email = "mfellows+test@dius.com.au"
  active = true
}

resource "pact_role" "somebody_admin" {
  role = "administrator"
  user = pact_user.somebody.uuid
}