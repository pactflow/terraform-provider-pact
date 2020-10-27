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
resource "pact_role_v1" "billy_admin" {
  role = "administrator"
  user = pact_user.billy.uuid
}
```

## Argument Reference

The following arguments are supported:

* `role` - (Required, string) The string name of a role to assign. Currently the only option is `administrator`.
* `user` - (Required, string) The UUID of a user to apply the role to. Can refer to the `uuid` output of the User resource, or of a known ID in the system.

## Default Roles

The following roles are available by default:
* Administrator: `cf75d7c2-416b-11ea-af5e-53c3b1a4efd8`
* test maintainer: `e9282e22-416b-11ea-a16e-57ee1bb61d18`
* CI/CD: `c1878b8e-d09e-11ea-8fde-af02c4677eb7`
* Viewer: `9fa50562-a42b-4771-aa8e-4bb3d623ae60`

## Importing

TBC