package npm

import (
	"github.com/hashicorp/terraform/helper/schema"
  "github.com/hashicorp/terraform/terraform"
  "github.com/alchemistake/go-npm"
)

//Provider Name
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NPM_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"npm_membership": resourceNPMUser(),
		},
    	ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
client := npm.NewTokenClient(d.Get("token").(string))
	return client, nil
}
