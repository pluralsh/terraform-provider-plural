package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/model"
)

var _ resource.Resource = &ServiceDeploymentResource{}
var _ resource.ResourceWithImportState = &ServiceDeploymentResource{}

func NewServiceDeploymentResource() resource.Resource {
	return &ServiceDeploymentResource{}
}

// ServiceDeploymentResource defines the ServiceDeployment resource implementation.
type ServiceDeploymentResource struct {
	client *client.Client
}

func (r *ServiceDeploymentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_deployment"
}

func (r *ServiceDeploymentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "ServiceDeployment resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal identifier of this ServiceDeployment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name of this ServiceDeployment.",
				Required:            true,
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "Namespace to deploy this ServiceDeployment.",
				Required:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Semver version of this service ServiceDeployment.",
				Optional:            true,
			},
			"docs_path": schema.StringAttribute{
				MarkdownDescription: "Path to the documentation in the target git repository.",
				Optional:            true,
			},
			"protect": schema.BoolAttribute{
				MarkdownDescription: "If true, deletion of this service is not allowed.",
				Optional:            true,
			},
			"kustomize": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"path": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Path to the kustomize file in the target git repository.",
					},
				},
				MarkdownDescription: "Kustomize related service metadata.",
				Optional:            true,
			},
			"configuration": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required:  true,
							Sensitive: true,
						},
					},
				},
				MarkdownDescription: "List of [name, value] secrets used to alter this ServiceDeployment configuration.",
				Optional:            true,
			},
			"cluster": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRoot("cluster").AtName("handle")),
							stringvalidator.ExactlyOneOf(path.MatchRoot("cluster").AtName("handle")),
						},
						Optional: true,
					},
					"handle": schema.StringAttribute{
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRoot("cluster").AtName("id")),
							stringvalidator.ExactlyOneOf(path.MatchRoot("cluster").AtName("id")),
						},
						Optional: true,
					},
				},
				MarkdownDescription: "Unique cluster id/handle to deploy this ServiceDeployment",
				Required:            true,
			},
			"repository": schema.ObjectAttribute{
				AttributeTypes: map[string]attr.Type{
					"id":     types.StringType,
					"ref":    types.StringType,
					"folder": types.StringType,
				},
				MarkdownDescription: "Repository information used to pull ServiceDeployment.",
				Required:            true,
			},
			"bindings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"read": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"group_id": schema.StringAttribute{
									Optional: true,
								},
								"id": schema.StringAttribute{
									Optional: true,
								},
								"user_id": schema.StringAttribute{
									Optional: true,
								},
							},
						},
						MarkdownDescription: "Read policies of this ServiceDeployment.",
						Optional:            true,
					},
					"write": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"group_id": schema.StringAttribute{
									Optional: true,
								},
								"id": schema.StringAttribute{
									Optional: true,
								},
								"user_id": schema.StringAttribute{
									Optional: true,
								},
							},
						},
						MarkdownDescription: "Write policies of this ServiceDeployment.",
						Optional:            true,
					},
				},
				MarkdownDescription: "Read and write policies of this ServiceDeployment.",
				Optional:            true,
			},
			"sync_config": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"diff_normalizer": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"group": schema.SingleNestedAttribute{
								Optional: true,
							},
							"json_patches": schema.SingleNestedAttribute{
								Optional: true,
							},
							"kind": schema.SingleNestedAttribute{
								Optional: true,
							},
							"name": schema.SingleNestedAttribute{
								Optional: true,
							},
							"namespace": schema.SingleNestedAttribute{
								Optional: true,
							},
						},
						Optional: true,
					},
					"namespace_metadata": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"annotations": schema.MapAttribute{
								ElementType: types.StringType,
								Optional:    true,
							},
							"labels": schema.MapAttribute{
								ElementType: types.StringType,
								Optional:    true,
							},
						},
						Optional: true,
					},
				},
				MarkdownDescription: "Settings for advanced tuning of the sync process.",
				Optional:            true,
			},
		},
	}
}

func (r *ServiceDeploymentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*model.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ServiceDeployment Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *ServiceDeploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := new(model.ServiceDeployment)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("service attributes: %+v", data.Attributes()))

	sd, err := r.client.CreateServiceDeployment(ctx, data.Cluster.Id.ValueStringPointer(), data.Cluster.Handle.ValueStringPointer(), data.Attributes())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create ServiceDeployment, got error: %s", err))
		return
	}

	data.Id = types.StringValue(sd.ID)

	tflog.Trace(ctx, "created a ServiceDeployment")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceDeploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := new(model.ServiceDeployment)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetServiceDeployment(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ServiceDeployment, got error: %s", err))
		return
	}

	data.From(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceDeploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(model.ServiceDeployment)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("service deployment: %+v", data))

	_, err := r.client.UpdateServiceDeployment(ctx, data.Id.ValueString(), data.UpdateAttributes())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update ServiceDeployment, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceDeploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(model.ServiceDeployment)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteServiceDeployment(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete ServiceDeployment, got error: %s", err))
		return
	}
}

func (r *ServiceDeploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
