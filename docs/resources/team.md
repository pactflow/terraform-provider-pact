# Team resource

This resource allows you to create teams for grouping users and applications.

https://docs.pactflow.io/docs/user-interface/settings/teams for documentation on managing teams.

## Compatibility

-> This feature is only available for the Pactflow platform.

## Example Usage

The following examples show the basic usage of the resource. We are creating a custom role that allows the user permissions to manage people and teams:

```hcl
resource "pact_team" "Futurama" {
  name = "Futurama"
  pacticipants = [
    pact_pacticipant.GraphQLAPI.name,
    pact_pacticipant.example.name
  ]
  users = [
    pact_user.bender_system_user.uuid,
  ]
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required, string) The name of the team.
- `pacticipants` - (Optional, list of strings) The set of UUIDs for each application to assign the team.
- `users` - (Optional, list of strings) The set of UUIDs for each user to assign to the team.
- `administrators` - (Optional, list of strings) The set of user UUIDs to assign as Admins to the team.

## Importing

As per the [docs](https://www.terraform.io/docs/import/usage.html), the ID used for importingis simply the name of the application.

You need to first obtain the existing role uuid, which you can obtain this through the Teams API (`GET /admin/teams`) or via the HAL browser.

1. Create the shell for the application to be imported into, ensuring the scopes are what you intend it to be:

```hcl
resource "pact_team" "Futurama" {
  name = "Futurama"
  pacticipants = [
    pact_pacticipant.GraphQLAPI.name,
    pact_pacticipant.example.name
  ]
  users = [
    pact_user.bender_system_user.uuid,
  ]
  administrators = [
    pact_user.leela_system_user.uuid,
  ]
}
```

2. Import the resource

```sh
terraform import pact_team.Futurama <team uuid>
```

3. Apply any new changes

```sh
teraform apply
```
