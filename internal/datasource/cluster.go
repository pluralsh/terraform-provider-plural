package datasource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/model"
	"terraform-provider-plural/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console-client-go"
)

func NewClusterDataSource() datasource.DataSource {
	return &clusterDataSource{}
}

type clusterDataSource struct {
	client *client.Client
}

func (d *clusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (d *clusterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A representation of a cluster you can deploy to.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Internal identifier of this cluster.",
			},
			"inserted_at": schema.StringAttribute{
				MarkdownDescription: "Creation date of this cluster.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name of this cluster, that also translates to cloud resource name.",
				Computed:            true,
			},
			"handle": schema.StringAttribute{
				MarkdownDescription: "A short, unique human-readable name used to identify this cluster. Does not necessarily map to the cloud resource name.",
				Optional:            true,
				Computed:            true,
			},
			"cloud": schema.StringAttribute{
				MarkdownDescription: "The cloud provider used to create this cluster.",
				Computed:            true,
			},
			"protect": schema.BoolAttribute{
				MarkdownDescription: "If set to `true` then this cluster cannot be deleted.",
				Computed:            true,
			},
			"tags": schema.MapAttribute{
				MarkdownDescription: "Key-value tags used to filter clusters.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *clusterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*provider.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Cluster Resource Configure Type",
			fmt.Sprintf("Expected *provider.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = data.Client
}

func (d *clusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model.Cluster
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.IsNull() && data.Handle.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Cluster ID and Handle",
			"The provider could not read cluster data. ID or handle needs to be specified.",
		)
	}

	// First try to fetch cluster by ID if it was provided.
	var cluster *console.ClusterFragment
	if !data.Id.IsNull() {
		if c, err := d.client.GetCluster(ctx, data.Id.ValueStringPointer()); err != nil {
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to read cluster by ID, got error: %s", err))
		} else {
			cluster = c.Cluster
		}
	}

	// If cluster was not fetched yet and handle was provided then try to use it to fetch cluster.
	if cluster == nil && !data.Handle.IsNull() {
		if c, err := d.client.GetClusterByHandle(ctx, data.Handle.ValueStringPointer()); err != nil {
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to read cluster by handle, got error: %s", err))
		} else {
			cluster = c.Cluster
		}
	}

	if cluster == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster, see warnings for more information"))
		return
	}

	data.Id = types.StringValue(cluster.ID)
	data.InseredAt = types.StringPointerValue(cluster.InsertedAt)
	data.Name = types.StringValue(cluster.Name)
	data.Handle = types.StringPointerValue(cluster.Handle)
	data.Protect = types.BoolPointerValue(cluster.Protect)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
