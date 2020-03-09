package broker

type Secret struct {
	UUID        string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

// curl 'https://dius.pact.dius.com.au/secrets' '{"name":"foo","description":"bar","value":"baz"}' --compressed
type SecretResponseCreate struct {
	Secret
	HalDoc
}
