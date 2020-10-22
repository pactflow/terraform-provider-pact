package main

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pact-foundation/terraform/broker"
	"github.com/pact-foundation/terraform/client"
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func getTeamFromResourceData(d *schema.ResourceData) broker.Team {
	name := d.Get("name").(string)
	uuid := d.Id()
	team := broker.Team{
		UUID: uuid,
		Name: name,
	}

	pacticipants, ok := d.Get("pacticipants").([]interface{})
	log.Println("[DEBUG] resource_team.go pacticipants?", pacticipants, ok)
	if ok && len(pacticipants) > 0 {
		log.Println("[DEBUG] resource_team.go have pacticipants", len(pacticipants), pacticipants)
		items := make([]broker.Pacticipant, len(pacticipants))
		for i, p := range pacticipants {
			items[i] = broker.Pacticipant{
				Name: p.(string),
			}
		}
		team.Embedded.Pacticipants = items
	}

	// users, ok := d.Get("users").([]interface{})
	// log.Println("[DEBUG] resource_team.go users?", users, ok)
	// if ok && len(users) > 0 {
	// 	log.Println("[DEBUG] resource_team.go have users", len(users), users)
	// 	items := make([]broker.User, len(users))
	// 	for i, p := range pacticipants {
	// 		items[i] = broker.Pacticipant{
	// 			Name: p.(string),
	// 		}
	// 	}
	// 	team.Embedded.Pacticipants = items
	// }
	return team
}

func teamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	team := getTeamFromResourceData(d)

	log.Println("[DEBUG] creating team", team)

	created, err := client.CreateTeam(team)

	if err == nil {
		team.UUID = created.UUID
		d.SetId(created.UUID)
		setTeamState(d, *created)
	}

	return err
}

func teamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	team := getTeamFromResourceData(d)

	log.Println("[DEBUG] updating team", team)

	updated, err := client.UpdateTeam(team)

	if err == nil {
		setTeamState(d, *updated)
	}

	return err
}

func teamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	uuid := d.Id()

	log.Println("[DEBUG] reading team", uuid)

	team, err := client.ReadTeam(uuid)

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
	pacticipants := make([]string, len(team.Embedded.Pacticipants))
	for _, p := range team.Embedded.Pacticipants {
		pacticipants = append(pacticipants, p.Name)
	}

	if len(pacticipants) > 0 {
		if err := d.Set("pacticipants", pacticipants); err != nil {
			log.Println("[ERROR] error setting key 'pacticipants'", err)
			return err
		}
	}

	return nil
}
