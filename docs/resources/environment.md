# Environment resource

This resource allows you to create custom Environments for assigning to Users.

See https://docs.pact.io/pact_broker/recording_deployments_and_releases#environments and https://docs.pactflow.io/docs/user-interface/settings/environments for further documentation on managing environments.

## Example Usage

The following example shows the basic usage of the resource. We are creating a custom environment that allows the user permissions to manage people and teams:

```hcl
resource "pact_environment" "UAT" {
  name = "UAT"
  display_name = "User Acceptance Testing"
  production = false
  teams = ["4ac05ed8-9e3b-4159-96c0-ad19e3b93658"]
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required, string) The string name of a environment to create. This must be contain only alphanumeric strings.
- `display_name` - (Required, string) The visible display name of the environment
- `production` - (Required, boolean) Whether or not the environment is a "production" environment or not
- `team_uuids` - (Optional, list of strings) The list of teams to assign to the team

## Importing

As per the [docs](https://www.terraform.io/docs/import/usage.html), the ID used for importingis simply the name of the application.

You need to first obtain the existing environment uuid, which you can find via the API/HAL browser.

1. Create the shell for the application to be imported into, ensuring the scopes are what you intend it to be:

```hcl
resource "pact_environment" "UAT" {
  name = "UAT"
  display_name = "User Acceptance Testing"
  production = false
  teams = ["4ac05ed8-9e3b-4159-96c0-ad19e3b93658"]
}
```

2. Import the resource

```sh
terraform import pact_environment.production <environment uuid>
```

3. Apply any new changes

```sh
teraform apply
```
