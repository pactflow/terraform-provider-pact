package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pact-foundation/terraform/broker"
	"github.com/pact-foundation/terraform/client"
)

const (
	administratorRole = "administrator"
	userRole          = "user"
)

var allowedRoles = map[string]string{
	administratorRole: "cf75d7c2-416b-11ea-af5e-53c3b1a4efd8",
	// userRole:          "e9282e22-416b-11ea-a16e-57ee1bb61d18",
}

func validateRoles(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if _, ok := allowedRoles[v]; !ok {
		errs = append(errs, fmt.Errorf("%q must be one of the allowed pre-existing roles %v, got %v", key, allowedRoles, v))
	}

	return
}

func role() *schema.Resource {
	return &schema.Resource{
		Create: roleCreate,
		Read:   roleRead,
		Delete: roleDelete,
		Schema: map[string]*schema.Schema{
			"role": {
				Type:         schema.TypeString,
				Description:  "Role to apply to the user",
				ValidateFunc: validateRoles,
				Required:     true,
				ForceNew:     true,
			},
			"user": {
				Type:        schema.TypeString,
				Description: "UUID of the user of which to apply the role",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the Role",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of API token",
			},
		},
	}
}

func roleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	userUUID := d.Get("user").(string)

	// NOTE: we only support the admin role at this time
	log.Println("[DEBUG] creating role for user with UUID:", userUUID)
	_, err := client.AddAdminRoleToUser(broker.User{
		UUID: userUUID,
	})

	if err == nil {
		d.SetId(allowedRoles["administrator"])
		d.Set("name", "Administrator")
	}

	return err
}

func roleRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func roleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	userUUID := d.Get("user").(string)

	log.Println("[DEBUG] deleting role for user with UUID:", userUUID)

	_, err := client.RemoveAdminRoleToUser(broker.User{
		UUID: userUUID,
	})

	if err != nil {
		d.SetId("")
	}

	return err
}
