package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	console "github.com/pluralsh/console-client-go"
)

type ClusterBindings struct {
	Read  types.Set `tfsdk:"read"`
	Write types.Set `tfsdk:"write"`
}

func (cb *ClusterBindings) ReadAttributes(ctx context.Context, d diag.Diagnostics) []*console.PolicyBindingAttributes {
	if cb == nil {
		return nil
	}

	return clusterPolicyBindingAttributes(cb.Read, ctx, d)
}

func (cb *ClusterBindings) WriteAttributes(ctx context.Context, d diag.Diagnostics) []*console.PolicyBindingAttributes {
	if cb == nil {
		return nil
	}

	return clusterPolicyBindingAttributes(cb.Write, ctx, d)
}

func clusterPolicyBindingAttributes(bindings types.Set, ctx context.Context, d diag.Diagnostics) []*console.PolicyBindingAttributes {
	if bindings.IsNull() {
		return nil
	}

	result := make([]*console.PolicyBindingAttributes, 0, len(bindings.Elements()))
	elements := make([]ClusterPolicyBinding, len(bindings.Elements()))
	d.Append(bindings.ElementsAs(ctx, &elements, false)...)

	for _, binding := range elements {
		result = append(result, &console.PolicyBindingAttributes{
			ID:      binding.ID.ValueStringPointer(),
			UserID:  binding.UserID.ValueStringPointer(),
			GroupID: binding.GroupID.ValueStringPointer(),
		})
	}

	return result
}

func (cb *ClusterBindings) From(readBindings []*console.PolicyBindingFragment, writeBindings []*console.PolicyBindingFragment, ctx context.Context, d diag.Diagnostics) {
	if cb == nil {
		return
	}

	cb.Read = clusterBindingsFrom(readBindings, ctx, d)
	cb.Write = clusterBindingsFrom(writeBindings, ctx, d)
}

func clusterBindingsFrom(bindings []*console.PolicyBindingFragment, ctx context.Context, d diag.Diagnostics) types.Set {
	if bindings == nil {
		return types.SetNull(basetypes.ObjectType{AttrTypes: ClusterPolicyBindingAttrTypes})
	}

	values := make([]attr.Value, len(bindings))
	for i, binding := range bindings {
		value := ClusterPolicyBinding{
			ID: types.StringPointerValue(binding.ID),
		}

		if binding.User != nil {
			value.UserID = types.StringValue(binding.User.ID)
		}

		if binding.Group != nil {
			value.GroupID = types.StringValue(binding.Group.ID)
		}

		objValue, diags := types.ObjectValueFrom(ctx, ClusterPolicyBindingAttrTypes, value)
		values[i] = objValue
		d.Append(diags...)
	}

	setValue, diags := types.SetValue(basetypes.ObjectType{AttrTypes: ClusterPolicyBindingAttrTypes}, values)
	d.Append(diags...)
	return setValue
}

type ClusterPolicyBinding struct {
	GroupID types.String `tfsdk:"group_id"`
	ID      types.String `tfsdk:"id"`
	UserID  types.String `tfsdk:"user_id"`
}

var ClusterPolicyBindingAttrTypes = map[string]attr.Type{
	"group_id": types.StringType,
	"id":       types.StringType,
	"user_id":  types.StringType,
}

func (cpb *ClusterPolicyBinding) Attributes() *console.PolicyBindingAttributes {
	return &console.PolicyBindingAttributes{
		ID:      cpb.ID.ValueStringPointer(),
		UserID:  cpb.UserID.ValueStringPointer(),
		GroupID: cpb.GroupID.ValueStringPointer(),
	}
}
