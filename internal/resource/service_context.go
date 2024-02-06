package resource

import (
	"context"
	"fmt"
	"time"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"k8s.io/apimachinery/pkg/util/wait"

	"terraform-provider-plural/internal/client"
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
				MarkdownDescription: "Internal identifier of this provider.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable name of this service context. Globally unique.",
				MarkdownDescription: "Human-readable name of this service context. Globally unique.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"configuration": schema.MapAttribute{
				Description:         "",
				MarkdownDescription: "",
				Optional:            true,
				ElementType:         types.StringType,
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
	data := new(serviceContext)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := data.Attributes(ctx, resp.Diagnostics)
	sc, err := r.client.SaveServiceContext(ctx, data.Name.String(), attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service context, got error: %s", err))
		return
	}

	data.From(sc.SaveServiceContext, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *ServiceContextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := new(serviceContext)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetServiceContext(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service context, got error: %s", err))
		return
	}

	data.From(response.ServiceContext, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *ServiceContextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(serviceContext)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.SaveServiceContext(ctx, data.Name.ValueString(), data.Attributes(ctx, resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update service context, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *ServiceContextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(serviceContext)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteServiceContext(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service context, got error: %s", err))
		return
	}

	err = wait.WaitForWithContext(ctx, client.Ticker(5*time.Second), func(ctx context.Context) (bool, error) {
		_, err := r.client.GetServiceContext(ctx, data.Name.ValueString())
		if client.IsNotFound(err) {
			return true, nil
		}

		return false, err
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during watiting for service context to be deleted, got error: %s", err))
		return
	}
}

func (r *ServiceContextResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
