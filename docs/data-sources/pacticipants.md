# Pacticipants Data Source

This data source lists all _Pacticipants_ (applications) registered on the Pact Broker.

## Compatibility

-> This feature is available to both Pactflow and OSS users

## Example Usage

```hcl
data "pact_pacticipants" "all" {}

output "pacticipant_names" {
  value = [for p in data.pact_pacticipants.all.pacticipants : p.name]
}
```

## Attributes Reference

- `pacticipants` - A list of pacticipants, each with:
  - `name` - The name of the pacticipant.
  - `repository_url` - The URL or location of the VCS repository.
  - `main_branch` - The main (default) branch.
  - `display_name` - The display name of the pacticipant.
