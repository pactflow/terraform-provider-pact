package main

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pact-foundation/terraform/broker"
	"github.com/pact-foundation/terraform/client"
)

func teamAssignment() *schema.Resource {
	return &schema.Resource{
		Create:   teamAssignmentCreate,
		Update:   teamAssignmentUpdate,
		Read:     teamAssignmentRead,
		Delete:   teamAssignmentDelete,
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
		Schema: map[string]*schema.Schema{
			"team": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of team",
			},
			"users": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func interfaceToStringArray(o interface{}) []string {
	items := o.([]interface{})
	res := make([]string, len(items))
	for i, item := range items {
		res[i] = item.(string)
	}

	return res
}

func teamAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	uuid := d.Get("team").(string)
	d.SetId(uuid)

	log.Println("[DEBUG] creating team assignment", uuid)

	if d.HasChange("users") {
		old, new := d.GetChange("users")
		log.Println("[DEBUG] teamAssignmentCreate - change. old:", old, "new:", new)

		usersToAdd := diff(interfaceToStringArray(old), interfaceToStringArray(new))
		usersToDelete := diff(interfaceToStringArray(new), interfaceToStringArray(old))
		log.Println("[DEBUG] teamAssignmentCreate - deleting:", usersToDelete, "adding:", usersToAdd)

		// // TODO: fix https://dius.slack.com/archives/GENCG4LAU/p1603334036004100
		err := client.DeleteTeamAssignments(broker.TeamsAssignmentRequest{
			UUID:  uuid,
			Users: usersToDelete,
		})

		if err != nil {
			return err
		}

		// TODO: partial
		req := broker.TeamsAssignmentRequest{
			UUID:  uuid,
			Users: usersToAdd,
		}

		_, err = client.UpdateTeamAssignments(req)

		if err != nil {
			return err
		}

		res, err := client.ReadTeamAssignments(broker.Team{
			UUID: uuid,
		})

		if err != nil {
			return err
		}

		// return setTeamAssignmentState(d, *res)
		return setTeamAssignmentState(d, res)
	}
	return nil

}

func teamAssignmentUpdate(d *schema.ResourceData, meta interface{}) error {
	return teamAssignmentCreate(d, meta)
}

func teamAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	uuid := d.Id()

	log.Println("[DEBUG] reading team assignment", uuid)

	res, err := client.ReadTeamAssignments(broker.Team{
		UUID: uuid,
	})

	log.Println("[DEBUG] have team assignment for READ", res)

	if err == nil {
		d.SetId(uuid)
		return setTeamAssignmentState(d, res)
	}

	return err
}

func teamAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	uuid := d.Id()
	req := broker.TeamsAssignmentRequest{
		UUID:  uuid,
		Users: extractUsersFromState(d),
	}

	log.Println("[DEBUG] deleting team assignments", team)

	err := client.DeleteTeamAssignments(req)

	if err != nil {
		d.SetId("")
	}

	return err
}

func setTeamAssignmentState(d *schema.ResourceData, team *broker.TeamsAssignmentResponse) error {
	log.Printf("[DEBUG] setting team assignment state: %+v \n", team)

	if team != nil {

		if err := d.Set("users", extractUsersFromApiResponse(team)); err != nil {
			log.Println("[ERROR] error setting key 'users'", err)
			return err
		}
	}

	return nil
}

// Use this to find the delta, and delete them from the team
func extractUsersFromState(d *schema.ResourceData) []string {
	usersRaw := d.Get("users").([]interface{})
	users := make([]string, len(usersRaw))
	for i, u := range usersRaw {
		users[i] = u.(string)
	}

	return users
}

// Use this to find the delta, and delete them from the team
func extractUsersFromApiResponse(response *broker.TeamsAssignmentResponse) []string {
	users := make([]string, len(response.Embedded.Users))
	for i, u := range response.Embedded.Users {
		users[i] = u.UUID
	}

	return users
}

// // Use this to find the delta, and delete them from the team
// func usersToDelete(d *schema.ResourceData, response broker.TeamsAssignmentResponse) []string {
// 	wantUsers := extractUsersFromState(d)
// 	actualUsers := extractUsersFromApiResponse(&response)

// 	return diff(wantUsers, actualUsers)
// }

// // Use this to find the delta, and delete them from the team
// func usersToAdd(d *schema.ResourceData, response broker.TeamsAssignmentResponse) []string {
// 	wantUsers := extractUsersFromState(d)
// 	actualUsers := extractUsersFromApiResponse(&response)

// 	return diff(actualUsers, wantUsers)
// }

// Finds the items in b that don't exist in a
func diff(a, b []string) []string {
	diff := make([]string, 0)
	m := make(map[string]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}

	return diff
}

// Pros of a "team_assignment"
// -> Single API call to attach all users to a team

// Cons of a team assignment
// -> it's not actually a single resource, so no UUID etc.
