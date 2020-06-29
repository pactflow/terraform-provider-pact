# User resource

This resource manages the lifecycle of a user.

See https://docs.pactflow.io/docs/user-interface/#settings---users for documentation on managing users within Pactflow.

## Example Usage
The following examples show the basic usage of the resource.

```hcl
resource "pact_user" "billy" {
  name = "Billy Sampson"
  email = "billy@sampson.co"
  active = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the user.
* `email` - (Required, string) The email address of the user to invite.
* `active` - (Optional, bool) Whether or not the user should be able to access the platform.

## Outputs

* `uuid` - (string) The unique ID in Pactflow for this user.

## Lifecycle

* `Create`: On an initial create, a user will be created in the local Pactflow account. If a user is not in any Pactflow organisation, they will receive an email with a temporary token for them to reset their credentials.
* `Update`: Changes to the user will be applied as expected.
* `Delete`: Users will not be removed in the system, they will simply be disabled.