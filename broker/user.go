package broker

// User represents a user login to Pactflow
type User struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
	Roles     []Role `json:"roles"`
}

// Users is a list of User objects to manage
type Users struct {
	Users []User `json:"users"`
}

// "users": [
//   {
//     "updatedAt": "2020-06-17T05:15:47+00:00",
//     "createdAt": "2020-06-17T14:55:39+00:00",
//     "uuid": "69e42b09-6228-4191-8505-508cee156b6d",
//     "identityProviderId": "8de836f4-ac35-4c83-acb2-9b25b881d133",
//     "name": "Matt Fellows",
//     "active": true,
//     "lastLogin": "2020-06-17T15:15:47.905+10:00",
//     "email": "matt.fellows@onegeek.com.au",
//     "roles": [
//       {
//         "createdAt": "2020-06-17T04:55:38+00:00",
//         "uuid": "cf75d7c2-416b-11ea-af5e-53c3b1a4efd8",
//         "name": "Administrator"
//       },
//       {
//         "createdAt": "2020-06-17T04:55:38+00:00",
//         "uuid": "e9282e22-416b-11ea-a16e-57ee1bb61d18",
//         "name": "User"
//       }
//     ],
//     "_links": {
//       "self": {
//         "title": "User",
//         "href": "http://localhost:9292/admin/users/69e42b09-6228-4191-8505-508cee156b6d"
//       }
//     }
//   }
// ]
