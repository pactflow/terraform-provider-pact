# Pact Broker - Terraform Provider

Terraform Provider for [Pact Broker](https://github.com/pact-foundation/pact_broker) and [Pactflow](https://pactflow.io).

[![Build Status](https://travis-ci.org/pactflow/terraform.svg?branch=master)](https://travis-ci.org/pactflow/terraform)
[![Coverage Status](https://coveralls.io/repos/github/pactflow/terraform/badge.svg?branch=master)](https://coveralls.io/github/pactflow/terraform?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pactflow/terraform)](https://goreportcard.com/report/github.com/pactflow/terraform)
[![GoDoc](https://godoc.org/github.com/pactflow/terraform?status.svg)](https://godoc.org/github.com/pactflow/terraform)
[![slack](http://slack.pact.io/badge.svg)](http://slack.pact.io)

## Introduction

<p align="center">
  <img width="880" src="https://raw.githubusercontent.com/pactflow/terraform/master/examples/tf-run.svg?sanitize=true">
</p>

Example:
```hcl
terraform {
  required_providers {
    pact = {
      source = "pactflow/pact"
      version = "0.1.2"
    }
  }
}

variable "token" {
  type = string
}

provider "pact" {
  # For the OSS Broker
  # host = "http://localhost"
  # basic_auth_username = "pact_broker"
  # basic_auth_password = "pact_broker"

  # For a Pactflow Broker
  host = "https://mybroker.pact.dius.com.au"
  access_token = var.token
}

# Create a Pacticipant for our Product API
resource "pact_application" "product_api" {
  name = "product_api"
  repository_url = "github.com/foo/product_api"
}

# Create a Pacticipant for our Ecommerce Website
resource "pact_application" "ecommerce_web" {
  name = "ecommerce_web"
  repository_url = "github.com/foo/ecommerce_web"
}

# Create a secret
resource "pact_secret" "ci_token" {
  name = "CIToken"
  description = "Jenkins token for Pactflow"
  value = "my super secret value"
}

# Create a webhook to trigger a build when a contract changes
# Uses the Jenkins token from above
resource "pact_webhook" "product_api_changed" {
  description = "Trigger build when a contract changes between ecommerce_web and product_api"
  webhook_provider = {
    name = "product_api"
  }
  webhook_consumer = {
    name = "ecommerce_web"
  }
  request {
    url = "https://foo.com/some/endpoint"
    method = "POST"
    headers = {
      "Content-Type" = "application/json",
      "Authoriation" = "Bearer $${user.CIToken}"
    }
    body = <<EOF
{
  "pact": "$${pactbroker.pactUrl}"
}
EOF
  }

  events = ["contract_content_changed", "contract_published"]
  depends_on = [pact_application.product_api, pact_application.ecommerce_web, pact_secret.ci_token]
}

# Create a user
resource "pact_user" "somebody" {
  name = "Tester McPerson"
  email = "pact@dius.com.au"
  active = true
}

# Assign admin role to user
resource "pact_role" "somebody_admin" {
  role = "administrator"
  user = pact_user.somebody.uuid
}
```


## Installing

Download the latest [release](https://github.com/pactflow/terraform/releases) and install into your Terraform plugin directory.

### Linux or Mac OSX

Run the following to have the provider installed for you automatically:

```sh
curl -fsSL https://raw.githubusercontent.com/pactflow/terraform/master/scripts/install.sh | bash
```

### Windows

Dowload the plugin to `%APPDATA%\terraform.d\plugins`.


### Installation notes

To use a released provider in your Terraform environment, run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider. To specify a particular provider version when installing released providers, see the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

To instead use a custom-built provider in your Terraform environment (e.g. the provider binary from the build instructions above), follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) After placing the custom-built provider into your plugins directory,  run `terraform init` to initialize it.

For either installation method, documentation about the provider specific configuration options can be found on the [provider's website](https://www.terraform.io/docs/providers/aws/index.html).

## Using the plugin


| Plugin      | Type     | Platform Support       | Description |
| ----------- | -------- | ---------------------- | ----------- |
| [Pact](docs/index.md)        | Provider | Pact Broker + Pactflow | Configures a target Pact Broker (such as a pactflow.io account) |
| [Pacticipant](docs/resources/pacticipant.md)  | Resource | Pact Broker + Pactflow | Create applications (known as Pacticipants) |
| [Webhook](docs/resources/webhook.md)     | Resource | Pact Broker + Pactflow | Configures a webhook to trigger on certain platform events |
| [Secret](docs/resources/secret.md)      | Resource | Pactflow               | Create an encrypted secret for use in Webhooks |
| [API Token](docs/resources/token.md)   | Resource | Pactflow               | Manage Pactflow API Tokens |
| [Users](docs/resources/user.md)   | Resource | Pactflow               | Manage Pactflow Users |
| [Roles](docs/resources/role.md)   | Resource | Pactflow               | Manage Pactflow Roles |

See our [Docs](./docs) folder for all plugins.

## Importing Resources

All resources (except for Roles) support [importing](https://www.terraform.io/docs/import/usage.html).

## Developing

### Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.10+
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

### Building locally
*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH (e.g. `$HOME/development/terraform-providers/`).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `./bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-aws
...
```

## Roadmap

Plan for the next few months:

- [x] Pacticipants
- [x] Webhooks
- [x] Secrets (Pactflow only)
- [x] API Tokens (Pactflow only)
- [x] User + Basic Role Management (Pactflow only)
- [ ] Teams
- [ ] Assign users+teams
- [ ] System Users
- [ ] Role-based access ()
- [ ] Role creation
- [ ] Better error messages for HTTP / runtime failures
- [x] Acceptance tests
- [ ] Better code coverage
- [ ] Extract `Client` into separate SDK package
- [ ] Publish 1.0.0

Want to see something else here? Have you say on our [Pact Feature Request](https://pact.canny.io/feature-requests/p/create-a-terraform-provider) board.