package cluster

import (
	"context"
	"fmt"
	"time"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

var _ resource.Resource = &clusterResource{}
var _ resource.ResourceWithImportState = &clusterResource{}

func NewClusterResource() resource.Resource {
	return &clusterResource{}
}

// ClusterResource defines the cluster resource implementation.
type clusterResource struct {
	client          *client.Client
	consoleUrl      string
	operatorHandler *OperatorHandler
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
			"version": schema.StringAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
			},
			"provider_id": schema.StringAttribute{
				MarkdownDescription: "",
				Required:            true,
			},
			"cloud": schema.StringAttribute{
				MarkdownDescription: "The cloud provider used to create this cluster.",
				Required:            true,
				Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
					model.CloudBYOK.String(), model.CloudAWS.String(), model.CloudAzure.String(), model.CloudGCP.String())},
			},
			"cloud_settings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"aws":   AWSCloudSettingsSchema(),
					"azure": AzureCloudSettingsSchema(),
					"gcp":   GCPCloudSettingsSchema(),
					"byok":  BYOKCloudSettingsSchema(),
				},
				MarkdownDescription: "Cloud-specific settings for this cluster.",
				Required:            true,
			},
			"protect": schema.BoolAttribute{
				MarkdownDescription: "If set to `true` then this cluster cannot be deleted.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
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

	data, ok := req.ProviderData.(*model.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Cluster Resource Configure Type",
			fmt.Sprintf("Expected *model.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = data.Client
	r.consoleUrl = data.ConsoleUrl
}

func (r *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.Cluster
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateCluster(ctx, data.CreateAttributes(resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create cluster, got error: %s", err))
		return
	}

	if model.IsCloud(data.Cloud.ValueString(), model.CloudBYOK) {
		if result.CreateCluster.DeployToken == nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to fetch cluster deploy token"))
			return
		}

		handler, err := NewOperatorHandler(ctx, &data.CloudSettings.BYOK.Kubeconfig, r.consoleUrl)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to init operator handler, got error: %s", err))
			return
		}

		err = handler.InstallOrUpgrade(*result.CreateCluster.DeployToken)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to install operator, got error: %s", err))
			return
		}
	}

	data.FromCreate(result, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.Cluster
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetCluster(ctx, data.Id.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster, got error: %s", err))
		return
	}

	data.From(result.Cluster, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data model.Cluster
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateCluster(ctx, data.Id.ValueString(), data.UpdateAttributes())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update cluster, got error: %s", err))
		return
	}

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

	err = wait.WaitForWithContext(ctx, client.Ticker(5*time.Second), func(ctx context.Context) (bool, error) {
		_, err := r.client.GetCluster(ctx, data.Id.ValueStringPointer())
		if client.IsNotFound(err) {
			return true, nil
		}

		return false, err
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error while watiting for cluster to be deleted, got error: %s", err))
		return
	}
}

func (r *clusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
