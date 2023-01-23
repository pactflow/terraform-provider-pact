package broker

// AuthenticationSettings contains PactFlow tenant authentication settings
type AuthenticationSettings struct {
	Providers AuthenticationProviders `json:"authenticationProviders"`
}

// AuthenticationProviders for the current tenant
type AuthenticationProviders struct {
	Google GoogleAuthenticationSettings `json:"Google,omitempty"`
	Github GithubAuthenticationSettings `json:"GitHub,omitempty"`
}

// GoogleAuthenticationSettings configures the allowed email domains to authenticate to PactFlow
type GoogleAuthenticationSettings struct {
	EmailDomains []string `json:"EmailDomains"`
}

// GithubAuthenticationSettings configures the allowed organisations that may authenticate to PactFlow
// NOTE: this does not perform any Github OAuth process, which must be confirmed via the UI after enabling
type GithubAuthenticationSettings struct {
	Organizations []string `json:"GithubOrganizations"`
}
