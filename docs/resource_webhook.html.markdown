# Webhook Resource

This resource manages the lifecycle of a _Webhook__.

Webhooks allow you to trigger an HTTP request when a pact is changed, a pact is published, or a verification is published. The most common use case for webhooks is to trigger a provider build every time a pact changes, and to trigger a consumer build every time a verification is published.

Webhooks can be used in conjunction with the [can-i-deploy](https://github.com/pact-foundation/pact_broker-client#can-i-deploy) tool \(a CLI that allows you to easily check the verification status of your pacts\), to allow you to fully automate the CI/CD process for all the applications that use the Pact Broker, ensuring both sides of the contract are fulfilled before deploying.

### The 'contract content changed' event

The broker uses the following logic to determine if a pact has changed:

* If the relevant consumer version has any tags, then for each tag, check if the content is different from the previous latest version for that tag. It is 'changed' if any of the checks are true. One side effect of this is that brand new tags will trigger a pact changed event, even if the content is the same as a previous version.
* If the relevant consumer version has no tags, then check if the content has changed since the previous latest version.

### The 'contract published' event

This is triggered every time a pact is published, regardless of whether it has changed or not.

### The 'verification published' event.

This is triggered every time a verification is published.

See [Webhooks](http://docs.pact.io/pact_broker/advanced_topics/webhooks/) for more information on configuring Webhooks.

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

  events = ["contract_content_changed", "contract_published"]
  depends_on = [pact_pacticipant.AdminService, pact_pacticipant.ProductService]
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required, string) A human readable description of the Webhooks purpose.
* `webhook_provider` - (Optional, block) A provider to scope events to. See [Pacticipant](#pacticipant) below for details.
* `webhook_consumer` - (Optional, block) A consumer to scope events to. See [Pacticipant](#pacticipant) below for details.
* `request` - (Required, block) The request to send when a webhook is fired. See [Request](#request) below for details.
* `events` - (Required, list of strings) one of	`"contract_content_changed"`, `"contract_published"` or `"provider_verification_published"` (see [Webhooks](http://docs.pact.io/pact_broker/advanced_topics/webhooks/) for more on this).


<a id="pacticipant"></a>
### Pacticipant

A pacticipant may be used as the consumer, provider, none or both in the webhook relationship.

* `name` - (Required, string) The name of the Pacticipant that should

<!-- start task-spec -->
<a id="request"></a>
### Request

`request` is a block within the configuration that can be repeated only **once** to specify the outgoing HTTP Request that should be sent for the Webhook.

* `url` (Required, string) A valid URL for the Webhook. This URL will be invoked on the configured events.
* `method` (Required, string) One of `POST`, `GET`, `PUT`, `PATCH`, or `DELETE`. Note that by default _only_ `POST` is supported. Other methods need to be explicitly opted in (this configuration is not currently supported by the provider)
* `username` (Optional, string) Basic auth username to send along with the request.
* `password` (Optional, string) Basic auth password to send along with the request.
* `headers` (Required, block) HTTP Headers as key/value pairs to send with the request.
* `body` (Required, string) A string body to be sent. JSON body validation will be checked and will produce a warning if invalid (it will _not_ fail validation).