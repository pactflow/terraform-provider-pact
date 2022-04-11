package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/pactflow/terraform/broker"
	"github.com/pactflow/terraform/client"
)

func environment() *schema.Resource {
	return &schema.Resource{
		Create:   environmentCreate,
		Update:   environmentUpdate,
		Read:     environmentRead,
		Delete:   environmentDelete,
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ValidateFunc: validateAlphaNumeric,
				Required:     true,
				Description:  "Name of the Environment",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the Environment",
			},
			"production": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Is this environment a production environment?",
			},
			"teams": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of teams (as uuids) that may use the environment",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of environment",
			},
		},
	}
}

func validateAlphaNumeric(v interface{}, k string) (warns []string, errs []error) {
	return validation.All(
		validation.StringLenBetween(1, 255),
		validation.StringMatch(
			regexp.MustCompile("[a-zA-Z0-9]+"),
			"The environment name must consist of alphanumerics"),
	)(v, k)
}

func environmentToCRUD(environment broker.Environment, teams []string) broker.EnvironmentCreateOrUpdateRequest {
	return broker.EnvironmentCreateOrUpdateRequest{
		DisplayName: environment.DisplayName,
		Name:        environment.Name,
		Production:  environment.Production,
		Teams:       teams,
		UUID:        environment.UUID,
	}
}

func environmentFromCRUD(environment broker.EnvironmentCreateOrUpdateResponse) broker.Environment {
	return broker.Environment{
		DisplayName: environment.DisplayName,
		Name:        environment.Name,
		Production:  environment.Production,
		CreatedAt:   environment.CreatedAt,
		UpdatedAt:   environment.UpdatedAt,
		UUID:        environment.UUID,
		Embedded:    environment.Embedded,
	}
}

func environmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	environment := getEnvironmentFromState(d)

	teams := d.Get("teams").([]interface{})
	log.Println("[DEBUG] creating environment", environment, teams)

	created, err := client.CreateEnvironment(environmentToCRUD(environment, arrayInterfaceToArrayString(teams)))

	if err != nil {
		return err
	}

	d.SetId(created.UUID)
	setEnvironmentState(d, environmentFromCRUD(*created))

	return nil
}

func environmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	environment := getEnvironmentFromState(d)
	teams := d.Get("teams").([]interface{})

	log.Println("[DEBUG] updating environment", environment)

	updated, err := client.UpdateEnvironment(environmentToCRUD(environment, arrayInterfaceToArrayString(teams)))

	if err != nil {
		return err
	}

	setEnvironmentState(d, environmentFromCRUD(*updated))

	return err
}

func environmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	uuid := d.Id()

	log.Println("[DEBUG] reading environment", uuid)

	environment, err := client.ReadEnvironment(uuid)

	if err == nil {
		d.SetId(environment.UUID)
		setEnvironmentState(d, *environment)
	}

	return err
}

func environmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)

	log.Println("[DEBUG] deleting environment", d.Id())

	err := client.DeleteEnvironment(getEnvironmentFromState(d))

	if err != nil {
		d.SetId("")
		return fmt.Errorf("unable to delete environment %s: %w", d.Id(), err)
	}

	return nil
}

func setEnvironmentState(d *schema.ResourceData, environment broker.Environment) error {
	log.Printf("[DEBUG] setting environment state: %+v \n", environment)

	if err := d.Set("name", environment.Name); err != nil {
		log.Println("[ERROR] error setting key 'name'", err)
		return err
	}
	if err := d.Set("display_name", environment.DisplayName); err != nil {
		log.Println("[ERROR] error setting key 'displayName'", err)
		return err
	}
	if err := d.Set("production", environment.Production); err != nil {
		log.Println("[ERROR] error setting key 'production'", err)
		return err
	}
	if err := d.Set("uuid", environment.UUID); err != nil {
		log.Println("[ERROR] error setting key 'uuid'", err)
		return err
	}
	if err := d.Set("teams", teamsFromEnvironment(environment)); err != nil {
		log.Println("[ERROR] error setting key 'teams'", err)
		return err
	}

	return nil
}

func teamsFromEnvironment(u broker.Environment) []string {
	teams := make([]string, len(u.Embedded.Teams))

	for i, r := range u.Embedded.Teams {
		teams[i] = r.UUID
	}

	sort.Strings(teams)

	return teams
}

func teamsFromStateChange(d *schema.ResourceData) []string {
	_, after := d.GetChange("teams")
	teams, ok := after.([]interface{})
	if !ok {
		return []string{}
	}
	return arrayInterfaceToArrayString(teams)
}

func getEnvironmentFromState(d *schema.ResourceData) broker.Environment {
	uuid := d.Id()
	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	production := d.Get("production").(bool)

	return broker.Environment{
		UUID:        uuid,
		DisplayName: displayName,
		Production:  production,
		Name:        name,
	}
}
