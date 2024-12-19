package datasource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func NewServiceContextDataSource() datasource.DataSource {
	return &serviceContextDataSource{}
}

type serviceContextDataSource struct {
	client *client.Client
}

func (d *serviceContextDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_context"
}

func (d *serviceContextDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A representation of a service context that can be reused during service deployment templating.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Internal identifier of this service context.",
				MarkdownDescription: "Internal identifier of this service context.",
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable name of this service context.",
				MarkdownDescription: "Human-readable name of this service context.",
				Required:            true,
			},
			"configuration": schema.StringAttribute{
				Description:         "Configuration in JSON format. Use 'jsondecode' method to decode data.",
				MarkdownDescription: "Configuration in JSON format. Use `jsondecode` method to decode data.",
				Computed:            true,
			},
		},
	}
}

func (d *serviceContextDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Service Context Data Source Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = data.Client
}

func (d *serviceContextDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := new(model.ServiceContext)
	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.client.GetServiceContext(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service context, got error: %s", err))
		return
	}

	data.From(response.ServiceContext, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
