package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pact-foundation/terraform/broker"
	"github.com/pact-foundation/terraform/client"
)

func role() *schema.Resource {
	return &schema.Resource{
		Create: roleCreate,
		Read:   roleRead,
		Update: roleUpdate,
		Delete: roleDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Role",
			},
			"scopes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "The pre-defined scope to add to the role",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the Role",
			},
		},
	}
}

func validateScopes(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	for _, scope := range broker.AllowedScopes {
		if scope == v {
			return
		}
	}
	errs = append(errs, fmt.Errorf("%q must be one of the allowed scopes %v, got %v", key, broker.AllowedScopes, v))

	return
}

func roleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	name := d.Get("name").(string)
	raw, ok := d.Get("scopes").([]interface{})

	permissions := make([]broker.Permission, len(raw))
	if ok && len(raw) > 0 {
		for i, s := range raw {
			permissions[i] = broker.Permission{
				Scope: s.(string),
			}
		}
	}

	role := broker.Role{
		Name:        name,
		Permissions: permissions,
	}

	created, err := client.CreateRole(role)

	if err == nil {
		d.SetId(created.UUID)
		d.Set("name", created.Name)
		d.Set("uuid", created.UUID)
	}

	return err
}

func roleRead(d *schema.ResourceData, meta interface{}) error {
	// TODO:
	return nil
}

func roleUpdate(d *schema.ResourceData, meta interface{}) error {
	// TODO:
	return nil
}

func roleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	uuid := d.Get("uuid").(string)

	log.Println("[DEBUG] deleting role for user with UUID:", uuid)

	err := client.DeleteRole(broker.Role{
		UUID: uuid,
	})

	if err != nil {
		d.SetId("")
	}

	return err
}

// curl -X POST -H"content-type: application/json" $PACT_BROKER_BASE_URL/admin/roles -d '{ "name": "FooRole", "permissions": [ { "name": "Manage users", "scope": "user:manage:*" } ] }' -H"Authorization: bearer $PACT_BROKER_TOKEN" -v  | jq .
