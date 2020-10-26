# Pacticipant Resource

This resource manages the lifecycle of a _Pacticipant_. A Pacticipant is an application that may perform the role of a consumer or a provider in the Pact ecosystem.

## Example Usage
The following examples show the basic usage of the resource.

```hcl
resource "pact_pacticipant" "admin" {
  name = "AdminService"
  repository_url = "github.com/company/admin"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the Pacticipant.
* `repository_url` - (Optional, string) A URL to the repository

## Importing

As per the [docs](https://www.terraform.io/docs/import/usage.html), the ID used for importingis simply the name of the Pacticipant.

1. Create the shell for the pacticipant to be imported into:

```tf
resource "pact_pacticipant" "Wiffle" {
  name = "Wiffle"
  repository_url = "github.com/company/admin"
}
```

2. Import the resource
```sh
terraform import pact_pacticipant.Wiffle Wiffle
```

3. Apply any new changes
```sh
teraform apply
```