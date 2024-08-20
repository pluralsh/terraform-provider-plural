package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/model"
)

var _ resource.ResourceWithConfigure = &stackRunTriggerResource{}

func NewStackRunTriggerResource() resource.Resource {
	return &stackRunTriggerResource{}
}

type stackRunTriggerResource struct {
	client *client.Client
}

func (in *stackRunTriggerResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_stack_run_trigger"
}

func (in *stackRunTriggerResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the Infrastructure Stack that should be used to start a new run from the newest SHA in the stack's run history.",
				MarkdownDescription: "ID of the Infrastructure Stack that should be used to start a new run from the newest SHA in the stack's run history.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"retrigger_key": schema.StringAttribute{
				Description:         "Every time this key changes stack run will be retriggered.",
				MarkdownDescription: "Every time this key changes stack run will be retriggered.",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (in *stackRunTriggerResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	data, ok := request.ProviderData.(*common.ProviderData)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Project Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
		return
	}

	in.client = data.Client
}

func (in *stackRunTriggerResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	data := new(model.StackRunTrigger)
	response.Diagnostics.Append(request.Plan.Get(ctx, data)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := in.client.TriggerRun(
		ctx,
		data.ID.ValueString(),
	)
	if err != nil {
		response.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to trigger stack run, got error: %s", err))
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (in *stackRunTriggerResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// Since this is only a trigger, there is no read API. Ignore.
}

func (in *stackRunTriggerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state model.StackRunTrigger
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.RetriggerKey.Equal(state.RetriggerKey) {
		_, err := in.client.TriggerRun(
			ctx,
			data.ID.ValueString(),
		)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to trigger stack run, got error: %s", err))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (in *stackRunTriggerResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Since this is only a trigger, there is no delete API. Ignore.
}
