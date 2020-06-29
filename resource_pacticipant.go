package main

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pact-foundation/terraform/broker"
	"github.com/pact-foundation/terraform/client"
)

func pacticipant() *schema.Resource {
	return &schema.Resource{
		Create:   pacticipantCreate,
		Update:   pacticipantUpdate,
		Read:     pacticipantRead,
		Delete:   pacticipantDelete,
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Pacticipant",
			},
			"repository_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL or location of the VCS repository",
			},
		},
	}
}

func pacticipantCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	name := d.Get("name").(string)
	url := d.Get("repository_url").(string)

	log.Println("[DEBUG] creating pacticipant", name)

	pacticipant := broker.Pacticipant{
		Name:          name,
		RepositoryURL: url,
	}
	_, err := client.CreatePacticipant(pacticipant)

	if err == nil {
		d.SetId(name)
		d.Set("repository_url", pacticipant.RepositoryURL)
	}

	return err
}
func pacticipantUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	name := d.Get("name").(string)
	url := d.Get("repository_url").(string)

	log.Println("[DEBUG] updating pacticipant", name)

	pacticipant := broker.Pacticipant{
		Name:          name,
		RepositoryURL: url,
	}
	_, err := client.UpdatePacticipant(pacticipant)

	if err == nil {
		d.SetId(name)
		d.Set("repository_url", pacticipant.RepositoryURL)
	}

	return err
}

func pacticipantRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	log.Println("[DEBUG] reading pacticipant", d.Id())

	pacticipant, err := client.ReadPacticipant(d.Id())

	log.Println("[DEBUG] have pacticipant for READ", pacticipant)

	if err == nil {
		d.SetId(pacticipant.Name)
		d.Set("repository_url", pacticipant.RepositoryURL)
	}

	return err
}

func pacticipantDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	name := d.Get("name").(string)
	url := d.Get("repository_url").(string)

	log.Println("[DEBUG] deleting pacticipant", name)

	err := client.DeletePacticipant(broker.Pacticipant{
		Name:          name,
		RepositoryURL: url,
	})

	if err != nil {
		d.SetId("")
	}

	return err
}
