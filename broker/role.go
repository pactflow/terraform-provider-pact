package broker

type Role struct {
	UUID        string       `json:"uuid"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions,omitempty"`
}

type Permission struct {
	Name  string `json:"name"`
	Scope string `json:"scope"`
}
