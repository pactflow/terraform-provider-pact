# Authentication Settings Resource

This resource manages the authentication settings on a Pactflow account.

-> This is currently only supported for the Pactflow.io platform, and does not apply to the on-premise version

## Example Usage
The following examples show the basic usage of the resource.

```hcl
resource "pact_authentication" "authentication" {
  github_organizations = ["DiUS", "pactflow"]
  google_domains = ["dius.com.au", "onegeek.com.au"]
}
```

## Argument Reference

The following arguments are supported:

* `github_organizations` - (Optional, list of strings) The Github organisations allowed access to the account
* `google_domains` - (Optional, list of strings) The list of Google domains (e.g. foo.com) allowed access to the account

## Importing

Import is not supported, as it's not useful. Simply copy the settings from the UI into the resource and you should be able to apply the settings over the top (the same as an import, except without needing to first perform the import step).