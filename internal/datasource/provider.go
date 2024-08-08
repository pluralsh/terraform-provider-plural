package datasource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
				Optional:            true,
				Computed:            true,
				Description:         "Internal identifier of this provider.",
				MarkdownDescription: "Internal identifier of this provider.",
				Validators:          []validator.String{stringvalidator.ExactlyOneOf(path.MatchRoot("cloud"))},
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
				Optional:            true,
				Computed:            true,
				Validators:          []validator.String{stringvalidator.ExactlyOneOf(path.MatchRoot("id"))},
			},
		},
	}
}

func (p *providerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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

	var id string
	if !data.Id.IsNull() {
		id = data.Id.ValueString()
	} else {
		if providers, err := p.client.ListProviders(ctx); err == nil {
			for _, pv := range providers.ClusterProviders.Edges {
				if pv.Node.Cloud == data.Cloud.ValueString() {
					id = pv.Node.ID
					break
				}
			}
		} else {
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to get providers, got error: %s", err))
		}
	}
	if id == "" {
		resp.Diagnostics.AddError("Client Error", "Unable to determine provider ID")
		return
	}

	result, err := p.client.GetClusterProvider(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read provider, got error: %s", err))
		return
	}
	if result == nil && result.ClusterProvider == nil {
		resp.Diagnostics.AddError("Not Found", "Unable to find provider")
		return
	}

	data.From(result.ClusterProvider)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
