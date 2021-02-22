package broker

// AuthenticationSettings contains Pactflow tenant authentication settings
type AuthenticationSettings struct {
	Providers AuthenticationProviders `json:"authenticationProviders"`
}

// AuthenticationProviders for the current tenant
type AuthenticationProviders struct {
	Google GoogleAuthenticationSettings `json:"Google,omitempty"`
	Github GithubAuthenticationSettings `json:"GitHub,omitempty"`
}

// GoogleAuthenticationSettings configures the allowed email domains to authenticate to Pactflow
type GoogleAuthenticationSettings struct {
	EmailDomains []string `json:"EmailDomains"`
}

// GithubAuthenticationSettings configures the allowed organisations that may authenticate to Pactflow
// NOTE: this does not perform any Github OAuth process, which must be confirmed via the UI after enabling
type GithubAuthenticationSettings struct {
	Organizations []string `json:"GithubOrganizations"`
}

// TODO: this must be fixed on the Pactflow side

// AuthenticationSettingsResponse contains Pactflow tenant authentication settings
type AuthenticationSettingsResponse struct {
	Providers AuthenticationProviders `json:"authenticationProviders"`
}

// AuthenticationProvidersResponse for the current tenant
type AuthenticationProvidersResponse struct {
	Google GoogleAuthenticationSettings `json:"Google,omitempty"`
	Github GithubAuthenticationSettings `json:"GitHub,omitempty"`
}

// GoogleAuthenticationSettingsResponse configures the allowed email domains to authenticate to Pactflow
type GoogleAuthenticationSettingsResponse struct {
	EmailDomains []string `json:"emailDomains"`
}

// GithubAuthenticationSettingsResponse configures the allowed organisations that may authenticate to Pactflow
// NOTE: this does not perform any Github OAuth process, which must be confirmed via the UI after enabling
type GithubAuthenticationSettingsResponse struct {
	Organizations []string `json:"githubOrganizations"`
}
