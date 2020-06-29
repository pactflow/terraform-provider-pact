# Secret Resource

This resource manages the lifecycle of a _Secret_. A Secret is an application that may perform the role of a consumer or a provider in the Pact ecosystem.

_NOTE_: this is currently only supported for the Pactflow.io platform.

## Example Usage
The following examples show the basic usage of the resource.

```hcl
resource "pact_secret" "some_jenkins_token" {
  name = "JenkinsToken"
  description = "A token for jenkins webhooks"
  value = "super secret thing"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the Secret (alphanumeric characters only)
* `description` - (Required, string) A human readable description of the Secret.
* `value` - (Required, string) The actual secret to store.

## Outputs

* `uuid` - (string) The unique ID in Pactflow for this secret.