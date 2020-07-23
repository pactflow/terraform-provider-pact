package broker

// WebhookEvent represents the types of events that trigger a Webhook
type WebhookEvent struct {
	Name string `json:"name"`
}

// Request is an HTTP request structure
type Request struct {
	Method   string  `json:"method"`
	URL      string  `json:"url"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Headers  Headers `json:"headers"`
	Body     string  `json:"body"`
}

// Webhook represents a webhook configured in the broker
type Webhook struct {
	ID          string         `json:"-"`
	Description string         `json:"description"`
	Enabled     bool           `json:"enabled"`
	CreatedAt   string         `json:"createdAt,omitempty"`
	Provider    *Pacticipant   `json:"provider"`
	Consumer    *Pacticipant   `json:"consumer"`
	Events      []WebhookEvent `json:"events,omitempty"`
	Request     Request        `json:"request"`
}

// WebhookResponse is the response body for any CRU methods
type WebhookResponse struct {
	Webhook
	HalDoc
}
