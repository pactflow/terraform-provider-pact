# Webhook Resource

This resource manages the lifecycle of a _Webhook_.

Webhooks allow you to trigger an HTTP request when a pact is changed, a pact is published, or a verification is published. The most common use case for webhooks is to trigger a provider build every time a pact changes, and to trigger a consumer build every time a verification is published.

Webhooks can be used in conjunction with the [can-i-deploy](https://github.com/pact-foundation/pact_broker-client#can-i-deploy) tool \(a CLI that allows you to easily check the verification status of your pacts\), to allow you to fully automate the CI/CD process for all the applications that use the Pact Broker, ensuring both sides of the contract are fulfilled before deploying.

See [Webhooks](http://docs.pact.io/pact_broker/advanced_topics/webhooks/) for more information on configuring Webhooks.

## Compatibility

-> This feature is available to both PactFlow and OSS users

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
- `username` (Optional, string) Basic auth username to send along with the request.
- `password` (Optional, string) Basic auth password to send along with the request.
- `headers` (Required, block) HTTP Headers as key/value pairs to send with the request.
- `body` (Required, string) A string body to be sent. JSON body validation will be checked and will produce a warning if invalid (it will _not_ fail validation).

## Outputs

- `uuid` - (string) The unique ID in PactFlow for this webhook.

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
