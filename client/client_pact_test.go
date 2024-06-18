package client

import (
	"testing"

	"fmt"
	"net/url"

	"github.com/mitchellh/copystructure"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/log"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/pactflow/terraform/broker"
	"github.com/stretchr/testify/assert"
)

func TestClientPact(t *testing.T) {
	assert.Equal(t, true, true)
}

func TestTerraformClientPact(t *testing.T) {
	log.SetLogLevel("ERROR")

	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "terraform-client",
		Provider: "pactflow-application-saas",
		Host:     "127.0.0.1",
	})
	assert.NoError(t, err)

	pacticipant := broker.Pacticipant{
		Name:          "terraform-client",
		RepositoryURL: "https://github.com/pactflow/new-terraform-provider-pact",
		MainBranch:    "Main",
		DisplayName:   "Terraform Client",
	}

	t.Run("Pacticipant", func(t *testing.T) {
		// pacticipant := broker.Pacticipant{
		// 	Name:          "terraform-client",
		// 	RepositoryURL: "https://github.com/pactflow/terraform-provider-pact",
		// 	MainBranch:    "Main",
		// 	DisplayName:   "Terraform Client",
		// }

		t.Run("CreatePacticipant", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				UponReceiving("a request to create a pacticipant").
				WithRequest("POST", "/pacticipants", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(pacticipant))
				}).
				WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(pacticipant))
				})

			err := mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.CreatePacticipant(pacticipant)
				assert.NoError(t, e)
				assert.Equal(t, "terraform-client", res.Name)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("ReadPacticipant", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a pacticipant with name terraform-client exists").
				UponReceiving("a request to get a pacticipant").
				WithRequest("GET", "/pacticipants/terraform-client", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.BodyMatch(&broker.Pacticipant{})
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.ReadPacticipant("terraform-client")
				assert.NoError(t, e)
				assert.Equal(t, "terraform-client", res.Name)
				assert.NotEmpty(t, res.RepositoryURL)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("UpdatePacticipant", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a pacticipant with name terraform-client exists").
				UponReceiving("a request to update a pacticipant").
				WithRequest("PATCH", "/pacticipants/terraform-client", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(pacticipant))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(pacticipant))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.UpdatePacticipant(pacticipant)
				assert.NoError(t, e)
				assert.Equal(t, "terraform-client", res.Name)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("DeletePacticipant", func(t *testing.T) {
			newPacticipant := broker.Pacticipant{
				Name:          "terraform-client",
				RepositoryURL: "https://github.com/pactflow/new-terraform-provider-pact",
			}

			mockProvider.
				AddInteraction().
				Given("a pacticipant with name terraform-client exists").
				UponReceiving("a request to delete a pacticipant").
				WithRequest("DELETE", "/pacticipants/terraform-client", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
				}).
				WillRespondWith(204)

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				return client.DeletePacticipant(newPacticipant)
			})
			assert.NoError(t, err)
		})
	})

	t.Run("Team", func(t *testing.T) {
		create := broker.TeamCreateOrUpdateRequest{
			Name:             "terraform-team",
			PacticipantNames: []string{pacticipant.Name},
		}

		created := broker.Team{
			Name: create.Name,
			UUID: "99643109-adb0-4e68-b25f-7b14d6bcae16",
		}

		team := broker.Team{
			Name: create.Name,
			UUID: created.UUID,
			Embedded: broker.TeamEmbeddedItems{
				Pacticipants: []broker.Pacticipant{
					{
						Name: "Pactflow Saas",
					},
				},
				Members: []broker.TeamUser{
					{
						UUID:   "4c260344-b170-41eb-b01e-c0ff10c72f25",
						Active: true,
					},
				},
				Administrators: []broker.TeamUser{
					{
						UUID: "4c260344-b170-41eb-b01e-c0ff10c72f25",
					},
				},
				Environments: []broker.TeamEnvironment{
					{
						UUID: "8000883c-abf0-4b4c-b993-426f607092a9",
					},
				},
			},
		}
		
		updated := broker.Team{
			Name: create.Name,
			UUID: created.UUID,
			Embedded: broker.TeamEmbeddedItems{
				Pacticipants: []broker.Pacticipant{
					{
						Name: "Terraform Client",
					},
				},
				Administrators: []broker.TeamUser{
					{
						UUID: "4c260344-b170-41eb-b01e-c0ff10c72f25",
					},
				},
				Environments: []broker.TeamEnvironment{
					{
						UUID: "8000883c-abf0-4b4c-b993-426f607092a9",
					},
				},
			},
		}

		update := broker.TeamCreateOrUpdateRequest{
			UUID:               updated.UUID,
			Name:               updated.Name,
			PacticipantNames:   []string{pacticipant.Name},
			AdministratorUUIDs: []string{"4c260344-b170-41eb-b01e-c0ff10c72f25"},
			EnvironmentUUIDs:   []string{"8000883c-abf0-4b4c-b993-426f607092a9"},
		}

		t.Run("ReadTeam", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a team with uuid 99643109-adb0-4e68-b25f-7b14d6bcae16 exists").
				UponReceiving("a request to get a team").
				WithRequest("GET", "/admin/teams/99643109-adb0-4e68-b25f-7b14d6bcae16", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))

				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(team))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.ReadTeam(team)

				assert.NoError(t, e)
				assert.NotNil(t, res)
				assert.Equal(t, "terraform-team", res.Name)
				assert.Equal(t, "99643109-adb0-4e68-b25f-7b14d6bcae16", res.UUID)
				assert.Len(t, res.Embedded.Members, 1)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("CreateTeam", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a pacticipant with name terraform-client exists").
				UponReceiving("a request to create a team").
				WithRequest("POST", "/admin/teams", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(create))
				}).
				WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(created))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.CreateTeam(create)

				assert.NoError(t, e)
				assert.NotNil(t, res)
				assert.Equal(t, "terraform-team", res.Name)
				assert.Equal(t, updated.UUID, res.UUID)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("UpdateTeam", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a team with uuid 99643109-adb0-4e68-b25f-7b14d6bcae16 exists").
				UponReceiving("a request to update a team").
				WithRequest("PUT", "/admin/teams/99643109-adb0-4e68-b25f-7b14d6bcae16", func(and *consumer.V2RequestBuilder) {
					and.Header("Content-Type", S("application/json"))
					and.Header("Authorization", Like("Bearer 1234"))
					and.JSONBody(Like(update))
				}).
				WillRespondWith(200, func(and *consumer.V2ResponseBuilder) {
					and.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					and.JSONBody(Like(updated))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.UpdateTeam(update)

				assert.NoError(t, e)
				assert.Equal(t, "terraform-team", res.Name)
				assert.Len(t, res.Embedded.Administrators, 1)
				assert.Len(t, res.Embedded.Environments, 1)
				assert.Len(t, res.Embedded.Pacticipants, 1)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("DeleteTeam", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a team with name terraform-team exists").
				UponReceiving("a request to delete a team").
				WithRequest("DELETE", "/admin/teams/99643109-adb0-4e68-b25f-7b14d6bcae16", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
				}).
				WillRespondWith(204)

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				return client.DeleteTeam(updated)
			})
			assert.NoError(t, err)
		})

		t.Run("UpdateTeamAssignments", func(t *testing.T) {
			req := broker.TeamsAssignmentRequest{
				UUID: updated.UUID,
				Users: []string{
					"05064a18-229d-4dfd-b37c-f00ec9673a49",
				},
			}

			mockProvider.
				AddInteraction().
				Given("a team with name terraform-team and user with uuid 05064a18-229d-4dfd-b37c-f00ec9673a49 exists").
				UponReceiving("a request to update team assignments").
				WithRequest("PUT", "/admin/teams/99643109-adb0-4e68-b25f-7b14d6bcae16/users", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(req))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(broker.TeamsAssignmentResponse{
						Embedded: broker.EmbeddedUsers{
							Users: []broker.TeamUser{
								{
									UUID:   "05064a18-229d-4dfd-b37c-f00ec9673a49",
									Active: true,
								},
							},
						},
					})
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, err := client.UpdateTeamAssignments(req)

				assert.Len(t, res.Embedded.Users, 1)

				return err
			})
			assert.NoError(t, err)
		})
	})

	t.Run("Secret", func(t *testing.T) {
		secret := broker.Secret{
			Name:        "terraformSecret",
			Description: "terraform secret",
			Value:       "supersecret",
			TeamUUID:    "1da4bc0e-8031-473f-880b-3b3951683284",
		}

		created := broker.Secret{
			UUID:        "b6af03cd-018c-4f1b-9546-c778d214f305",
			Name:        secret.Name,
			Description: secret.Description,
		}

		update := broker.Secret{
			UUID:        created.UUID,
			Name:        secret.Name,
			Description: "updated description",
			Value:       "topsecret",
		}

		updated := broker.Secret{
			UUID:        created.UUID,
			Name:        secret.Name,
			Description: "updated description",
		}

		t.Run("CreateSecret", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a team with uuid 1da4bc0e-8031-473f-880b-3b3951683284 exists").
				UponReceiving("a request to create a secret").
				WithRequest("POST", "/secrets", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(secret))
				}).
				WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(created))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.CreateSecret(secret)
				assert.NoError(t, e)
				assert.Equal(t, "terraformSecret", res.Name)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("UpdateSecret", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a secret with uuid b6af03cd-018c-4f1b-9546-c778d214f305 exists").
				UponReceiving("a request to update a secret").
				WithRequest("PUT", "/secrets/b6af03cd-018c-4f1b-9546-c778d214f305", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(update))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(updated))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.UpdateSecret(update)
				assert.NoError(t, e)
				assert.Equal(t, "terraformSecret", res.Name)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("DeleteSecret", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a secret with uuid b6af03cd-018c-4f1b-9546-c778d214f305 exists").
				UponReceiving("a request to delete a secret").
				WithRequest("DELETE", "/secrets/b6af03cd-018c-4f1b-9546-c778d214f305", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
				}).
				WillRespondWith(204, func(b *consumer.V2ResponseBuilder) {})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				return client.DeleteSecret(created)
			})
			assert.NoError(t, err)
		})
	})

	t.Run("Role", func(t *testing.T) {
		role := broker.Role{
			Name: "terraform-role",
			Permissions: []broker.Permission{
				{
					Scope:       "user:manage:*",
					Label:       "permission label",
					Description: "premission description",
				},
			},
		}

		created := broker.Role{
			UUID:        "e1407277-2a25-4559-8fed-4214dd12a1e8",
			Name:        role.Name,
			Permissions: role.Permissions,
		}

		update := broker.Role{
			UUID:        created.UUID,
			Name:        role.Name,
			Permissions: created.Permissions,
		}

		t.Run("CreateRole", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				UponReceiving("a request to create a role").
				WithRequest("POST", "/admin/roles", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(role))
				}).
				WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(created))

				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.CreateRole(role)
				assert.NoError(t, e)
				assert.Equal(t, "terraform-role", res.Name)
				assert.Len(t, res.Permissions, 1)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("ReadRole", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a role with uuid e1407277-2a25-4559-8fed-4214dd12a1e8 exists").
				UponReceiving("a request to get a role").
				WithRequest("GET", "/admin/roles/e1407277-2a25-4559-8fed-4214dd12a1e8", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(created))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.ReadRole(created.UUID)
				assert.NoError(t, e)
				assert.Equal(t, "terraform-role", res.Name)
				assert.Len(t, res.Permissions, 1)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("UpdateRole", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a role with uuid e1407277-2a25-4559-8fed-4214dd12a1e8 exists").
				UponReceiving("a request to update a role").
				WithRequest("PUT", "/admin/roles/e1407277-2a25-4559-8fed-4214dd12a1e8", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(update))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(update))

				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.UpdateRole(update)
				assert.NoError(t, e)
				assert.Equal(t, "terraform-role", res.Name)
				assert.Len(t, res.Permissions, 1)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("DeleteRole", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a role with uuid e1407277-2a25-4559-8fed-4214dd12a1e8 exists").
				UponReceiving("a request to delete a role").
				WithRequest("DELETE", "/admin/roles/e1407277-2a25-4559-8fed-4214dd12a1e8", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
				}).
				WillRespondWith(204)

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				return client.DeleteRole(created)
			})
			assert.NoError(t, err)
		})
	})

	t.Run("User", func(t *testing.T) {
		user := broker.User{
			Name:   "terraform user",
			Email:  "terraform.user@some.domain",
			Active: true,
			Type:   broker.RegularUser,
		}

		created := broker.User{
			UUID:   "819f6dbf-dd7a-47ff-b369-e3ed1d2578a0",
			Name:   user.Name,
			Email:  user.Email,
			Active: user.Active,
			Type:   user.Type,
			Embedded: struct {
				Roles []broker.Role `json:"roles,omitempty"`
				Teams []broker.Team `json:"teams,omitempty"`
			}{
				Roles: []broker.Role{
					{
						Name: "terraform-role",
						UUID: "84f66fab-1c42-4351-96bf-88d3a09d7cd2",
						Permissions: []broker.Permission{
							{
								Scope:       "user:manage:*",
								Label:       "permission label",
								Description: "permission description",
							},
						},
					},
				},
			},
		}

		update := created

		setUserRoles := broker.SetUserRolesRequest{
			Roles: []string{
				"84f66fab-1c42-4351-96bf-88d3a09d7cd2",
			},
		}

		t.Run("CreateUser", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				UponReceiving("a request to create a user").
				WithRequest("POST", "/admin/users/invite-user", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(user))

				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(created))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.CreateUser(user)
				assert.NoError(t, e)
				assert.Equal(t, "terraform user", res.Name)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("ReadUser", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a user with uuid 819f6dbf-dd7a-47ff-b369-e3ed1d2578a0 exists").
				UponReceiving("a request to get a user").
				WithRequest("GET", "/admin/users/819f6dbf-dd7a-47ff-b369-e3ed1d2578a0", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(created))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.ReadUser(created.UUID)
				assert.NoError(t, e)
				assert.Equal(t, "terraform user", res.Name)
				assert.Len(t, res.Embedded.Roles, 1)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("UpdateUser", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a user with uuid 819f6dbf-dd7a-47ff-b369-e3ed1d2578a0 exists").
				UponReceiving("a request to update a user").
				WithRequest("PUT", "/admin/users/819f6dbf-dd7a-47ff-b369-e3ed1d2578a0", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(update))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(update))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.UpdateUser(update)
				assert.NoError(t, e)
				assert.Equal(t, "terraform user", res.Name)
				assert.Len(t, res.Embedded.Roles, 1)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("DeleteUser", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a user with uuid 819f6dbf-dd7a-47ff-b369-e3ed1d2578a0 exists").
				UponReceiving("a request to delete a user").
				WithRequest("PUT", "/admin/users/819f6dbf-dd7a-47ff-b369-e3ed1d2578a0", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(update))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(update))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				return client.DeleteUser(created)
			})
			assert.NoError(t, err)
		})

		t.Run("SetUserRoles", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a user with uuid 819f6dbf-dd7a-47ff-b369-e3ed1d2578a0 exists").
				UponReceiving("a request to set user's roles").
				WithRequest("PUT", "/admin/users/819f6dbf-dd7a-47ff-b369-e3ed1d2578a0/roles", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(setUserRoles))

				}).
				WillRespondWith(200)

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				return client.SetUserRoles(created.UUID, setUserRoles)
			})
			assert.NoError(t, err)
		})
	})

	t.Run("SystemAccount", func(t *testing.T) {
		user := broker.User{
			Name:   "system account",
			Active: true,
			Type:   broker.SystemAccount,
		}

		created := broker.User{
			UUID:   "71a5be7d-bb9c-427b-ba49-ee8f1df0ae58",
			Name:   user.Name,
			Active: user.Active,
			Type:   user.Type,
			Embedded: struct {
				Roles []broker.Role `json:"roles,omitempty"`
				Teams []broker.Team `json:"teams,omitempty"`
			}{
				Roles: []broker.Role{
					{
						Name: "terraform-role",
						UUID: "84f66fab-1c42-4351-96bf-88d3a09d7cd2",
						Permissions: []broker.Permission{
							{
								Scope:       "user:manage:*",
								Label:       "permission label",
								Description: "permission description",
							},
						},
					},
				},
			},
		}

		update := created

		setUserRoles := broker.SetUserRolesRequest{
			Roles: []string{
				"84f66fab-1c42-4351-96bf-88d3a09d7cd2",
			},
		}

		t.Run("CreateSystemAccount", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				UponReceiving("a request to create a system account").
				WithRequest("POST", "/admin/system-accounts", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(user))

				}).
				WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.Header("Location", S(fmt.Sprintf("https://foo.com/path/to/%s", created.UUID)))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.CreateSystemAccount(user)
				assert.NoError(t, e)
				assert.Equal(t, created.UUID, res.UUID)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("ReadSystemAccount", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a system account with uuid 71a5be7d-bb9c-427b-ba49-ee8f1df0ae58 exists").
				UponReceiving("a request to get a system account").
				WithRequest("GET", "/admin/users/71a5be7d-bb9c-427b-ba49-ee8f1df0ae58", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(created))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.ReadUser(created.UUID)
				assert.NoError(t, e)
				assert.Equal(t, "system account", res.Name)
				assert.Len(t, res.Embedded.Roles, 1)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("UpdateSystemAccount", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a system account with uuid 71a5be7d-bb9c-427b-ba49-ee8f1df0ae58 exists").
				UponReceiving("a request to update a system account").
				WithRequest("PUT", "/admin/users/71a5be7d-bb9c-427b-ba49-ee8f1df0ae58", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(update))

				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(update))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.UpdateUser(update)
				assert.NoError(t, e)
				assert.Equal(t, "system account", res.Name)
				assert.Len(t, res.Embedded.Roles, 1)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("DeleteSystemAccount", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a system account with uuid 71a5be7d-bb9c-427b-ba49-ee8f1df0ae58 exists").
				UponReceiving("a request to delete a system account").
				WithRequest("PUT", "/admin/users/71a5be7d-bb9c-427b-ba49-ee8f1df0ae58", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(update))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(update))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				return client.DeleteUser(created)
			})
			assert.NoError(t, err)
		})

		t.Run("SetSystemAccountRoles", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a system account with uuid 71a5be7d-bb9c-427b-ba49-ee8f1df0ae58 exists").
				UponReceiving("a request to set a system account's roles").
				WithRequest("PUT", "/admin/users/71a5be7d-bb9c-427b-ba49-ee8f1df0ae58/roles", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(setUserRoles))
				}).
				WillRespondWith(200)

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				return client.SetUserRoles(created.UUID, setUserRoles)
			})
			assert.NoError(t, err)
		})
	})

	t.Run("Token", func(t *testing.T) {
		readOnlytoken := broker.APIToken{
			UUID:        "e068b1ea-d064-4719-971f-98af49cdf3f7",
			Description: "Read only token (developer)",
			Value:       "1234",
		}
		readWriteToken := broker.APIToken{
			UUID:        "cb32752b-0f1b-4f6f-817f-4b12b7cc8592",
			Description: "Read/write token (CI)",
			Value:       "5678",
		}

		t.Run("ReadTokens", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a token with uuid e068b1ea-d064-4719-971f-98af49cdf3f7 exists").
				UponReceiving("a request to get the tokens for the current account").
				WithRequest("GET", "/settings/tokens", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(map[string]interface{}{
						"_embedded": map[string]interface{}{
							"items": []interface{}{
								map[string]interface{}{
									"uuid":        Like(readOnlytoken.UUID),
									"description": readOnlytoken.Description,
									"value":       Like(readOnlytoken.Value),
								},
								map[string]interface{}{
									"uuid":        Like(readWriteToken.UUID),
									"description": readWriteToken.Description,
									"value":       Like(readWriteToken.Value),
								},
							},
						},
					}))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.ReadToken("e068b1ea-d064-4719-971f-98af49cdf3f7")

				assert.Equal(t, "1234", res.Value)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("RegenerateToken", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a token with uuid e068b1ea-d064-4719-971f-98af49cdf3f7 exists").
				UponReceiving("a request to regenerate a token").
				WithRequest("POST", "/settings/tokens/e068b1ea-d064-4719-971f-98af49cdf3f7/regenerate", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
				}).
				WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(readOnlytoken))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.RegenerateToken(readOnlytoken)
				assert.Equal(t, "1234", res.Value)

				return e
			})
			assert.NoError(t, err)
		})
	})

	t.Run("Webhook", func(t *testing.T) {
		webhook := broker.Webhook{
			Description: "terraform webhook",
			TeamUUID:    "607fba87-8209-4aff-a7d2-d8e9f92b94a2",
			Enabled:     true,
			Events: []broker.WebhookEvent{
				{Name: "contract_content_changed"},
				{Name: "contract_published"},
				{Name: "provider_verification_failed"},
				{Name: "provider_verification_published"},
				{Name: "provider_verification_succeeded"},
			},
			Provider: &broker.Pacticipant{
				Name: "terraform-provider",
			},
			Consumer: &broker.Pacticipant{
				Name: "terraform-consumer",
			},
			Request: broker.Request{
				Method:   "POST",
				URL:      "https://postman-echo.com/post",
				Username: "user",
				Password: "password",
				Headers: broker.Headers{
					"content-type": "application/json",
				},
				Body: map[string]string{
					"pact": "$${pactbroker.pactUrl}",
				},
			},
		}

		c, _ := copystructure.Copy(&webhook)
		created := c.(*broker.Webhook)
		created.ID = "2e4bf0e6-b0cf-451f-b05b-69048955f019"
		created.Request.Password = ""

		t.Run("CreateWebhook", func(t *testing.T) {

			mockProvider.
				AddInteraction().
				Given("a team with uuid 607fba87-8209-4aff-a7d2-d8e9f92b94a2 exists").
				UponReceiving("a request to create a webhook").
				WithRequest("POST", "/webhooks", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(webhook))

				}).
				WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(created))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.CreateWebhook(webhook)
				assert.NoError(t, e)
				assert.Equal(t, "terraform webhook", res.Description)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("ReadWebhook", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a webhook with ID 2e4bf0e6-b0cf-451f-b05b-69048955f019 exists").
				UponReceiving("a request to get a webhook").
				WithRequest("GET", "/webhooks/2e4bf0e6-b0cf-451f-b05b-69048955f019", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(created))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.ReadWebhook("2e4bf0e6-b0cf-451f-b05b-69048955f019")
				assert.NoError(t, e)
				assert.Equal(t, "terraform webhook", res.Description)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("UpdateWebhook", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a webhook with ID 2e4bf0e6-b0cf-451f-b05b-69048955f019 exists").
				UponReceiving("a request to update a webhook").
				WithRequest("PUT", "/webhooks/2e4bf0e6-b0cf-451f-b05b-69048955f019", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(created))

				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(created))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.UpdateWebhook(*created)
				assert.NoError(t, e)
				assert.Equal(t, "terraform webhook", res.Description)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("DeleteWebhook", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a webhook with ID 2e4bf0e6-b0cf-451f-b05b-69048955f019 exists").
				UponReceiving("a request to delete a webhook").
				WithRequest("DELETE", "/webhooks/2e4bf0e6-b0cf-451f-b05b-69048955f019", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
				}).
				WillRespondWith(204)

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				return client.DeleteWebhook(*created)
			})
			assert.NoError(t, err)
		})
	})

	t.Run("AuthenticationSettings", func(t *testing.T) {
		authSettings := broker.AuthenticationSettings{
			Providers: broker.AuthenticationProviders{
				Google: broker.GoogleAuthenticationSettings{
					EmailDomains: []string{"pactflow.io"},
				},
				Github: broker.GithubAuthenticationSettings{
					Organizations: []string{"pactflow"},
				},
			},
		}

		t.Run("SetTenantAuthenticationSettings", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				UponReceiving("a request to update authentication settings").
				WithRequest("PUT", "/admin/tenant/authentication-settings", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(authSettings))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/json;charset=utf-8"))
					b.JSONBody(Like(authSettings))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.SetTenantAuthenticationSettings(authSettings)
				assert.NoError(t, e)
				assert.Contains(t, res.Providers.Google.EmailDomains, "pactflow.io")
				assert.Contains(t, res.Providers.Github.Organizations, "pactflow")

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("ReadTenantAuthenticationSettings", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				UponReceiving("a request to get authentication settings").
				WithRequest("GET", "/admin/tenant/authentication-settings").
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/json;charset=utf-8"))
					b.JSONBody(Like(authSettings))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.ReadTenantAuthenticationSettings()
				assert.NoError(t, e)
				assert.Contains(t, res.Providers.Google.EmailDomains, "pactflow.io")
				assert.Contains(t, res.Providers.Github.Organizations, "pactflow")

				return e
			})
			assert.NoError(t, err)
		})
	})

	t.Run("Environment", func(t *testing.T) {
		environment := broker.Environment{
			UUID:        "8000883c-abf0-4b4c-b993-426f607092a9",
			Name:        "TerraformEnvironment",
			Production:  true,
			DisplayName: "terraform environment",
			Embedded: broker.EnvironmentEmbeddedItems{
				Teams: []broker.EnvironmentEmbeddedTeams{
					{
						Name: "terraform-team",
						UUID: "99643109-adb0-4e68-b25f-7b14d6bcae16",
					},
				},
			},
		}

		create := broker.EnvironmentCreateOrUpdateRequest{
			DisplayName: environment.DisplayName,
			Name:        environment.Name,
			Production:  environment.Production,
			Teams: []string{
				"99643109-adb0-4e68-b25f-7b14d6bcae16",
			},
		}

		created := broker.EnvironmentCreateOrUpdateResponse{
			UUID:        environment.UUID,
			Name:        environment.Name,
			DisplayName: environment.DisplayName,
			Production:  environment.Production,
			Teams: []string{
				"99643109-adb0-4e68-b25f-7b14d6bcae16",
			},
			Embedded: broker.EnvironmentEmbeddedItems{
				Teams: []broker.EnvironmentEmbeddedTeams{
					{
						Name: "terraform-team",
						UUID: "99643109-adb0-4e68-b25f-7b14d6bcae16",
					},
				},
			},
		}

		update := broker.EnvironmentCreateOrUpdateRequest{
			UUID:        created.UUID,
			DisplayName: environment.DisplayName,
			Name:        "terraform-updated-environment",
			Production:  environment.Production,
			Teams: []string{
				"99643109-adb0-4e68-b25f-7b14d6bcae16",
			},
		}

		updated := broker.EnvironmentCreateOrUpdateResponse{
			UUID:        environment.UUID,
			Name:        update.Name,
			DisplayName: environment.DisplayName,
			Embedded: broker.EnvironmentEmbeddedItems{
				Teams: []broker.EnvironmentEmbeddedTeams{
					{
						Name: "terraform-team",
						UUID: "99643109-adb0-4e68-b25f-7b14d6bcae16",
					},
				},
			},
		}

		t.Run("CreateEnvironment", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("a team with uuid 99643109-adb0-4e68-b25f-7b14d6bcae16 exists").
				UponReceiving("a request to create an environment").
				WithRequest("POST", "/environments", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(create))
				}).
				WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(created))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.CreateEnvironment(create)
				assert.NoError(t, e)
				assert.Equal(t, "TerraformEnvironment", res.Name)
				assert.Len(t, res.Embedded.Teams, 1)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("ReadEnvironment", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("an environment with uuid 8000883c-abf0-4b4c-b993-426f607092a9 exists").
				UponReceiving("a request to get an environment").
				WithRequest("GET", "/environments/8000883c-abf0-4b4c-b993-426f607092a9", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/hal+json;charset=utf-8"))
					b.JSONBody(Like(environment))
				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.ReadEnvironment(created.UUID)
				assert.NoError(t, e)
				assert.Equal(t, "TerraformEnvironment", res.Name)
				assert.Len(t, res.Embedded.Teams, 1)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("UpdateEnvironment", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("an environment with uuid 8000883c-abf0-4b4c-b993-426f607092a9 exists").
				UponReceiving("a request to update an environment").
				WithRequest("PUT", "/environments/8000883c-abf0-4b4c-b993-426f607092a9", func(b *consumer.V2RequestBuilder) {
					b.Header("Content-Type", S("application/json"))
					b.Header("Authorization", Like("Bearer 1234"))
					b.JSONBody(Like(update))

				}).
				WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
					b.Header("Content-Type", S("application/json;charset=utf-8"))
					b.JSONBody(Like(updated))

				})

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				res, e := client.UpdateEnvironment(update)
				assert.NoError(t, e)
				assert.Equal(t, "terraform-updated-environment", res.Name)
				assert.Len(t, res.Embedded.Teams, 1)

				return e
			})
			assert.NoError(t, err)
		})

		t.Run("DeleteEnvironment", func(t *testing.T) {
			mockProvider.
				AddInteraction().
				Given("an environment with uuid 8000883c-abf0-4b4c-b993-426f607092a9 exists").
				UponReceiving("a request to delete an environment").
				WithRequest("DELETE", "/environments/8000883c-abf0-4b4c-b993-426f607092a9", func(b *consumer.V2RequestBuilder) {
					b.Header("Authorization", Like("Bearer 1234"))
				}).
				WillRespondWith(204)

			err = mockProvider.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := clientForPact(config)

				return client.DeleteEnvironment(broker.Environment{
					UUID: created.UUID,
				})
			})
			assert.NoError(t, err)
		})
	})
}

func clientForPact(config consumer.MockServerConfig) *Client {
	baseURL, err := url.Parse(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
	if err != nil {
		panic(fmt.Sprintf("unable to create client for pact test: %s", err))
	}

	return NewClient(nil, Config{
		AccessToken: "1234",
		BaseURL:     baseURL,
	})
}

var Like = matchers.Like
var EachLike = matchers.EachLike
var Term = matchers.Term
var Regex = matchers.Regex
var HexValue = matchers.HexValue
var Identifier = matchers.Identifier
var IPAddress = matchers.IPAddress
var IPv6Address = matchers.IPv6Address
var Timestamp = matchers.Timestamp
var Date = matchers.Date
var Time = matchers.Time
var UUID = matchers.UUID

type S = matchers.String

var ArrayMinLike = matchers.ArrayMinLike
