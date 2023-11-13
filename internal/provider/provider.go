package provider

import (
	"fmt"

	"cd-terraform-provider/internal/console"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func PluralCDProvider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"console_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PLURALCD_CONSOLE_URL", nil),
				Description: "Plural Console URL, i.e. https://console.demo.onplural.sh.",
			},
			"access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PLURALCD_ACCESS_TOKEN", nil),
				Description: "Plural Console access token.",
			},
			"use_cli": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				DefaultFunc: schema.EnvDefaultFunc("PLURALCD_USE_CLI", nil),
				Description: "Use `plural cd login` command for authentication.",
			},
		},
		ResourcesMap:   map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (any, error) {
		url := d.Get("console_url").(string)
		token := d.Get("access_token").(string)
		useCli := d.Get("use_cli").(bool)

		if useCli {
			config := console.ReadConfig()

			token = config.Token
			if token == "" {
				return nil, fmt.Errorf("you have not set up a console login, run `plural cd login` to save your credentials")
			}

			url = config.Url
			if config.Url == "" {
				return nil, fmt.Errorf("you have not set up a console login, run `plural cd login` to save your credentials")
			}

		}

		return console.NewClient(url, token)
	}

	return provider
}
