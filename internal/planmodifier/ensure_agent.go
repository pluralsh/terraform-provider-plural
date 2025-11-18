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

func (in ensureAgentPlanModifier) PlanModifyInt32(ctx context.Context, req planmodifier.Int32Request, resp *planmodifier.Int32Response) {
	if req.State.Raw.IsNull() {
		return
	}

	var agentDeployed types.Bool
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("agent_deployed"), &agentDeployed)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If agent is not deployed, mark reapply key as unknown to force update
	if !agentDeployed.IsNull() && !agentDeployed.ValueBool() {
		resp.PlanValue = types.Int32Unknown()
	}
}

func EnsureAgent() planmodifier.Int32 {
	return ensureAgentPlanModifier{}
}
