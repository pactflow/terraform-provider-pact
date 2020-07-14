# Token Resource

This resource manages the lifecycle of an _API Token_. A Token can be used to make API calls to the Pactflow platform.

**It is highly recommended that this resource only be used to import existing tokens, and not be used to update existing tokens**

_NOTE_: this is currently only supported for the Pactflow.io platform.



## Example Usage

The following examples show the basic usage of the resource.

```hcl
resource "pact_token" "read_only_api_token" {
  type = "read-only"
  name = "Local dev token"
}
resource "pact_token" "read_write_api_token" {
  type = "read-write"
  name = "CI token"
}
```

**NOTE**: There can be at most 1 of each type of token, as shown above. Our [roadamp](https://pactflow.io/pactflow-feature-roadmap/) includes expanded support for API tokens (multiple named tokens at the user and administration level).

**NOTE**: If you change the `read-write` token, it will generate a new token and invalidate the existing token. You will need to use the new value returned to run Terraform again. For example, you may want to extract the `value` property using the Terraform [Output](https://www.terraform.io/docs/configuration/outputs.html) feature.

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the token (for documentation purposes, this won't affect anything in the UI). Changing the name will generate a new token.
* `type` - (Required, string) One of 'read-only' or 'read-write'. Read only tokens are not allowed to modify any resources, whilst write tokens are able to modify any resource.

## Outputs

* `uuid` (string) The UUID of the token for use in API calls.
* `description` (string) The description of the token.
* `value` (sensitive, string) The actual API token for use in authenticated calls.

## Lifecycle

* `Create`: On an initial create, the plugin simply fetches the remote token (i.e. it imports the value, but does not change it)
* `Update`: Changing the `name` property will regenerate the token, resulting in a new local and remote value.
* `Delete`: API tokens (currently) cannot be deleted. This operation simply detaches the local state from the remote broker.

## Importing

As per the [docs](https://www.terraform.io/docs/import/usage.html), the ID used for importing is the UUID of the token.. You can obtain this through the API.

1. Create the shell for the user to be imported into:

```tf
resource "pact_token" "read_only_api_token" {
  type = "read-only"
  name = "Local dev token"
}
```

2. Import the resource
```sh
terraform import pact_token.read_only_api_token j3xYRnn9dgkkSWrXB1oaXw
```

The token value has now been extracted and can now be managed by terraform.