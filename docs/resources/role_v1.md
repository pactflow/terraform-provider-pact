# Role v1 resource

!> This resource is deprecated, roles may now be added directly to a `pact_user` resource

This resource manages assigns existing platform roles to a given User.

See https://docs.pactflow.io/docs/user-interface/#settings---users for documentation on managing users and roles within Pactflow.

## Compatibility

-> This feature is only available for the Pactflow platform.

## Example Usage
The following examples show the basic usage of the resource. Here, we first create a user "billy" and attach a role to them, referencing the unique id (`uuid`) of the user.

```hcl
resource "pact_user" "billy" {
  name = "Billy Sampson"
  email = "billy@sampson.co"
  active = false
}

# Assign to the billy the user by referencing its uuid
resource "pact_role_v1" "billy_admin" {
  role = "administrator"
  user = pact_user.billy.uuid
}
```

## Argument Reference

The following arguments are supported:

* `role` - (Required, string) The string name of a role to assign. Currently the only option is `administrator`.
* `user` - (Required, string) The UUID of a user to apply the role to. Can refer to the `uuid` output of the User resource, or of a known ID in the system.

## Importing

This is not supported for Roles.