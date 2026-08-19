package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
)

type Policy struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	Policy      types.String `tfsdk:"policy"`
	ProjectId   types.String `tfsdk:"project_id"`
}

func (p *Policy) Attributes() gqlclient.PolicyAttributes {
	return gqlclient.PolicyAttributes{
		Name:        p.Name.ValueStringPointer(),
		Type:        policyTypePointer(p.Type),
		Description: p.Description.ValueStringPointer(),
		Policy:      p.Policy.ValueStringPointer(),
		ProjectID:   p.ProjectId.ValueStringPointer(),
	}
}

func (p *Policy) From(response *gqlclient.PolicyFragment) {
	p.Id = types.StringValue(response.ID)
	p.Name = types.StringValue(response.Name)
	p.Type = types.StringValue(string(response.Type))
	p.Description = types.StringPointerValue(response.Description)
	p.Policy = types.StringValue(response.Policy)

	if response.Project == nil {
		p.ProjectId = types.StringNull()
		return
	}

	p.ProjectId = types.StringValue(response.Project.ID)
}

func policyTypePointer(value types.String) *gqlclient.PolicyType {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	policyType := gqlclient.PolicyType(value.ValueString())
	return &policyType
}
