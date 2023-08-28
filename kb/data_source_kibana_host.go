// Return the connection settings of Kibana
// Supported version:
//  - v7

package kb

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kibana "github.com/samuel-deal-tisseo/go-kibana-rest/v7"
)

func dataSourceKibanaHost() *schema.Resource {
	return &schema.Resource{
		Description: "`kibana_host` can be used to retrieve the Kibana connection settings.",
		Read:        dataSourceKibanaHostRead,

		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kibana URL",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username to use to connect to Kibana using basic auth",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Password to use to connect to Kibana using basic auth",
			},
			"method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "basic or form",
			},
		},
	}
}

func dataSourceKibanaHostRead(d *schema.ResourceData, m interface{}) error {
	var url string
	var username string
	var password string

	conf := m.(*kibana.Client)

	url = conf.Client.HostURL
	username = conf.Client.UserInfo.Username
	password = conf.Client.UserInfo.Password

	d.SetId(url)
	d.Set("url", url)
	d.Set("username", username)
	d.Set("password", password)
	if conf.AuthMethod == kibana.AUTH_FORM {
		d.Set("method", "form")
	} else {
		d.Set("method", "basic")
	}
	return nil
}
