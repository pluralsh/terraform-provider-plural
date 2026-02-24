package resource

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type uppercaseStringPlanModifier struct{}

func (m uppercaseStringPlanModifier) Description(_ context.Context) string {
	return "Normalizes string values to uppercase."
}

func (m uppercaseStringPlanModifier) MarkdownDescription(_ context.Context) string {
	return "Normalizes string values to uppercase."
}

func (m uppercaseStringPlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	resp.PlanValue = types.StringValue(strings.ToUpper(req.PlanValue.ValueString()))
}

func UppercaseString() planmodifier.String {
	return uppercaseStringPlanModifier{}
}
