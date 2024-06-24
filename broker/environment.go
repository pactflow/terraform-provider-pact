package broker

// Environment represents an environment and application may be deployed to
type Environment struct {
	Name        string                   `json:"name,omitempty"`
	Production  bool                     `json:"production"`
	DisplayName string                   `json:"displayName,omitempty"`
	CreatedAt   string                   `json:"createdAt,omitempty"`
	UpdatedAt   string                   `json:"updatedAt,omitempty"`
	UUID        string                   `json:"uuid,omitempty"`
	Embedded    EnvironmentEmbeddedItems `json:"_embedded,omitempty"`
}

type EnvironmentCreateOrUpdateRequest struct {
	UUID        string   `json:"-"`
	DisplayName string   `json:"displayName,omitempty"`
	Name        string   `json:"name,omitempty"`
	Production  bool     `json:"production"`
	Teams       []string `json:"teamUuids"`
}

type EnvironmentCreateOrUpdateResponse struct {
	Name        string                   `json:"name,omitempty"`
	Production  bool                     `json:"production"`
	DisplayName string                   `json:"displayName,omitempty"`
	CreatedAt   string                   `json:"createdAt,omitempty"`
	UpdatedAt   string                   `json:"updatedAt,omitempty"`
	UUID        string                   `json:"uuid,omitempty"`
	Teams       []string                 `json:"teamUuids"`
	Embedded    EnvironmentEmbeddedItems `json:"_embedded,omitempty"`
}

type EnvironmentEmbeddedItems struct {
	Teams []EnvironmentEmbeddedTeams `json:"teams,omitempty"`
}

type EnvironmentEmbeddedTeams struct {
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
}

// POST environments
//  {"uuid":"2739c79b-a6ba-4398-be7a-85ec96f79fbe","name":"test1","displayName":"test1 with teams","production":false,"createdAt":"2022-03-07T12:22:05+00:00","teamUuids":["6d746ad5-919f-49e3-84c0-648cafc5d912"],"_embedded":{"teams":[{"uuid":"6d746ad5-919f-49e3-84c0-648cafc5d912","name":"Pactflow Demos","_links":{"self":{"title":"Team","href":"https://testdemo.pactflow.io/admin/teams/6d746ad5-919f-49e3-84c0-648cafc5d912"}}}]},"_links":{"self":{"title":"Environment","name":"test1","href":"https://testdemo.pactflow.io/environments/2739c79b-a6ba-4398-be7a-85ec96f79fbe"},"pb:currently-deployed-deployed-versions":{"title":"Versions currently deployed to test1 with teams environment","href":"https://testdemo.pactflow.io/environments/2739c79b-a6ba-4398-be7a-85ec96f79fbe/deployed-versions/currently-deployed"},"pb:currently-supported-released-versions":{"title":"Versions released and supported in test1 with teams environment","href":"https://testdemo.pactflow.io/environments/2739c79b-a6ba-4398-be7a-85ec96f79fbe/released-versions/currently-supported"},"pb:environments":{"title":"Environments","href":"https://testdemo.pactflow.io/environments"}}}

// TODO: why does the GET not return `teamUuids` also?

// GET environments/:env

// {
// 	"uuid": "e188cc63-7e49-4243-8fba-a8e5437d7d5c",
// 	"name": "aoeuaoeua",
// 	"displayName": "aoeuaoeu",
// 	"production": true,
// 	"updatedAt": "2022-03-02T09:29:02+00:00",
// 	"createdAt": "2022-03-02T09:28:34+00:00",
// 	"_embedded": {
// 			"teams": [
// 					{
// 							"uuid": "219af87b-7a6f-4efb-89dd-9cdd79a050d1",
// 							"name": "O'reilly superstream",
// 							"_links": {
// 									"self": {
// 											"title": "Team",
// 											"href": "https://testdemo.pactflow.io/admin/teams/219af87b-7a6f-4efb-89dd-9cdd79a050d1"
// 									}
// 							}
// 					}
// 			]
// 	},
// 	"_links": {
// 			"self": {
// 					"title": "Environment",
// 					"name": "aoeuaoeua",
// 					"href": "https://testdemo.pactflow.io/environments/e188cc63-7e49-4243-8fba-a8e5437d7d5c"
// 			},
// 			"pb:currently-deployed-deployed-versions": {
// 					"title": "Versions currently deployed to aoeuaoeu environment",
// 					"href": "https://testdemo.pactflow.io/environments/e188cc63-7e49-4243-8fba-a8e5437d7d5c/deployed-versions/currently-deployed"
// 			},
// 			"pb:currently-supported-released-versions": {
// 					"title": "Versions released and supported in aoeuaoeu environment",
// 					"href": "https://testdemo.pactflow.io/environments/e188cc63-7e49-4243-8fba-a8e5437d7d5c/released-versions/currently-supported"
// 			},
// 			"pb:environments": {
// 					"title": "Environments",
// 					"href": "https://testdemo.pactflow.io/environments"
// 			}
// 	}
// }
