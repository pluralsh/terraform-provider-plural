package datasource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewProviderDataSource() datasource.DataSource {
	return &providerDataSource{}
}

type providerDataSource struct {
	client *client.Client
}

func (p *providerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider"
}

func (p *providerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A representation of a provider you can deploy your clusters to.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "Internal identifier of this provider.",
				MarkdownDescription: "Internal identifier of this provider.",
			},
			"editable": schema.BoolAttribute{
				Description:         "Whether this provider is editable.",
				MarkdownDescription: "Whether this provider is editable.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable name of this provider. Globally unique.",
				MarkdownDescription: "Human-readable name of this provider. Globally unique.",
				Computed:            true,
			},
			"namespace": schema.StringAttribute{
				Description:         "The namespace the Cluster API resources are deployed into.",
				MarkdownDescription: "The namespace the Cluster API resources are deployed into.",
				Computed:            true,
			},
			"cloud": schema.StringAttribute{
				Description:         "The name of the cloud service for this provider.",
				MarkdownDescription: "The name of the cloud service for this provider.",
				Computed:            true,
			},
		},
	}
}

func (p *providerDataSource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*model.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Resource Configure Type",
			fmt.Sprintf("Expected *model.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	p.client = data.Client
}

func (p *providerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model.Provider
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Provider ID",
			"The provider could not read provider data. ID needs to be specified.",
		)
	}

	result, err := p.client.GetClusterProvider(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read provider, got error: %s", err))
		return
	}
	if result == nil || result.ClusterProvider == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Unable to find provider"))
		return
	}

	data.From(result.ClusterProvider)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
