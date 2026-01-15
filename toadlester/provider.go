package toadlester

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TOADLESTER_URL", nil),
				Description: "URL to access the Toadlester API",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"toadlester": resourceType(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"toadlester_type": dataSourceType(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diag diag.Diagnostics
	baseURL := d.Get("base_url").(string)

	client := NewAPIClient(baseURL)

	return &Config{
		Client: client,
	}, diag
}
