package resource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"
)

var _ resource.Resource = &rbacResource{}
var _ resource.ResourceWithImportState = &rbacResource{}

func NewRbacResource() resource.Resource {
	return &rbacResource{}
}

// rbacResource defines the rbac resource implementation.
type rbacResource struct {
	client     *client.Client
	consoleUrl string
}

func (r *rbacResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rbac"
}

func (r *rbacResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A representation of rbac settings for a provider or cluster.",
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Description:         "The cluster id for these rbac settings",
				MarkdownDescription: "The cluster id for these rbac settings",
				Optional:            true,
			},
			"service_id": schema.StringAttribute{
				Description:         "The service id for these rbac settings",
				MarkdownDescription: "The service id for these rbac settings",
				Optional:            true,
			},
			"bindings": schema.SingleNestedAttribute{
				Description:         "Read and write policies of this resource.",
				MarkdownDescription: "Read and write policies of this resource.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"read": schema.SetNestedAttribute{
						Optional:            true,
						Description:         "Read policies of this resource.",
						MarkdownDescription: "Read policies of this resource.",
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
						Description:         "Write policies of this resource.",
						MarkdownDescription: "Write policies of this resource.",
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
			},
		},
	}
}

func (r *rbacResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = data.Client
	r.consoleUrl = data.ConsoleUrl
}

func (r *rbacResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.RBAC
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateRbac(ctx, data.Attributes(ctx, resp.Diagnostics), data.ServiceId.ValueStringPointer(), data.ClusterId.ValueStringPointer(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update rbac, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rbacResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// ignore
}

func (r *rbacResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data model.RBAC
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateRbac(ctx, data.Attributes(ctx, resp.Diagnostics), data.ServiceId.ValueStringPointer(), data.ClusterId.ValueStringPointer(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update rbac, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rbacResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// ignore
}

func (r *rbacResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// ignore
}
