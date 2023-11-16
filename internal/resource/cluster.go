package resource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	console "github.com/pluralsh/console-client-go"
)

var _ resource.Resource = &clusterResource{}
var _ resource.ResourceWithImportState = &clusterResource{}

func NewClusterResource() resource.Resource {
	return &clusterResource{}
}

// ClusterResource defines the cluster resource implementation.
type clusterResource struct {
	client *client.Client
}

func (r *clusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (r *clusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A representation of a cluster you can deploy to.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal identifier of this cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"inserted_at": schema.StringAttribute{
				MarkdownDescription: "Creation date of this cluster.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name of this cluster, that also translates to cloud resource name.",
				Required:            true,
			},
			"handle": schema.StringAttribute{
				MarkdownDescription: "A short, unique human-readable name used to identify this cluster. Does not necessarily map to the cloud resource name.",
				Optional:            true,
				Computed:            true,
			},
			"cloud": schema.StringAttribute{
				MarkdownDescription: "The cloud provider used to create this cluster.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOfCaseInsensitive("byok")},
			},
			"protect": schema.BoolAttribute{
				MarkdownDescription: "If set to `true` then this cluster cannot be deleted.",
				Optional:            true,
			},
			"tags": schema.MapAttribute{
				MarkdownDescription: "Key-value tags used to filter clusters.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *clusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Cluster Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.Cluster
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := console.ClusterAttributes{
		Name:    data.Name.ValueString(),
		Handle:  data.Handle.ValueStringPointer(),
		Protect: data.Protect.ValueBoolPointer(),
	}
	cluster, err := r.client.CreateCluster(ctx, attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create cluster, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a cluster")

	data.Name = types.StringValue(cluster.CreateCluster.Name)
	data.Handle = types.StringPointerValue(cluster.CreateCluster.Handle)
	data.Id = types.StringValue(cluster.CreateCluster.ID)
	data.Protect = types.BoolPointerValue(cluster.CreateCluster.Protect)
	data.InseredAt = types.StringPointerValue(cluster.CreateCluster.InsertedAt)

	if data.Cloud.ValueString() == "byok" {
		if cluster.CreateCluster.DeployToken == nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to fetch cluster deploy token"))
			return
		}

		// TODO:
		//   deployToken := *cluster.CreateCluster.DeployToken
		//   url := fmt.Sprintf("%s/ext/gql", p.ConsoleClient.Url())
		//   p.doInstallOperator(url, deployToken)

		tflog.Trace(ctx, "installed the cluster operator")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.Cluster
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cluster, err := r.client.GetCluster(ctx, data.Id.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster, got error: %s", err))
		return
	}

	data.Id = types.StringValue(cluster.Cluster.ID)
	data.InseredAt = types.StringPointerValue(cluster.Cluster.InsertedAt)
	data.Name = types.StringValue(cluster.Cluster.Name)
	data.Handle = types.StringPointerValue(cluster.Cluster.Handle)
	data.Protect = types.BoolPointerValue(cluster.Cluster.Protect)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data model.Cluster
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := console.ClusterUpdateAttributes{
		Handle:  data.Handle.ValueStringPointer(),
		Protect: data.Protect.ValueBoolPointer(),
	}
	cluster, err := r.client.UpdateCluster(ctx, data.Id.ValueString(), attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update cluster, got error: %s", err))
		return
	}

	data.Handle = types.StringPointerValue(cluster.UpdateCluster.Handle)
	data.Protect = types.BoolPointerValue(cluster.UpdateCluster.Protect)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data model.Cluster
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteCluster(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete cluster, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted the cluster")
}

func (r *clusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
