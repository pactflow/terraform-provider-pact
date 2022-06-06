package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pactflow/terraform/broker"
	"github.com/pactflow/terraform/client"
)

func team() *schema.Resource {
	return &schema.Resource{
		Create:   teamCreate,
		Update:   teamUpdate,
		Read:     teamRead,
		Delete:   teamDelete,
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Team",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of team",
			},
			"pacticipants": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"users": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"administrators": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func getTeamFromResourceData(d *schema.ResourceData) broker.Team {
	name := d.Get("name").(string)
	uuid := d.Id()
	team := broker.Team{
		UUID: uuid,
		Name: name,
	}

	pacticipants := ExpandStringSet(d.Get("pacticipants").(*schema.Set))
	log.Println("[DEBUG] resource_team.go pacticipants?", pacticipants)
	if len(pacticipants) > 0 {
		log.Println("[DEBUG] resource_team.go have pacticipants", len(pacticipants), pacticipants)
		items := make([]broker.Pacticipant, len(pacticipants))
		for i, p := range pacticipants {
			items[i] = broker.Pacticipant{
				Name: p,
			}
		}
		team.Embedded.Pacticipants = items
	}

	administrators := ExpandStringSet(d.Get("administrators").(*schema.Set))
	log.Println("[DEBUG] resource_team.go administrators?", administrators)
	if len(administrators) > 0 {
		log.Println("[DEBUG] resource_team.go have administrators", len(administrators), administrators)
		items := make([]broker.User, len(administrators))
		for i, p := range administrators {
			items[i] = broker.User{
				UUID: p,
			}
		}
		team.Embedded.Administrators = items
	}

	return team
}

// Removes any users from the team that shouldn't be there, and adds those that should
func assignTeamUsers(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	uuid := d.Id()

	log.Println("[DEBUG] assigning users to team", uuid)

	if d.HasChange("users") {
		old, new := d.GetChange("users")
		log.Println("[DEBUG] teamAssignmentCreate - change. old:", old, "new:", new)

		usersToAdd := ExpandStringSet(new.(*schema.Set))
		log.Println("[DEBUG] teamAssignmentCreate - setting users:", usersToAdd)

		req := broker.TeamsAssignmentRequest{
			UUID:  uuid,
			Users: usersToAdd,
		}

		res, err := client.UpdateTeamAssignments(req)

		if err != nil {
			return err
		}

		return setTeamAssignmentState(d, res)
	}

	return nil
}

func teamToCRUDRequest(t broker.Team) broker.TeamCreateOrUpdateRequest {
	pacticipants := make([]string, len(t.Embedded.Pacticipants))
	for i, a := range t.Embedded.Pacticipants {
		pacticipants[i] = a.Name
	}

	administrators := make([]string, len(t.Embedded.Administrators))
	for i, a := range t.Embedded.Administrators {
		administrators[i] = a.UUID
	}

	return broker.TeamCreateOrUpdateRequest{
		Name:               t.Name,
		UUID:               t.UUID,
		PacticipantNames:   pacticipants,
		AdministratorUUIDs: administrators,
	}
}

func teamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	team := getTeamFromResourceData(d)
	create := teamToCRUDRequest(team)

	log.Println("[DEBUG] creating team", team)

	created, err := client.CreateTeam(create)

	if err != nil {
		return fmt.Errorf("error creating team: %w", err)
	}

	team.UUID = created.UUID
	d.SetId(created.UUID)
	setTeamState(d, *created)

	err = assignTeamUsers(d, meta)
	if err != nil {
		d.Partial(true)
		log.Printf("\n\n[DEBUG] error assigning team users: %v \n\n", err)
		return fmt.Errorf("error assigning team users: %w", err)
	}

	return err
}

func teamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	team := getTeamFromResourceData(d)
	update := teamToCRUDRequest(team)

	log.Println("[DEBUG] updating team", team)

	updated, err := client.UpdateTeam(update)

	if err == nil {
		setTeamState(d, *updated)
	}

	err = assignTeamUsers(d, meta)
	if err != nil {
		d.Partial(true)
		return fmt.Errorf("error assigning team users: %w", err)
	}

	return err
}

func teamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	read := getTeamFromResourceData(d)

	log.Println("[DEBUG] reading team", read)

	team, err := client.ReadTeam(read)

	log.Println("[DEBUG] have team for READ", team)

	if err == nil {
		d.SetId(team.UUID)
		setTeamState(d, *team)
	}

	return err
}

func teamDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	uuid := d.Id()
	team := broker.Team{
		UUID: uuid,
	}

	log.Println("[DEBUG] deleting team", team)

	err := client.DeleteTeam(team)

	if err != nil {
		d.SetId("")
	}

	return err
}

func setTeamState(d *schema.ResourceData, team broker.Team) error {
	log.Printf("[DEBUG] setting team state: %+v \n", team)

	if err := d.Set("name", team.Name); err != nil {
		log.Println("[ERROR] error setting key 'name'", err)
		return err
	}
	if err := d.Set("uuid", team.UUID); err != nil {
		log.Println("[ERROR] error setting key 'uuid'", err)
		return err
	}

	if len(team.Embedded.Pacticipants) > 0 {
		pacticipants := make([]string, len(team.Embedded.Pacticipants))
		for i, p := range team.Embedded.Pacticipants {
			pacticipants[i] = p.Name
		}

		sort.Strings(pacticipants)

		if err := d.Set("pacticipants", pacticipants); err != nil {
			log.Println("[ERROR] error setting key 'pacticipants'", err)
			return err
		}
	}

	if len(team.Embedded.Members) > 0 {
		members := make([]string, len(team.Embedded.Members))
		for i, m := range team.Embedded.Members {
			log.Println("[DEBUG] adding team member with UUID", m.UUID)
			members[i] = m.UUID
		}

		sort.Strings(members)

		if err := d.Set("users", members); err != nil {
			log.Println("[ERROR] error setting key 'users'", err)
			return err
		}
	}

	if len(team.Embedded.Administrators) > 0 {
		administrators := make([]string, len(team.Embedded.Administrators))
		for i, a := range team.Embedded.Administrators {
			log.Println("[DEBUG] adding administrator with UUID", a.UUID)
			administrators[i] = a.UUID
		}

		sort.Strings(administrators)

		if err := d.Set("Administrators", administrators); err != nil {
			log.Println("[ERROR] error setting key 'administrators'", err)
			return err
		}
	}

	return nil
}

func setTeamAssignmentState(d *schema.ResourceData, team *broker.TeamsAssignmentResponse) error {
	log.Printf("[DEBUG] setting team assignment state: %+v \n", team)

	if team != nil {
		if err := d.Set("users", extractUsersFromAPIResponse(team)); err != nil {
			log.Println("[ERROR] error setting key 'users'", err)
			return err
		}
	}

	return nil
}

// Use this to find the delta, and delete them from the team
func extractUsersFromAPIResponse(response *broker.TeamsAssignmentResponse) []string {
	users := make([]string, len(response.Embedded.Users))
	for i, u := range response.Embedded.Users {
		users[i] = u.UUID
	}

	sort.Strings(users)

	return users
}
