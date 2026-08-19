package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

type BindingPolicy struct {
	Id           types.String          `tfsdk:"id"`
	PolicyId     types.String          `tfsdk:"policy_id"`
	BindPolicyId types.String          `tfsdk:"bind_policy_id"`
	Type         types.String          `tfsdk:"type"`
	Interval     types.String          `tfsdk:"interval"`
	Matches      *BindingPolicyMatches `tfsdk:"matches"`
}

type BindingPolicyMatches struct {
	Workbench *WorkbenchPolicyMatches `tfsdk:"workbench"`
}

type WorkbenchPolicyMatches struct {
	Regexes types.List `tfsdk:"regexes"`
}

func (p *BindingPolicy) Attributes(ctx context.Context, d *diag.Diagnostics) gqlclient.BindingPolicyAttributes {
	return gqlclient.BindingPolicyAttributes{
		PolicyID:     p.PolicyId.ValueString(),
		BindPolicyID: p.BindPolicyId.ValueString(),
		Type:         gqlclient.BindingPolicyType(p.Type.ValueString()),
		Interval:     p.Interval.ValueStringPointer(),
		Matches:      p.matchesAttributes(ctx, d),
	}
}

func (p *BindingPolicy) UpdateAttributes(ctx context.Context, d *diag.Diagnostics) gqlclient.BindingPolicyUpdateAttributes {
	return gqlclient.BindingPolicyUpdateAttributes{
		PolicyID:     p.PolicyId.ValueStringPointer(),
		BindPolicyID: p.BindPolicyId.ValueStringPointer(),
		Type:         gqlclient.BindingPolicyType(p.Type.ValueString()),
		Interval:     p.Interval.ValueStringPointer(),
		Matches:      p.matchesAttributes(ctx, d),
	}
}

func (p *BindingPolicy) From(response *gqlclient.BindingPolicyFragment, ctx context.Context, d *diag.Diagnostics) {
	p.Id = types.StringValue(response.ID)
	p.Type = types.StringValue(string(response.Type))
	p.Interval = types.StringValue(response.Interval)
	p.PolicyId = policyIDFrom(response.Policy)
	p.BindPolicyId = policyIDFrom(response.BindPolicy)
	p.Matches = bindingPolicyMatchesFrom(response.Matches, ctx, d)
}

func (p *BindingPolicy) matchesAttributes(ctx context.Context, d *diag.Diagnostics) *gqlclient.BindingPolicyMatchesAttributes {
	if p.Matches == nil || p.Matches.Workbench == nil {
		return nil
	}

	regexes := make([]string, 0)
	d.Append(p.Matches.Workbench.Regexes.ElementsAs(ctx, &regexes, false)...)
	if d.HasError() {
		return nil
	}

	return &gqlclient.BindingPolicyMatchesAttributes{
		Workbench: &gqlclient.WorkbenchPolicyMatchesAttributes{Regexes: lo.ToSlicePtr(regexes)},
	}
}

func bindingPolicyMatchesFrom(response *gqlclient.BindingPolicyFragment_Matches, ctx context.Context, d *diag.Diagnostics) *BindingPolicyMatches {
	if response == nil || response.Workbench == nil {
		return nil
	}

	regexes, diagnostics := types.ListValueFrom(ctx, types.StringType, response.Workbench.Regexes)
	d.Append(diagnostics...)

	return &BindingPolicyMatches{
		Workbench: &WorkbenchPolicyMatches{Regexes: regexes},
	}
}

func policyIDFrom(policy *gqlclient.TinyPolicyFragment) types.String {
	if policy == nil {
		return types.StringNull()
	}

	return types.StringValue(policy.ID)
}
