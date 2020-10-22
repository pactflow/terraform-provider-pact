package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pact-foundation/terraform/broker"
	"github.com/pact-foundation/terraform/client"
)

func user() *schema.Resource {
	return &schema.Resource{
		Create:   userCreate,
		Update:   userUpdate,
		Read:     userRead,
		Delete:   userDelete,
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the User",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Email address of the user",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Active status of the user",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of user",
			},
			"type": {
				Type:         schema.TypeString,
				Default:      allowedUserTypes[userType],
				Description:  "The type of user (regular/system)",
				Optional:     true,
				ValidateFunc: validateUserType,
			},

			// List of UUIDs
			"teams": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					// ValidateFunc: validateEvents,
				},
			},
		},
	}
}

const (
	userType   = "user"
	systemType = "system"
)

var allowedUserTypes = map[string]broker.UserType{
	userType:   broker.RegularUser,
	systemType: broker.SystemAccount,
}

func userTypeAsString(t broker.UserType) string {
	for k, v := range allowedUserTypes {
		if v == t {
			return k
		}
	}

	return ""
}

func validateUserType(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)

	if _, ok := allowedUserTypes[v]; !ok {
		errs = append(errs, fmt.Errorf("%q must be one of the allowed types %v, got %v", key, allowedUserTypes, v))
	}

	return
}

func userCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	user := getUserFromState(d)

	log.Println("[DEBUG] creating user", user)

	created, err := client.CreateUser(user)

	if err == nil {
		user.UUID = created.UUID
		d.SetId(created.UUID)

		// Update user (any other params aren't set on first user creation - e.g. name, active etc.)
		// This is ineffectual, as these details come from Cognito and will be reset each login :P
		// _, err = client.UpdateUser(user)

		setUserState(d, *created)
	}

	return err
}

func userUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	user := getUserFromState(d)

	log.Println("[DEBUG] updating user", user)

	updated, err := client.UpdateUser(user)

	if err == nil {
		setUserState(d, *updated)
	}

	// TODO: Team assignments

	return err
}

func userRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	uuid := d.Id()

	log.Println("[DEBUG] reading user", uuid)

	user, err := client.ReadUser(uuid)

	log.Println("[DEBUG] have user for READ", user)

	if err == nil {
		d.SetId(user.UUID)
		setUserState(d, *user)
	}

	return err
}

func userDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	uuid := d.Id()
	user := broker.User{
		UUID:   uuid,
		Active: false,
	}

	log.Println("[DEBUG] deleting user", user)

	err := client.DeleteUser(user)

	if err != nil {
		d.SetId("")
	}

	return err
}

func setUserState(d *schema.ResourceData, user broker.User) error {
	log.Printf("[DEBUG] setting user state: %+v \n", user)

	if err := d.Set("name", user.Name); err != nil {
		log.Println("[ERROR] error setting key 'name'", err)
		return err
	}
	if err := d.Set("email", user.Email); err != nil {
		log.Println("[ERROR] error setting key 'email'", err)
		return err
	}
	if err := d.Set("active", user.Active); err != nil {
		log.Println("[ERROR] error setting key 'active'", err)
		return err
	}
	if err := d.Set("uuid", user.UUID); err != nil {
		log.Println("[ERROR] error setting key 'uuid'", err)
		return err
	}
	if err := d.Set("type", userTypeAsString(user.Type)); err != nil {
		log.Println("[ERROR] error setting key 'type'", err)
		return err
	}

	return nil
}

func getUserFromState(d *schema.ResourceData) broker.User {
	uuid := d.Id()
	name := d.Get("name").(string)
	email := d.Get("email").(string)
	active := d.Get("active").(bool)
	userType := d.Get("type").(string)

	return broker.User{
		UUID:   uuid,
		Email:  email,
		Active: active,
		Name:   name,
		Type:   allowedUserTypes[userType],
	}
}
