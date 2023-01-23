package broker

// UserType describes the type of user represented in PactFlow
type UserType int

const (
	// RegularUser represents a human
	RegularUser UserType = iota

	// SystemAccount are designed for API access only and not for humans
	SystemAccount
)

// User represents a user login to PactFlow
type User struct {
	UUID               string   `json:"uuid,omitempty"`
	Name               string   `json:"name,omitempty"`
	FirstName          string   `json:"first_name,omitempty"`
	LastName           string   `json:"last_name,omitempty"`
	Email              string   `json:"email,omitempty"`
	Active             bool     `json:"active"`
	CreatedAt          string   `json:"createdAt,omitempty"`
	UpdatedAt          string   `json:"updatedAt,omitempty"`
	LastLogin          string   `json:"lastLogin,omitempty"`
	IdentityProviderID string   `json:"identityProviderId,omitempty"`
	Type               UserType `json:"type,omitempty"`
	TypeDescription    string   `json:"typeDescription,omitempty"`
	Embedded           struct {
		Roles []Role `json:"roles,omitempty"`
		Teams []Team `json:"teams,omitempty"`
	} `json:"_embedded,omitempty"`
}

// Users is a list of User objects to manage
type Users struct {
	Users []User `json:"users"`
}

// SetUserRolesRequest is used to set roles to a given user
type SetUserRolesRequest struct {
	Roles []string `json:"roles"`
}

