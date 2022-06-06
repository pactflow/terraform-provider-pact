package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pactflow/terraform/broker"
	"github.com/pactflow/terraform/client"
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
				Optional:    true,
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
				ForceNew:     true,
				ValidateFunc: validateUserType,
			},
			"roles": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A list of roles (as uuids) to apply to the user",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

const (
	userType   = "user"
	systemType = "system"
)

// Special role that can't be assigned via API: https://github.com/pactflow/terraform-provider-pact/issues/22
const TEAM_ADMINISTRATOR_ROLE = "d635f960-88f2-4f13-8043-4641a02dffa0"

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

	roles := ExpandStringSet(d.Get("roles").(*schema.Set))
	log.Println("[DEBUG] creating user", user, roles)

	var created *broker.User
	var err error
	if user.Type == broker.SystemAccount {
		created, err = client.CreateSystemAccount(user)
	} else {
		created, err = client.CreateUser(user)
	}

	if err != nil {
		return err
	}

	d.SetId(created.UUID)

	setUserState(d, *created)

	log.Println("[DEBUG] updating user roles", d.Id(), roles)

	err = client.SetUserRoles(d.Id(), broker.SetUserRolesRequest{
		Roles: roles,
	})

	if err != nil {
		// Creating a user is a non-atomic transaction, because roles is a separate API call
		d.Partial(true)
		log.Println("[ERROR] error updating user roles", err)
		return fmt.Errorf("error updating roles for user (%s / %s): %w", d.Id(), created.Email, err)
	}

	d.Set("roles", roles)

	return nil
}

func userUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	user := getUserFromState(d)

	log.Println("[DEBUG] updating user", user)

	updated, err := client.UpdateUser(user)

	if err != nil {
		return err
	}

	setUserState(d, *updated)

	if d.HasChange("roles") {
		roles := rolesFromStateChange(d)
		log.Println("[DEBUG] updating user roles", roles)

		err = client.SetUserRoles(d.Id(), broker.SetUserRolesRequest{
			Roles: roles,
		})

		if err != nil {
			d.Partial(true) // updating users is non-atomic, let the diff applier know this
			log.Println("[ERROR] error updating user roles", err)
			return fmt.Errorf("error updating roles for user (%s): %s", d.Id(), err)
		}

		d.Set("roles", roles)
	}

	return err
}

func userRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	uuid := d.Id()

	log.Println("[DEBUG] reading user", uuid)

	user, err := client.ReadUser(uuid)

	if err == nil {
		d.SetId(user.UUID)
		setUserState(d, *user)
	}

	return err
}

func userDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)
	uuid := d.Id()

	log.Println("[DEBUG] deleting user", d.Id())

	// TODO: Delete attached resources Roles and Teams, because a Users aren't deleted, but simply disabled
	user, err := client.ReadUser(uuid)
	if err != nil {
		log.Println("[ERROR] unable to fetch user for delete", user)
		return fmt.Errorf("unable to fetch user for delete: %w", err)
	}
	log.Println("[DEBUG] have user for delete", user)

	rolesToRemove := make([]string, len(user.Embedded.Roles))
	for i, r := range user.Embedded.Roles {
		rolesToRemove[i] = r.UUID
	}

	err = client.SetUserRoles(uuid, broker.SetUserRolesRequest{
		Roles: []string{},
	})

	if err != nil {
		return fmt.Errorf("unable to remove roles from user when deleting (disabling) user %s: %w", d.Id(), err)
	}

	for _, t := range user.Embedded.Teams {
		err = client.DeleteTeamAssignment(t, *user)
		if err != nil {
			return fmt.Errorf("unable to remove user %s (%s) from team %s (%s): %w", d.Id(), user.Email, t.UUID, t.Name, err)
		}
	}
	user.Embedded.Roles = nil
	user.Embedded.Teams = nil

	err = client.DeleteUser(*user)

	if err != nil {
		d.SetId("")
		return fmt.Errorf("unable to delete (disable) user %s: %w", d.Id(), err)
	}

	return nil
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
	if err := d.Set("type", user.Type); err != nil {
		log.Println("[ERROR] error setting key 'type'", err)
		return err
	}
	if err := d.Set("roles", rolesFromUser(user)); err != nil {
		log.Println("[ERROR] error setting key 'roles'", err)
		return err
	}

	return nil
}

func rolesFromUser(u broker.User) []string {
	roles := make([]string, 0)

	for _, r := range u.Embedded.Roles {
		// Exclude the team administrator role
		// See https://github.com/pactflow/terraform-provider-pact/issues/22
		if r.UUID != TEAM_ADMINISTRATOR_ROLE {
			roles = append(roles, r.UUID)
		}
	}

	return roles
}

func rolesFromStateChange(d *schema.ResourceData) []string {
	return ExpandStringSet(d.Get("roles").(*schema.Set))
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
