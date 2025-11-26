package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ensureAgentPlanModifier struct{}

func (in ensureAgentPlanModifier) Description(_ context.Context) string {
	return "Forces resource update when agent is not deployed"
}

func (in ensureAgentPlanModifier) MarkdownDescription(_ context.Context) string {
	return "Forces resource update when agent is not deployed"
}

func (in ensureAgentPlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.State.Raw.IsNull() {
		return
	}

	var agentDeployed types.Bool
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("agent_deployed"), &agentDeployed)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the agent is not deployed, force update by setting the field to unknown.
	if !agentDeployed.IsNull() && !agentDeployed.ValueBool() {
		resp.PlanValue = types.BoolUnknown()
	}
}

func EnsureAgent() planmodifier.Bool {
	return ensureAgentPlanModifier{}
}
