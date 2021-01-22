# Role resource

This resource allows you to create custom Roles for assigning to users..

See https://docs.pactflow.io/docs/permissions/predefined-roles for documentation on managing users and roles within Pactflow.

## Example Usage
The following examples show the basic usage of the resource. We are creating a custom role that allows the user permissions to manage people and teams:

```hcl
resource "pact_role" "special_role" {
  name = "CustomUserManagementRole"
  scopes = [
    "user:manage:*",
    "team:manage:*",
    "user:invite",
    "system_account:manage:*",
    "system_account:read:*",
    "user:read:*",
    "team:read:*",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `role` - (Required, string) The string name of a role to assign. Currently the only option is `administrator`.
* `user` - (Required, string) The UUID of a user to apply the role to. Can refer to the `uuid` output of the User resource, or of a known ID in the system.

## Available permissions

At the time of authoring, the following permissions are available. See https://docs.pactflow.io/docs/permissions/permissions for the definitive list.

* `user:manage:*`
* `team:manage:*`
* `user:invite`
* `system_account:manage:*`
* `system_account:read:*`
* `user:read:*`
* `team:read:*`
* `contract_data:manage:*`
* `contract_data:read:*`
* `contract_data:bulk_delete:*`
* `webhook:manage:*`
* `secret:manage:*`
* `role:manage:*`
* `role:read:*`
* `token:manage:own`
* `read_token:manage:ow`

## Importing

TBC