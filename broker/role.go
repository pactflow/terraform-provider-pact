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

type PatchResource struct {
	Uuid string `json:"uuid"`
}
type UpdateResourceAssignmentPatchRequest struct {
	Operation string        `json:"op"`
	Path      string        `json:"path"`
	Value     PatchResource `json:"value"`
}

// Default platform roles
const (
	ROLE_ADMINISTRATOR      = "cf75d7c2-416b-11ea-af5e-53c3b1a4efd8"
	ROLE_CI_CD              = "c1878b8e-d09e-11ea-8fde-af02c4677eb7"
	ROLE_GUEST              = "d6938de2-e37c-11eb-b80e-3f68328092ca"
	ROLE_TEAM_ADMINISTRATOR = "d635f960-88f2-4f13-8043-4641a02dffa0"
	ROLE_USER               = "e9282e22-416b-11ea-a16e-57ee1bb61d18"
	ROLE_VIEWER             = "9fa50562-a42b-4771-aa8e-4bb3d623ae60"
)
