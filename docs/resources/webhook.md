# Webhook Resource

This resource manages the lifecycle of a _Webhook_.

Webhooks allow you to trigger an HTTP request when a pact is changed, a pact is published, or a verification is published. The most common use case for webhooks is to trigger a provider build every time a pact changes, and to trigger a consumer build every time a verification is published.

Webhooks can be used in conjunction with the [can-i-deploy](https://github.com/pact-foundation/pact_broker-client#can-i-deploy) tool \(a CLI that allows you to easily check the verification status of your pacts\), to allow you to fully automate the CI/CD process for all the applications that use the Pact Broker, ensuring both sides of the contract are fulfilled before deploying.

See [Webhooks](http://docs.pact.io/pact_broker/advanced_topics/webhooks/) for more information on configuring Webhooks.

## Compatibility

-> This feature is available to both Pactflow and OSS users

## Example Usage

The following examples show the basic usage of the resource.

```hcl
resource "pact_webhook" "product_events" {
  description = "Trigger Product API verification build on contract changes for Admin UI"
  webhook_provider = {
    name = "ProductService"
  }
  webhook_consumer = {
    name = "AdminService"
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

  events = ["contract_published"]
  depends_on = [pact_pacticipant.AdminService, pact_pacticipant.ProductService]
}
```

### Write-only credentials (Terraform 1.11+)

Use `password_wo`, `username_wo` and/or `headers_wo` to send credentials without ever storing them in the Terraform plan or state. Because Terraform cannot detect changes to write-only values, bump the matching `*_wo_version` whenever a value changes:

```hcl
resource "pact_webhook" "product_events" {
  description = "Trigger Product API verification build on contract changes for Admin UI"
  request {
    url                 = "https://foo.com/some/endpoint"
    method              = "POST"
    username            = "test"
    password_wo         = var.webhook_password # can also be an ephemeral value
    password_wo_version = 1
    headers = {
      "Content-Type" = "application/json"
    }
    # merged with `headers`, but never stored in the plan or state
    headers_wo = jsonencode({
      "Authorization" = "Bearer ${var.webhook_token}"
    })
    headers_wo_version = 1
    body = jsonencode({ pact = "$${pactbroker.pactUrl}" })
  }

  events = ["contract_published"]
}
```

Write-only headers are not tracked by Terraform: drift on those headers (e.g. a header added outside Terraform) is not detected while `headers_wo` is in use.

## Argument Reference

The following arguments are supported:

- `description` - (Required, string) A human readable description of the Webhooks purpose.
- `webhook_provider` - (Optional, block) A provider to scope events to. See [Pacticipant](#pacticipant) below for details. Omitting the provider indicates the webhook should fire for all providers.

From https://docs.pact.io/pact_broker/advanced_topics/api_docs/webhooks#creating

> Both provider and consumer are optional - omitting either indicates that any pacticipant in that role will be matched.

- `webhook_consumer` - (Optional, block) A consumer to scope events to. See [Pacticipant](#pacticipant) below for details. Omitting the consumer indicates the webhook should fire for all consumers.
- `request` - (Required, block) The request to send when a webhook is fired. See [Request](#request) below for details.
- `events` - (Required, list of strings) one of `contract_requiring_verification_published`, `contract_content_changed`, `contract_published`, `provider_verification_published`, `provider_verification_succeeded` or `provider_verification_failed` (see [Webhooks](http://docs.pact.io/pact_broker/advanced_topics/webhooks/) for more on this).
- `team` - (Optional, string) The uuid of the team to assign to the webhook.

<a id="pacticipant"></a>

### Pacticipant

A pacticipant may be used as the consumer, provider, none or both in the webhook relationship.

- `name` - (Required, string) The name of the Pacticipant that should

<!-- start task-spec -->

<a id="request"></a>

### Request

`request` is a block within the configuration that can be repeated only **once** to specify the outgoing HTTP Request that should be sent for the Webhook.

- `url` (Required, string) A valid URL for the Webhook. This URL will be invoked on the configured events.
- `method` (Required, string) One of `POST`, `GET`, `PUT`, `PATCH`, or `DELETE`. Note that by default _only_ `POST` is supported. Other methods need to be explicitly opted in (this configuration is not currently supported by the provider)
- `username` (Optional, string) Basic auth username to send along with the request. Conflicts with `username_wo`.
- `username_wo` (Optional, string, write-only) Basic auth username to send along with the request, as a write-only argument (see `password_wo`). Requires Terraform 1.11+. Conflicts with `username`, must be used together with `username_wo_version`.
- `username_wo_version` (Optional, number) An arbitrary version number for `username_wo`. Change it (e.g. increment it) to update the username in the broker. Required when `username_wo` is set.
- `password` (Optional, string) Basic auth password to send along with the request. Stored (marked as sensitive) in the Terraform state. Conflicts with `password_wo`.
- `password_wo` (Optional, string, write-only) Basic auth password to send along with the request, as a [write-only argument](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments): it is never persisted in the Terraform plan or state, and can reference ephemeral values. Requires Terraform 1.11+. Conflicts with `password`, must be used together with `password_wo_version`.
- `password_wo_version` (Optional, number) An arbitrary version number for `password_wo`. Change it (e.g. increment it) to update the password in the broker. Required when `password_wo` is set.
- `headers` (Required, block) HTTP Headers as key/value pairs to send with the request.
- `headers_wo` (Optional, string, write-only) Additional HTTP headers to send with the request, as a JSON object of strings (use `jsonencode`). They are merged with `headers` (a header can't be set in both) and are never persisted in the Terraform plan or state. Requires Terraform 1.11+. Must be used together with `headers_wo_version`.
- `headers_wo_version` (Optional, number) An arbitrary version number for `headers_wo`. Change it (e.g. increment it) to update the write-only headers in the broker. Required when `headers_wo` is set.
- `body` (Required, string) A string body to be sent. JSON body validation will be checked and will produce a warning if invalid (it will _not_ fail validation).

## Outputs

- `uuid` - (string) The unique ID in Pactflow for this webhook.

## Importing

As per the [docs](https://www.terraform.io/docs/import/usage.html), the ID used for importing is the UUID of the webhook. You can obtain this through the API.

1. Create the shell for the user to be imported into:

```tf
resource "pact_webhook" "product_events" {
 ...
}
```

2. Import the resource

```sh
terraform  import pact_webhook.product_events ZBztO9l5poBdBDyUNewbNw
```

3. Plan any new changes

```sh
teraform plan
```
