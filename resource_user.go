package main

import (
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
		},
	}
}

func userCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	name := d.Get("name").(string)
	email := d.Get("email").(string)
	active := d.Get("active").(bool)

	log.Println("[DEBUG] creating user", name)

	user := broker.User{
		Name:      name,
		Email:     email,
		Active:    active,
		FirstName: "foo",
		LastName:  "foo",
	}

	created, err := client.CreateUser(user)

	if err == nil {
		user.UUID = created.UUID
		d.SetId(created.UUID)
		setUserState(d, *created)
	}

	// Update user (any other params aren't set on first user creation - e.g. name, active etc.)
	_, err = client.UpdateUser(user)

	return err
}

func userUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	uuid := d.Id()
	name := d.Get("name").(string)
	email := d.Get("email").(string)
	active := d.Get("active").(bool)
	user := broker.User{
		UUID:   uuid,
		Email:  email,
		Active: active,
		Name:   name,
	}

	log.Println("[DEBUG] updating user", user)

	updated, err := client.UpdateUser(user)

	if err == nil {
		setUserState(d, *updated)
	}

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
		UUID: uuid,
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

	return nil
}
