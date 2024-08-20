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
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/model"
)

var _ resource.ResourceWithConfigure = &prAutomationTriggerResource{}

func NewPrAutomationTriggerResource() resource.Resource {
	return &prAutomationTriggerResource{}
}

type prAutomationTriggerResource struct {
	client *client.Client
}

func (in *prAutomationTriggerResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_pr_automation_trigger"
}

func (in *prAutomationTriggerResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"pr_automation_id": schema.StringAttribute{
				Description:         "ID of the PR Automation that should be triggered.",
				MarkdownDescription: "ID of the PR Automation that should be triggered.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"repo_slug": schema.StringAttribute{
				Description:         "Repo slug of the repository PR Automation should be triggered against. If not provided PR Automation repo will be used. Example format for a github repository: <userOrOrg>/<repoName>",
				MarkdownDescription: "Repo slug of the repository PR Automation should be triggered against. If not provided PR Automation repo will be used. Example format for a github repository: <userOrOrg>/<repoName>",
				Optional:            true,
			},
			"pr_automation_branch": schema.StringAttribute{
				Description:         "Branch that should be created against PR Automation base branch.",
				MarkdownDescription: "Branch that should be created against PR Automation base branch.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"context": schema.MapAttribute{
				Description:         "PR Automation configuration context.",
				MarkdownDescription: "PR Automation configuration context.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"retrigger_key": schema.StringAttribute{
				Description:         "Every time this key changes PR automation will be retriggered.",
				MarkdownDescription: "Every time this key changes PR automation will be retriggered.",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (in *prAutomationTriggerResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (in *prAutomationTriggerResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	data := new(model.PrAutomationTrigger)
	response.Diagnostics.Append(request.Plan.Get(ctx, data)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := in.client.CreatePullRequest(
		ctx,
		data.PrAutomationID.ValueString(),
		data.RepoSlug.ValueStringPointer(),
		data.PrAutomationBranch.ValueStringPointer(),
		data.ContextJson(ctx, response.Diagnostics),
	)
	if err != nil {
		response.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create pull request, got error: %s", err))
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (in *prAutomationTriggerResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// Since this is only a trigger, there is no read API. Ignore.
}

func (in *prAutomationTriggerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state model.PrAutomationTrigger
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.RetriggerKey.Equal(state.RetriggerKey) {
		_, err := in.client.CreatePullRequest(
			ctx,
			data.PrAutomationID.ValueString(),
			data.RepoSlug.ValueStringPointer(),
			data.PrAutomationBranch.ValueStringPointer(),
			data.ContextJson(ctx, resp.Diagnostics),
		)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create pull request, got error: %s", err))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (in *prAutomationTriggerResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Since this is only a trigger, there is no delete API. Ignore.
}
