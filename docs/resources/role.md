# Role resource

This resource allows you to create custom Roles for assigning to Users.

See https://docs.pactflow.io/docs/permissions/predefined-roles for documentation on managing users and roles within Pactflow.

## Example Usage
The following example shows the basic usage of the resource. We are creating a custom role that allows the user permissions to manage people and teams:

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
* `scopes` - (Required, list of strings) The scopes to apply to the role (see below for the available scopes)

## Available scopes

At the time of authoring, the following permission scopes are available. See https://docs.pactflow.io/docs/permissions/permissions for the definitive list.

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

As per the [docs](https://www.terraform.io/docs/import/usage.html), the ID used for importingis simply the name of the application.

You need to first obtain the existing role uuid, which you can find via the API/HAL browser.

1. Create the shell for the application to be imported into, ensuring the scopes are what you intend it to be:


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

2. Import the resource
```sh
terraform import pact_role.special_role <role uuid>
```

3. Apply any new changes
```sh
teraform apply
```