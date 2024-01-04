package datasource

import (
	"context"
	"fmt"
	"os"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"

	"github.com/mitchellh/go-homedir"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"gopkg.in/yaml.v2"
)

type config struct {
	Email types.String `tfsdk:"email" yaml:"email"`
	Token types.String `tfsdk:"token" yaml:"email"`
}

func NewConfigDataSource() datasource.DataSource {
	return &configDataSource{}
}

type configDataSource struct {
	client *client.Client
}

func (d *configDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (d *configDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A representation of a config to authenticate to app.plural.sh",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The email used to authenticate to plural.",
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Access token used to authenticate to plural.",
				Computed:            true,
			},
		},
	}
}

func (d *configDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Config Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = data.Client
}

func (d *configDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data config
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p, err := homedir.Expand("~/.plural/config.yml")
	if err != nil {
		resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Could not read local plural config: %s", err))
		return
	}

	res, err := os.ReadFile(p)
	if err != nil {
		resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Could not read local plural config: %s", err))
		return
	}

	var conf struct {
		Spec config
	}

	if err := yaml.Unmarshal(res, &conf); err != nil {
		resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Could not parse local plural config: %s", err))
		return
	}

	data = conf.Spec
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
