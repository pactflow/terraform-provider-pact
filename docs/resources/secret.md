# Secret Resource

This resource manages the lifecycle of a _Secret_. A Secret is an application that may perform the role of a consumer or a provider in the Pact ecosystem.

## Compatibility

-> This feature is only available for the Pactflow platform.

## Example Usage
The following examples show the basic usage of the resource.

```hcl
resource "pact_secret" "some_jenkins_token" {
  name = "JenkinsToken"
  description = "A token for jenkins webhooks"
  value = "super secret thing"
}
```

### Write-only secret value (Terraform 1.11+)

Use `value_wo` to send the secret to Pactflow without ever storing it in the Terraform plan or state. Because Terraform cannot detect changes to write-only values, bump `value_wo_version` whenever the secret changes:

```hcl
ephemeral "vault_kv_secret_v2" "jenkins" {
  mount = "secret"
  name  = "jenkins"
}

resource "pact_secret" "some_jenkins_token" {
  name             = "JenkinsToken"
  description      = "A token for jenkins webhooks"
  value_wo         = ephemeral.vault_kv_secret_v2.jenkins.data.token
  value_wo_version = 1
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required, string) The name of the Secret (alphanumeric characters only)
- `description` - (Required, string) A human readable description of the Secret.
- `value` - (Optional, string) The actual secret to store. Stored (marked as sensitive) in the Terraform state. Exactly one of `value` or `value_wo` must be set.
- `value_wo` - (Optional, string, write-only) The actual secret to store, as a [write-only argument](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments): it is never persisted in the Terraform plan or state, and can reference ephemeral values. Requires Terraform 1.11+. Must be used together with `value_wo_version`.
- `value_wo_version` - (Optional, number) An arbitrary version number for `value_wo`. Change it (e.g. increment it) to update the secret in Pactflow. Required when `value_wo` is set.
- `team` - (Optional, string) The uuid of the team to assign to the secret.

## Outputs

- `uuid` - (string) The unique ID in Pactflow for this secret.

## Importing

_NOTE_: secrets cannot be extracted through the API. Whilst a resource itself can be imported and then updated, the original value of the secret is not accessible via the API.

As per the [docs](https://www.terraform.io/docs/import/usage.html), the ID used for importing is the UUID of the secret. You can obtain this through the API.

1. Create the shell for the user to be imported into:

```tf
resource "pact_secret" "somesecret" {
  name = "SomeSecret"
  description = "Some Description"
  value_wo = var.some_secret
  value_wo_version = 1
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
