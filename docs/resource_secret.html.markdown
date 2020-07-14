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

## Importing

_NOTE_: secrets cannot be extracted through the API. Whilst a resource itself can be imported and then updated, the original value of the secret is not accessible via the API.

As per the [docs](https://www.terraform.io/docs/import/usage.html), the ID used for importing is the UUID of the secret. You can obtain this through the API.

1. Create the shell for the user to be imported into:

```tf
resource "pact_secret" "somesecret" {
  name = "SomeSecret"
  description = "Some Description"
}
```

2. Import the resource
```sh
terraform import pact_secret.somesecret e8d4891d-5c96-4dbf-b320-5bb7e3238269
```

3. Apply any new changes
```sh
teraform apply
```