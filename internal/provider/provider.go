package provider

import (
	"context"

	"terraform-provider-plural/internal/console"
	"terraform-provider-plural/internal/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func PluralProvider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"console_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PLURAL_CONSOLE_URL", nil),
				Description: "Plural Console URL, i.e. https://console.demo.onplural.sh.",
			},
			"access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PLURAL_ACCESS_TOKEN", nil),
				Description: "Plural Console access token.",
			},
			"use_cli": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				DefaultFunc: schema.EnvDefaultFunc("PLURAL_USE_CLI", nil),
				Description: "Use `plural cd login` command for authentication.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"plural_cluster": resource.Cluster(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		url := d.Get("console_url").(string)
		token := d.Get("access_token").(string)
		useCli := d.Get("use_cli").(bool)

		if useCli {
			config := console.ReadConfig()

			token = config.Token
			if token == "" {
				return nil, diag.Errorf("you have not set up a console login, run `plural cd login` to save your credentials")
			}

			url = config.Url
			if config.Url == "" {
				return nil, diag.Errorf("you have not set up a console login, run `plural cd login` to save your credentials")
			}
		}

		return console.NewClient(url, token), nil
	}

	return provider
}
