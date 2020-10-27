package broker

type Role struct {
	UUID        string       `json:"uuid,omitempty"`
	Name        string       `json:"name,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
}

type Permission struct {
	Name        string `json:"name,omitempty"`
	Scope       string `json:"scope,omitempty"`
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
}

var AllowedScopes = []string{
	"user:manage:*",
	"team:manage:*",
	"user:invite",
	"system_account:manage:*",
	"system_account:read:*",
	"user:read:*",
	"team:read:*",
	"contract_data:manage:*",
	"contract_data:read:*",
	"contract_data:bulk_delete:*",
	"webhook:manage:*",
	"secret:manage:*",
	"role:manage:*",
	"role:read:*",
	"token:manage:own",
	"read_token:manage:own",
}
