package resource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/model"

	"terraform-provider-plural/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ServiceContextResource{}
var _ resource.ResourceWithImportState = &ServiceContextResource{}

func NewServiceContextResource() resource.Resource {
	return &ServiceContextResource{}
}

// ServiceContextResource defines the serviceContext resource implementation.
type ServiceContextResource struct {
	client *client.Client
}

func (r *ServiceContextResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_context"
}

func (r *ServiceContextResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A representation of a service context that can be reused during service deployment templating.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Internal identifier of this service context.",
				MarkdownDescription: "Internal identifier of this service context.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable name of this service context. Globally unique.",
				MarkdownDescription: "Human-readable name of this service context. Globally unique.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"configuration": schema.StringAttribute{
				Description:         "Configuration in JSON format. Use 'jsonencode' and 'jsondecode' methods to encode and decode data.",
				MarkdownDescription: "Configuration in JSON format. Use `jsonencode` and `jsondecode` methods to encode and decode data.",
				Optional:            true,
			},
			"secrets": schema.MapAttribute{
				Description:         "",
				MarkdownDescription: "",
				Optional:            true,
				Sensitive:           true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *ServiceContextResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Service context Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *ServiceContextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := new(model.ServiceContextExtended)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sc, err := r.client.SaveServiceContext(ctx, data.Name.ValueString(), data.Attributes(ctx, &resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service context, got error: %s", err))
		return
	}

	data.From(sc.SaveServiceContext, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *ServiceContextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := new(model.ServiceContextExtended)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetServiceContext(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service context, got error: %s", err))
		return
	}
	if response.ServiceContext == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to find service context, got no error")
		return
	}

	data.From(response.ServiceContext, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *ServiceContextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(model.ServiceContextExtended)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.SaveServiceContext(ctx, data.Name.ValueString(), data.Attributes(ctx, &resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update service context, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *ServiceContextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(model.ServiceContextExtended)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteServiceContext(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service context, got error: %s", err))
		return
	}
}

func (r *ServiceContextResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
