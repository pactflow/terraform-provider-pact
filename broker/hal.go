package broker

type Headers map[string]string

// Link represents a link to a resource
type Link struct {
	Href  string `json:"href"`
	Title string `json:"title"`
	Name  string `json:"name"`
}

// HalLinks represents the _links key in a HAL document.
type HalLinks map[string]Link

// HalDoc is a simple representation of the HAL response from a Pact Broker.
type HalDoc struct {
	Links HalLinks `json:"_links"`
}
