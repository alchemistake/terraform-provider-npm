package npm

import (
	"fmt"
	"strings"

	"github.com/alchemistake/go-npm"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNPMUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceNPMUserCreateOrUpdate,
		Read:   resourceNPMUserRead,
		Update: resourceNPMUserCreateOrUpdate,
		Delete: resourceNPMUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNPMUserImporter,
		},

		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"org": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"role": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateValueFunc([]string{"developer", "admin", "owner"}),
			},
		},
	}
}

func resourceNPMUserCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	user := d.Get("user").(string)
	org := d.Get("org").(string)
	role := d.Get("role").(string)
	client := m.(*npm.Client)
	err := client.AddUser(org, user, role)
	if err != nil {
		return err
	}

	d.SetId(user)
	return resourceNPMUserRead(d, m)

}

func resourceNPMUserRead(d *schema.ResourceData, m interface{}) error {
	user := d.Get("user").(string)
	org := d.Get("org").(string)

	client := m.(*npm.Client)
	t, err := client.GetUsers(org)
	if err != nil {
		d.SetId("")
		return nil
	}
	d.Set("user", user)
	d.Set("role", t[user])
	d.Set("org", org)
	return nil
}

func resourceNPMUserDelete(d *schema.ResourceData, m interface{}) error {
	user := d.Get("user").(string)
	org := d.Get("org").(string)
	client := m.(*npm.Client)
	err := client.DeleteUser(org, user)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceNPMUserImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	sParts := strings.Split(d.Id(), ":")

	if len(sParts) != 3 {
		return nil, fmt.Errorf("Invalid ID specified. Supplied ID must be written as user:role:org")
	}

	d.Set("user", sParts[0])
	d.Set("role", sParts[1])
	d.Set("org", sParts[2])

	d.SetId(sParts[0])

	return []*schema.ResourceData{d}, nil
}
