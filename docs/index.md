# Pact Provider

Pact Broker Provider sets up a connection to a hosted Broker.

-> We currently support both the Open Source Pact Broker and our Pactflow.io platform.

## Example Usage
The following examples show the basic usage of the resouproviderrce.

```hcl
provider "pact" {
  # host = "http://localhost"
  # basic_auth_username = "pact_broker"
  # basic_auth_password = "pact_broker"
  host = "https://dius.pact.dius.com.au"
  access_token = "oO_ITO-bummTj6_oJoMPmw"
  tls_insecure = true
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required, string) A fully qualified hostname (e.g. for a Pactflow account https://mybroker.pact.dius.com.au
* `basic_auth_username` - (Optional, string) A basic auth username to authenticate to a Pact Broker (not required for Pactflow users)
* `basic_auth_password` - (Optional, string) A basic auth password to authenticate to a Pact Broker (not required for Pactflow users)
* `access_token` - (Optional, string) An API Bearer token to authenticate to a Pactflow account (for Pactflow users only)
* `tls_insecure` - (Optional, bool) Disable TLS verification checks (useful for internal brokers with self-signed certificates)