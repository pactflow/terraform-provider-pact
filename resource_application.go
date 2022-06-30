package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pactflow/terraform/broker"
	"github.com/pactflow/terraform/client"
)

func application() *schema.Resource {
	return &schema.Resource{
		Create:   applicationCreate,
		Update:   applicationUpdate,
		Read:     applicationRead,
		Delete:   applicationDelete,
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
			"main_branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Main (default) branch",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The display name of the pacticipant",
			},
		},
	}
}

func applicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	name := d.Get("name").(string)
	url := d.Get("repository_url").(string)
	branch := d.Get("main_branch").(string)
	displayName := d.Get("display_name").(string)

	log.Println("[DEBUG] creating pacticipant", name)

	pacticipant := broker.Pacticipant{
		Name:          name,
		RepositoryURL: url,
		MainBranch:    branch,
		DisplayName:   displayName,
	}
	_, err := client.CreatePacticipant(pacticipant)

	if err != nil {
		return fmt.Errorf("error creating application: %w", err)
	}

	d.SetId(name)
	d.Set("name", pacticipant.Name)
	d.Set("repository_url", pacticipant.RepositoryURL)
	d.Set("main_branch", pacticipant.MainBranch)
	d.Set("display_name", pacticipant.DisplayName)

	return nil
}

func applicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	name := d.Get("name").(string)
	url := d.Get("repository_url").(string)
	branch := d.Get("main_branch").(string)
	displayName := d.Get("display_name").(string)

	log.Println("[DEBUG] updating pacticipant", name)

	pacticipant := broker.Pacticipant{
		Name:          name,
		RepositoryURL: url,
		MainBranch:    branch,
		DisplayName:   displayName,
	}
	_, err := client.UpdatePacticipant(pacticipant)

	if err != nil {
		return fmt.Errorf("error updating application: %w", err)
	}

	d.SetId(name)
	d.Set("name", pacticipant.Name)
	d.Set("repository_url", pacticipant.RepositoryURL)
	d.Set("main_branch", pacticipant.MainBranch)
	d.Set("display_name", pacticipant.DisplayName)

	return nil
}

func applicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	log.Println("[DEBUG] reading pacticipant", d.Id())

	pacticipant, err := client.ReadPacticipant(d.Id())

	log.Println("[DEBUG] have pacticipant for READ", pacticipant)

	if err != nil {
		return fmt.Errorf("error reading application: %w", err)
	}

	d.SetId(pacticipant.Name)
	d.Set("name", pacticipant.Name)
	d.Set("repository_url", pacticipant.RepositoryURL)
	d.Set("main_branch", pacticipant.MainBranch)
	d.Set("display_name", pacticipant.DisplayName)

	return nil
}

func applicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	name := d.Get("name").(string)

	log.Println("[DEBUG] deleting pacticipant", name)

	err := client.DeletePacticipant(broker.Pacticipant{
		Name: name,
	})

	if err != nil {
		d.SetId("")
		return fmt.Errorf("error deleting application: %w", err)
	}

	return nil

}
