package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	consoleClient "github.com/pluralsh/console-client-go"
)

var _ resource.Resource = &ClusterResource{}
var _ resource.ResourceWithImportState = &ClusterResource{}

func NewClusterResource() resource.Resource {
	return &ClusterResource{}
}

// ClusterResource defines the cluster resource implementation.
type ClusterResource struct {
	client *consoleClient.Client
}

// ClusterResourceModel describes the cluster resource data model.
type ClusterResourceModel struct {
	Id        types.String `tfsdk:"id"`
	InseredAt types.String `tfsdk:"inserted_at"`
	Name      types.String `tfsdk:"name"`
	Handle    types.String `tfsdk:"handle"`
	Cloud     types.String `tfsdk:"cloud"`
	Protect   types.Bool   `tfsdk:"protect"`
	Tags      types.Map    `tfsdk:"tags"`
}

func (r *ClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (r *ClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *ClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*consoleClient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Cluster Resource Configure Type",
			fmt.Sprintf("Expected *console.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := consoleClient.ClusterAttributes{
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

	data.Id = types.StringValue(cluster.CreateCluster.ID)

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

func (r *ClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClusterResourceModel
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
	data.InseredAt = types.StringUnknown() // TODO: Update client to return this field.
	data.Name = types.StringValue(cluster.Cluster.Name)
	data.Handle = types.StringValue(*cluster.Cluster.Handle)
	data.Protect = types.BoolUnknown() // TODO: Update client to return this field.

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := consoleClient.ClusterUpdateAttributes{
		Handle: data.Handle.ValueStringPointer(),
	}
	cluster, err := r.client.UpdateCluster(ctx, data.Id.ValueString(), attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update cluster, got error: %s", err))
		return
	}

	data.Handle = types.StringValue(*cluster.UpdateCluster.Handle)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClusterResourceModel
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

func (r *ClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
