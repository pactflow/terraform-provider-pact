package broker

// APIToken represents an individual API token
type APIToken struct {
	UUID        string `json:"uuid"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

// APITokensEmbedded contains the embedded links in the resource
type APITokensEmbedded struct {
	Items []APIToken `json:"items"`
}

// APITokenResponse is the response body for any CRU API calls
type APITokenResponse struct {
	APIToken
	HalDoc
}

// APITokensResponse is the response body for List API calls
type APITokensResponse struct {
	Embedded APITokensEmbedded `json:"_embedded"`
	HalDoc
}
