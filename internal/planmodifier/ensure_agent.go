package resource

import (
	"context"

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

func (in ensureAgentPlanModifier) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.State.Raw.IsNull() {
		return
	}

	if req.StateValue.IsUnknown() || req.StateValue.IsNull() {
		return
	}

	// If the agent is not deployed, force update by setting the field to unknown.
	if !req.StateValue.ValueBool() {
		resp.PlanValue = types.BoolUnknown()
		return
	}

	// Terraform plans computed values as unknown on update. Keep the existing
	// true value when the agent is already deployed.
	if req.PlanValue.IsUnknown() {
		resp.PlanValue = req.StateValue
	}
}

func EnsureAgent() planmodifier.Bool {
	return ensureAgentPlanModifier{}
}
