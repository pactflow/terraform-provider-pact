package broker

// RoleV1 is a role that a User may have in Pactflow
// {
// 	"updatedAt": "2020-06-28T23:31:51+00:00",
// 	"createdAt": "2020-06-28T23:31:51+00:00",
// 	"uuid": "cf75d7c2-416b-11ea-af5e-53c3b1a4efd8",
// 	"name": "Administrator"
// }
type RoleV1 struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
