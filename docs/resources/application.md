# Pacticipant Resource

This resource manages the lifecycle of a _Application_ (also known as a pacticipant). An Application may perform the role of a consumer or a provider in the Pact ecosystem.

## Compatibility

-> This feature is available to both Pactflow and OSS users

## Example Usage

The following examples show the basic usage of the resource.

```hcl
resource "pact_application" "admin" {
  name = "AdminService"
  repository_url = "github.com/company/admin"
}
```

-> Thisresource was renamed in order to simplify our language of core Pact concepts in the ecosystem. It is still available under its previous name `pact_pacticipant` for compatibility reasons.

## Argument Reference

The following arguments are supported:

- `name` - (Required, string) The name of the application.
- `repository_url` - (Optional, string) A URL to the repository
- `main_branch` - (Optional, string) The name of the main branch

## Importing

As per the [docs](https://www.terraform.io/docs/import/usage.html), the ID used for importingis simply the name of the application.

1. Create the shell for the application to be imported into:

```tf
resource "pact_application" "Wiffle" {
  name = "Wiffle"
  repository_url = "github.com/company/admin"
  main_branch = "main"
}
```

2. Import the resource

```sh
terraform import pact_application.Wiffle Wiffle
```

3. Apply any new changes

```sh
teraform apply
```
