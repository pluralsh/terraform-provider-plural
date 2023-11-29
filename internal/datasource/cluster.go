package datasource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
				Validators:          []validator.String{stringvalidator.ExactlyOneOf(path.MatchRoot("handle"))},
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
				Validators:          []validator.String{stringvalidator.ExactlyOneOf(path.MatchRoot("id"))},
			},
			"version": schema.StringAttribute{
				Description:         "Desired Kubernetes version for this cluster.",
				MarkdownDescription: "Desired Kubernetes version for this cluster.",
				Computed:            true,
			},
			"provider_id": schema.StringAttribute{
				Description:         "Provider used to create this cluster.",
				MarkdownDescription: "Provider used to create this cluster.",
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
			"node_pools": schema.ListNestedAttribute{
				Description:         "List of node pool specs managed by this cluster.",
				MarkdownDescription: "List of node pool specs managed by this cluster.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "Node pool name.",
							MarkdownDescription: "Node pool name.",
							Computed:            true,
						},
						"min_size": schema.Int64Attribute{
							Description:         "Minimum number of instances in this node pool.",
							MarkdownDescription: "Minimum number of instances in this node pool.",
							Computed:            true,
						},
						"max_size": schema.Int64Attribute{
							Description:         "Maximum number of instances in this node pool.",
							MarkdownDescription: "Maximum number of instances in this node pool.",
							Computed:            true,
						},
						"instance_type": schema.StringAttribute{
							Description:         "The type of used node. Usually cloud-specific.",
							MarkdownDescription: "The type of used node. Usually cloud-specific.",
							Computed:            true,
						},
						"labels": schema.MapAttribute{
							Description:         "Kubernetes labels applied to the nodes in this pool.",
							MarkdownDescription: "Kubernetes labels applied to the nodes in this pool.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"taints": schema.ListNestedAttribute{
							Description:         "Taints applied to a node.",
							MarkdownDescription: "Taints applied to a node.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.MapAttribute{
										ElementType: types.StringType,
										Computed:    true,
									},
									"value": schema.MapAttribute{
										ElementType: types.StringType,
										Computed:    true,
									},
									"effect": schema.MapAttribute{
										ElementType: types.StringType,
										Computed:    true,
									},
								},
							},
						},
						"cloud_settings": schema.SingleNestedAttribute{
							Description:         "Cloud-specific settings for this node pool.",
							MarkdownDescription: "Cloud-specific settings for this node pool.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"aws": schema.SingleNestedAttribute{
									Description:         "AWS node pool customizations.",
									MarkdownDescription: "AWS node pool customizations.",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"launch_template_id": schema.StringAttribute{
											Description:         "Custom launch template for your nodes. Useful for Golden AMI setups.",
											MarkdownDescription: "Custom launch template for your nodes. Useful for Golden AMI setups.",
											Computed:            true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *clusterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Cluster Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = data.Client
}

func (d *clusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cluster
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

	data.From(cluster, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
