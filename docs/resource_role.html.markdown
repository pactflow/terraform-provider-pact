# Role resource

This resource manages assigns existing platform roles to a given User.

See https://docs.pactflow.io/docs/user-interface/#settings---users for documentation on managing users and roles within Pactflow.

## Example Usage
The following examples show the basic usage of the resource. Here, we first create a user "billy" and attach a role to them, referencing the unique id (`uuid`) of the user.

```hcl
resource "pact_user" "billy" {
  name = "Billy Sampson"
  email = "billy@sampson.co"
  active = false
}

# Assignt to the billy the user by referencing its uuid
resource "pact_role" "billy_admin" {
  role = "administrator"
  user = pact_user.billy.uuid
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the user.
* `email` - (Required, string) The email address of the user to invite.
* `active` - (Optional, bool) Whether or not the user should be able to access the platform.