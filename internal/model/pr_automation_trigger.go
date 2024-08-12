package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-plural/internal/common"
)

type PrAutomationTrigger struct {
	PrAutomationID     types.String `tfsdk:"pr_automation_id"`
	RepoSlug           types.String `tfsdk:"repo_slug"`
	PrAutomationBranch types.String `tfsdk:"pr_automation_branch"`
	Context            types.Map    `tfsdk:"context"`
}

func (in *PrAutomationTrigger) ContextJson(ctx context.Context, d diag.Diagnostics) *string {
	triggerContext := make(map[string]types.String, len(in.Context.Elements()))
	in.Context.ElementsAs(ctx, &triggerContext, false)

	return common.AttributesJson(triggerContext, d)
}
