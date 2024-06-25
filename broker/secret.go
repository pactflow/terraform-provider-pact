package broker

type Secret struct {
	UUID        string `json:"-"`
	TeamUUID    string `json:"teamUuid,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value,omitempty"`
}

// curl 'https://dius.pact.dius.com.au/secrets' '{"name":"foo","description":"bar","value":"baz"}' --compressed
type SecretResponse struct {
	Secret
	HalDoc
}
