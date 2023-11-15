package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	ds "terraform-provider-plural/internal/datasource"
	r "terraform-provider-plural/internal/resource"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pluralsh/console-client-go"
	"github.com/pluralsh/plural-cli/pkg/console"
)

var _ provider.Provider = &PluralProvider{}

// PluralProvider defines the Plural provider implementation.
type PluralProvider struct {
	version string
}

// pluralProviderModel describes the Plural provider data model.
type pluralProviderModel struct {
	ConsoleUrl  types.String `tfsdk:"console_url"`
	AccessToken types.String `tfsdk:"access_token"`
	UseCli      types.Bool   `tfsdk:"use_cli"`
}

func (p *PluralProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "plural"
	resp.Version = p.version
}

func (p *PluralProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"console_url": schema.StringAttribute{
				MarkdownDescription: "Plural Console URL, i.e. `https://console.demo.onplural.sh`. Can be sourced from `PLURAL_CONSOLE_URL`.",
				Optional:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Plural Console access token. Can be sourced from `PLURAL_ACCESS_TOKEN`.",
				Optional:            true,
				Sensitive:           true,
			},
			"use_cli": schema.BoolAttribute{
				MarkdownDescription: "Use Plural CLI `plural cd login` command for authentication. Can be sourced from `PLURAL_USE_CLI`.",
				Optional:            true,
			},
		},
	}
}

func (p *PluralProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	consoleUrl := os.Getenv("PLURAL_CONSOLE_URL")
	accessToken := os.Getenv("PLURAL_ACCESS_TOKEN")
	useCli, _ := strconv.ParseBool(os.Getenv("PLURAL_USE_CLI"))

	var data pluralProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.ConsoleUrl.IsNull() {
		consoleUrl = data.ConsoleUrl.ValueString()
	}

	if !data.AccessToken.IsNull() {
		accessToken = data.AccessToken.ValueString()
	}

	if !data.UseCli.IsNull() {
		useCli = data.UseCli.ValueBool()
	}

	if useCli {
		config := console.ReadConfig()
		accessToken = config.Token
		consoleUrl = config.Url

		if consoleUrl == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("use_cli"),
				"Missing Plural Console URL",
				"The provider could not read Plural Console URL from Plural CLI. "+
					"Run `plural cd login` to save your credentials first. "+
					"You can also specify Plural Console URL and access token directly, see documentation for more information.",
			)
		}

		if accessToken == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("use_cli"),
				"Missing Plural Access Token",
				"The provider could not read Plural Console access token from Plural CLI. "+
					"Run `plural cd login` to save your credentials first. "+
					"You can also specify Plural Console URL and access token directly, see documentation for more information.",
			)
		}
	} else {
		if consoleUrl == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("console_url"),
				"Missing Plural Console URL",
				"The provider cannot create the Plural Console client as there is a missing or empty value for the Plural Console URL. "+
					"Set the URL value in the configuration or use the PLURAL_CONSOLE_URL environment variable. "+
					"If either is already set, ensure the value is not empty. "+
					"You can also use Plural CLI for authentication, see documentation for more information.",
			)
		}

		if accessToken == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("access_token"),
				"Missing Plural Console Access Token",
				"The provider cannot create the Plural Console client as there is a missing or empty value for the Plural Console access token. "+
					"Set the URL value in the configuration or use the PLURAL_ACCESS_TOKEN environment variable. "+
					"If either is already set, ensure the value is not empty. "+
					"You can also use Plural CLI for authentication, see documentation for more information.",
			)
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	consoleClient := client.NewClient(http.DefaultClient, fmt.Sprintf("%s/gql", consoleUrl), func(req *http.Request) {
		req.Header.Set("Authorization", fmt.Sprintf("Token %s", accessToken))
	})

	resp.ResourceData = consoleClient
	resp.DataSourceData = consoleClient
}

func (p *PluralProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		r.NewClusterResource,
	}
}

func (p *PluralProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		ds.NewClusterDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PluralProvider{
			version: version,
		}
	}
}