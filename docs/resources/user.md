# User resource

This resource manages the lifecycle of a user.

!> **This resource only works for Pactflow Cloud users, and is not compatible with the use of SSO providers (e.g. Google, SAML). That is, it will create users separate to any external Identity Provider you have configured**

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
* `email` - (Required for User, Optional for SystemAccount, string) The email address of the user to invite.
* `active` - (Optional, bool) Whether or not the user should be able to access the platform.
* `type` - (Optional, string) Whether or not to provision a standard user (`user`) or a System Account (`system`).

## Outputs

* `uuid` - (string) The unique ID in Pactflow for this user.

## Lifecycle

* `Create`: On an initial create, a user will be invited to Pactflow, and added to the local Pactflow account. If a user is not already in any Pactflow organisation, they will receive an email with a temporary token for them to reset their credentials.
* `Update`: Changes to the user will be applied as expected.
* `Delete`: Users will not be removed in the system, they will simply be disabled (Users are global in the Pactflow platform)

## Importing

As per the [docs](https://www.terraform.io/docs/import/usage.html), the ID used for importing is the UUID of the user. You can obtain this through the User API (`GET /admin/users`) and also through the user management screens.

1. Create the shell for the user to be imported into:

```tf
resource "pact_user" "someuser" {
  name = "Some User"
  active = true
  email = "foo@foo.com
}
```

2. Import the resource
```sh
terraform import pact_user.someuser e8d4891d-5c96-4dbf-b320-5bb7e3238269
```

3. Apply any new changes
```sh
teraform apply
```