// User:
//
// {
//   "updatedAt": "2020-10-20T06:53:11+00:00",
//   "createdAt": "2020-10-20T06:53:11+00:00",
//   "uuid": "4c260344-b170-41eb-b01e-c0ff10c72f25",
//   "identityProviderId": "d050ce6c-9fa7-4c71-80d3-95af9b3b5cb9",
//   "name": "Matt Fellows",
//   "active": true,
//   "lastLogin": "2020-10-20T09:05:07.623+00:00",
//   "email": "matt.fellows@onegeek.com.au",
//   "type": 0,
//   "typeDescription": "User",
//   "_embedded": {
//     "roles": [
//       {
//         "uuid": "e9282e22-416b-11ea-a16e-57ee1bb61d18",
//         "name": "Test Maintainer",
//         "permissions": [
//           {
//             "name": "Read system accounts",
//             "scope": "system_account:read:*"
//           },
//           {
//             "name": "Read roles",
//             "scope": "role:read:*"
//           },
//           {
//             "name": "Read users",
//             "scope": "user:read:*"
//           },
//           {
//             "name": "Read teams",
//             "scope": "team:read:*"
//           },
//           {
//             "name": "Manage API tokens",
//             "scope": "token:manage:own"
//           },
//           {
//             "name": "Manage contract data",
//             "scope": "contract_data:manage:*"
//           },
//           {
//             "name": "Read contract data",
//             "scope": "contract_data:read:*"
//           },
//           {
//             "name": "Manage webhooks",
//             "scope": "webhook:manage:*"
//           },
//           {
//             "name": "Manage secrets",
//             "scope": "secret:manage:*"
//           },
//           {
//             "name": "Bulk delete contract data",
//             "scope": "contract_data:bulk_delete:*"
//           }
//         ],
//         "_links": {
//           "self": {
//             "title": "User Role",
//             "name": "Test Maintainer",
//             "href": "https://testdemo.pactflow.io/admin/users/4c260344-b170-41eb-b01e-c0ff10c72f25/roles/e9282e22-416b-11ea-a16e-57ee1bb61d18"
//           }
//         }
//       },
//       {
//         "uuid": "cf75d7c2-416b-11ea-af5e-53c3b1a4efd8",
//         "name": "Administrator",
//         "permissions": [
//           {
//             "name": "Manage users",
//             "scope": "user:manage:*"
//           },
//           {
//             "name": "Manage teams",
//             "scope": "team:manage:*"
//           },
//           {
//             "name": "Invite users",
//             "scope": "user:invite"
//           },
//           {
//             "name": "Manage system accounts",
//             "scope": "system_account:manage:*"
//           },
//           {
//             "name": "Manage roles",
//             "scope": "role:manage:*"
//           },
//           {
//             "name": "Manage API tokens",
//             "scope": "token:manage:own"
//           },
//           {
//             "name": "Manage contract data",
//             "scope": "contract_data:manage:*"
//           },
//           {
//             "name": "Bulk delete contract data",
//             "scope": "contract_data:bulk_delete:*"
//           },
//           {
//             "name": "Manage webhooks",
//             "scope": "webhook:manage:*"
//           },
//           {
//             "name": "Manage secrets",
//             "scope": "secret:manage:*"
//           }
//         ],
//         "_links": {
//           "self": {
//             "title": "User Role",
//             "name": "Administrator",
//             "href": "https://testdemo.pactflow.io/admin/users/4c260344-b170-41eb-b01e-c0ff10c72f25/roles/cf75d7c2-416b-11ea-af5e-53c3b1a4efd8"
//           }
//         }
//       }
//     ],
//     "teams": [
//       {
//         "name": "PactFlow Demos",
//         "uuid": "6d746ad5-919f-49e3-84c0-648cafc5d912",
//         "_links": {
//           "self": {
//             "title": "Team",
//             "name": "PactFlow Demos",
//             "href": "https://testdemo.pactflow.io/teams/6d746ad5-919f-49e3-84c0-648cafc5d912"
//           },
//           "pf:membership": {
//             "title": "Team membership",
//             "href": "https://testdemo.pactflow.io/teams/6d746ad5-919f-49e3-84c0-648cafc5d912/users/4c260344-b170-41eb-b01e-c0ff10c72f25"
//           },
//           "pf:integrations": {
//             "title": "Team integrations",
//             "name": "PactFlow Demos",
//             "href": "https://testdemo.pactflow.io/integrations/team/6d746ad5-919f-49e3-84c0-648cafc5d912"
//           }
//         }
//       }
//     ]
//   },
//   "_links": {
//     "self": {
//       "title": "User",
//       "href": "https://testdemo.pactflow.io/admin/users/4c260344-b170-41eb-b01e-c0ff10c72f25"
//     }
//   },
//   "roles": [
//     {
//       "updatedAt": "2020-04-02T22:42:28+00:00",
//       "createdAt": "2020-04-02T22:42:28+00:00",
//       "uuid": "cf75d7c2-416b-11ea-af5e-53c3b1a4efd8",
//       "name": "Administrator",
//       "permissions": [
//         {
//           "name": "Manage users",
//           "scope": "user:manage:*"
//         },
//         {
//           "name": "Manage teams",
//           "scope": "team:manage:*"
//         },
//         {
//           "name": "Invite users",
//           "scope": "user:invite"
//         },
//         {
//           "name": "Manage system accounts",
//           "scope": "system_account:manage:*"
//         },
//         {
//           "name": "Manage roles",
//           "scope": "role:manage:*"
//         },
//         {
//           "name": "Manage API tokens",
//           "scope": "token:manage:own"
//         },
//         {
//           "name": "Manage contract data",
//           "scope": "contract_data:manage:*"
//         },
//         {
//           "name": "Bulk delete contract data",
//           "scope": "contract_data:bulk_delete:*"
//         },
//         {
//           "name": "Manage webhooks",
//           "scope": "webhook:manage:*"
//         },
//         {
//           "name": "Manage secrets",
//           "scope": "secret:manage:*"
//         }
//       ]
//     },
//     {
//       "updatedAt": "2020-04-02T22:42:28+00:00",
//       "createdAt": "2020-04-02T22:42:28+00:00",
//       "uuid": "e9282e22-416b-11ea-a16e-57ee1bb61d18",
//       "name": "Test Maintainer",
//       "permissions": [
//         {
//           "name": "Read system accounts",
//           "scope": "system_account:read:*"
//         },
//         {
//           "name": "Read roles",
//           "scope": "role:read:*"
//         },
//         {
//           "name": "Read users",
//           "scope": "user:read:*"
//         },
//         {
//           "name": "Read teams",
//           "scope": "team:read:*"
//         },
//         {
//           "name": "Manage API tokens",
//           "scope": "token:manage:own"
//         },
//         {
//           "name": "Manage contract data",
//           "scope": "contract_data:manage:*"
//         },
//         {
//           "name": "Read contract data",
//           "scope": "contract_data:read:*"
//         },
//         {
//           "name": "Manage webhooks",
//           "scope": "webhook:manage:*"
//         },
//         {
//           "name": "Manage secrets",
//           "scope": "secret:manage:*"
//         },
//         {
//           "name": "Bulk delete contract data",
//           "scope": "contract_data:bulk_delete:*"
//         }
//       ]
//     }
//   ]
// }

// curl -X PUT localhost:9292/admin/users/04c60c21-4b0b-45ea-b38f-ca1d8211323a/roles -H"Authorization: Bearer localhost" -d '{ "roles": [ "8bff8a01-0993-41b5-97c5-14dc9ce23268" ] }' | jq .
