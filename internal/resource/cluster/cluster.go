package cluster

import (
	"context"
	"fmt"
	"time"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"
	internalvalidator "terraform-provider-plural/internal/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
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
		Description:         "A representation of a cluster you can deploy to.",
		MarkdownDescription: "A representation of a cluster you can deploy to.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this cluster.",
				MarkdownDescription: "Internal identifier of this cluster.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"inserted_at": schema.StringAttribute{
				Description:         "Creation date of this cluster.",
				MarkdownDescription: "Creation date of this cluster.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable name of this cluster, that also translates to cloud resource name.",
				MarkdownDescription: "Human-readable name of this cluster, that also translates to cloud resource name.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"handle": schema.StringAttribute{
				Description:         "A short, unique human-readable name used to identify this cluster. Does not necessarily map to the cloud resource name.",
				MarkdownDescription: "A short, unique human-readable name used to identify this cluster. Does not necessarily map to the cloud resource name.",
				Optional:            true,
				Computed:            true,
			},
			"version": schema.StringAttribute{
				Description:         "Kubernetes version to use for this cluster. Leave empty for bring your own cluster. Supported version ranges can be found at https://github.com/pluralsh/console/tree/master/static/k8s-versions.",
				MarkdownDescription: "Kubernetes version to use for this cluster. Leave empty for bring your own cluster. Supported version ranges can be found at https://github.com/pluralsh/console/tree/master/static/k8s-versions.",
				Optional:            true,
				Validators: []validator.String{
					internalvalidator.ConflictsWithIf(internalvalidator.ConflictsIfTargetValueOneOf([]string{common.CloudBYOK.String()}),
						path.MatchRoot("cloud")),
				},
			},
			"desired_version": schema.StringAttribute{
				Description:         "Desired Kubernetes version for this cluster.",
				MarkdownDescription: "Desired Kubernetes version for this cluster.",
				Computed:            true,
			},
			"current_version": schema.StringAttribute{
				Description:         "Current Kubernetes version for this cluster.",
				MarkdownDescription: "Current Kubernetes version for this cluster.",
				Computed:            true,
			},
			"provider_id": schema.StringAttribute{
				Description:         "Provider used to create this cluster. Leave empty for bring your own cluster.",
				MarkdownDescription: "Provider used to create this cluster. Leave empty for bring your own cluster.",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					internalvalidator.ConflictsWithIf(internalvalidator.ConflictsIfTargetValueOneOf([]string{common.CloudBYOK.String()}),
						path.MatchRoot("cloud")),
				},
			},
			"cloud": schema.StringAttribute{
				Description:         "The cloud provider used to create this cluster.",
				MarkdownDescription: "The cloud provider used to create this cluster.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
					common.CloudBYOK.String(), common.CloudAWS.String(), common.CloudAzure.String(), common.CloudGCP.String()),
					internalvalidator.AlsoRequiresIf(internalvalidator.RequiresIfSourceValueOneOf([]string{
						common.CloudAWS.String(),
						common.CloudAzure.String(),
						common.CloudGCP.String(),
					}), path.MatchRoot("provider_id")),
				},
			},
			"cloud_settings": schema.SingleNestedAttribute{
				Description:         "Cloud-specific settings for this cluster.",
				MarkdownDescription: "Cloud-specific settings for this cluster.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"aws":   AWSCloudSettingsSchema(),
					"azure": AzureCloudSettingsSchema(),
					"gcp":   GCPCloudSettingsSchema(),
					"byok":  BYOKCloudSettingsSchema(),
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
			},
			"node_pools": schema.SetNestedAttribute{
				Description:         "Experimental, not ready for production use. List of node pool specs managed by this cluster. Leave empty for bring your own cluster.",
				MarkdownDescription: "**Experimental, not ready for production use.** List of node pool specs managed by this cluster. Leave empty for bring your own cluster.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "Node pool name. Must be unique.",
							MarkdownDescription: "Node pool name. Must be unique.",
							Required:            true,
						},
						"min_size": schema.Int64Attribute{
							Description:         "Minimum number of instances in this node pool.",
							MarkdownDescription: "Minimum number of instances in this node pool.",
							Required:            true,
						},
						"max_size": schema.Int64Attribute{
							Description:         "Maximum number of instances in this node pool.",
							MarkdownDescription: "Maximum number of instances in this node pool.",
							Required:            true,
						},
						"instance_type": schema.StringAttribute{
							Description:         "The type of node to use. Usually cloud-specific.",
							MarkdownDescription: "The type of node to use. Usually cloud-specific.",
							Required:            true,
						},
						"labels": schema.MapAttribute{
							Description:         "Kubernetes labels to apply to the nodes in this pool. Useful for node selectors.",
							MarkdownDescription: "Kubernetes labels to apply to the nodes in this pool. Useful for node selectors.",
							ElementType:         types.StringType,
							Optional:            true,
							Computed:            true,
						},
						"taints": schema.SetNestedAttribute{
							Description:         "Any taints you'd want to apply to a node, i.e. for preventing scheduling on spot instances.",
							MarkdownDescription: "Any taints you'd want to apply to a node, i.e. for preventing scheduling on spot instances.",
							Optional:            true,
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Required: true,
									},
									"value": schema.StringAttribute{
										Required: true,
									},
									"effect": schema.StringAttribute{
										Required: true,
									},
								},
							},
						},
						"cloud_settings": schema.SingleNestedAttribute{
							Description:         "Cloud-specific settings for this node pool.",
							MarkdownDescription: "Cloud-specific settings for this node pool.",
							Optional:            true,
							PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
							Attributes: map[string]schema.Attribute{
								"aws": schema.SingleNestedAttribute{
									Description:         "AWS node pool customizations.",
									MarkdownDescription: "AWS node pool customizations.",
									Optional:            true,
									PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
									Attributes: map[string]schema.Attribute{
										"launch_template_id": schema.StringAttribute{
											Description:         "Custom launch template for your nodes. Useful for Golden AMI setups.",
											MarkdownDescription: "Custom launch template for your nodes. Useful for Golden AMI setups.",
											Optional:            true,
											PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
										},
									},
								},
							},
						},
					},
				},
			},
			"protect": schema.BoolAttribute{
				Description:         "If set to \"true\" then this cluster cannot be deleted.",
				MarkdownDescription: "If set to `true` then this cluster cannot be deleted.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"tags": schema.MapAttribute{
				Description:         "Key-value tags used to filter clusters.",
				MarkdownDescription: "Key-value tags used to filter clusters.",
				Optional:            true,
				ElementType:         types.StringType,
				PlanModifiers:       []planmodifier.Map{mapplanmodifier.RequiresReplace()},
			},
			"bindings": schema.SingleNestedAttribute{
				Description:         "Read and write policies of this cluster.",
				MarkdownDescription: "Read and write policies of this cluster.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"read": schema.SetNestedAttribute{
						Optional:            true,
						Description:         "Read policies of this cluster.",
						MarkdownDescription: "Read policies of this cluster.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"group_id": schema.StringAttribute{
									Description:         "",
									MarkdownDescription: "",
									Optional:            true,
								},
								"id": schema.StringAttribute{
									Description:         "",
									MarkdownDescription: "",
									Optional:            true,
								},
								"user_id": schema.StringAttribute{
									Description:         "",
									MarkdownDescription: "",
									Optional:            true,
								},
							},
						},
					},
					"write": schema.SetNestedAttribute{
						Optional:            true,
						Description:         "Write policies of this cluster.",
						MarkdownDescription: "Write policies of this cluster.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"group_id": schema.StringAttribute{
									Description:         "",
									MarkdownDescription: "",
									Optional:            true,
								},
								"id": schema.StringAttribute{
									Description:         "",
									MarkdownDescription: "",
									Optional:            true,
								},
								"user_id": schema.StringAttribute{
									Description:         "",
									MarkdownDescription: "",
									Optional:            true,
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *clusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = data.Client
	r.consoleUrl = data.ConsoleUrl
}

func (r *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data cluster
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateCluster(ctx, data.Attributes(ctx, resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create cluster, got error: %s", err))
		return
	}

	if common.IsCloud(data.Cloud.ValueString(), common.CloudBYOK) {
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

	data.FromCreate(result, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data cluster
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetCluster(ctx, data.Id.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster, got error: %s", err))
		return
	}
	if result == nil || result.Cluster == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Unable to find cluster, it looks like it was deleted manually"))
		return
	}

	data.From(result.Cluster, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data cluster
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateCluster(ctx, data.Id.ValueString(), data.UpdateAttributes(ctx, resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update cluster, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data cluster
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
		response, err := r.client.GetCluster(ctx, data.Id.ValueStringPointer())
		if client.IsNotFound(err) || response.Cluster == nil {
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
