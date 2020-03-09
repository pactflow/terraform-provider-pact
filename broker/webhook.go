package broker

type WebhookEvent struct {
	Name string `json:"name"`
}

type RequestBody map[string]interface{}

type Request struct {
	Method   string      `json:"method"`
	URL      string      `json:"url"`
	Username string      `json:"username"`
	Password string      `json:"password"`
	Headers  Headers     `json:"headers"`
	Body     RequestBody `json:"body"`
}

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

type WebhookResponseCreate struct {
	Webhook
	HalDoc
}
