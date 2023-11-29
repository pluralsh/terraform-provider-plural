package cluster

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console-client-go"
)

type ClusterBindings struct {
	Read  []*ClusterPolicyBinding `tfsdk:"read"`
	Write []*ClusterPolicyBinding `tfsdk:"write"`
}

func (cb *ClusterBindings) ReadAttributes() []*console.PolicyBindingAttributes {
	if cb == nil {
		return []*console.PolicyBindingAttributes{}
	}

	return clusterPolicyBindingAttributes(cb.Read)
}

func (cb *ClusterBindings) WriteAttributes() []*console.PolicyBindingAttributes {
	if cb == nil {
		return []*console.PolicyBindingAttributes{}
	}

	return clusterPolicyBindingAttributes(cb.Write)
}

type ClusterPolicyBinding struct {
	GroupID types.String `tfsdk:"group_id"`
	ID      types.String `tfsdk:"id"`
	UserID  types.String `tfsdk:"user_id"`
}

func (c *ClusterPolicyBinding) Attributes() *console.PolicyBindingAttributes {
	return &console.PolicyBindingAttributes{
		ID:      c.ID.ValueStringPointer(),
		UserID:  c.UserID.ValueStringPointer(),
		GroupID: c.GroupID.ValueStringPointer(),
	}
}

func clusterPolicyBindingAttributes(bindings []*ClusterPolicyBinding) []*console.PolicyBindingAttributes {
	result := make([]*console.PolicyBindingAttributes, len(bindings))
	for i, b := range bindings {
		result[i] = b.Attributes()
	}

	return result
}
