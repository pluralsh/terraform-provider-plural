package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
)

type Bindings struct {
	Read  types.Set `tfsdk:"read"`
	Write types.Set `tfsdk:"write"`
}

func (cb *Bindings) ReadAttributes(ctx context.Context, d *diag.Diagnostics) []*console.PolicyBindingAttributes {
	if cb == nil {
		return nil
	}

	return policyBindingAttributes(cb.Read, ctx, d)
}

func (cb *Bindings) WriteAttributes(ctx context.Context, d *diag.Diagnostics) []*console.PolicyBindingAttributes {
	if cb == nil {
		return nil
	}

	return policyBindingAttributes(cb.Write, ctx, d)
}

func policyBindingAttributes(bindings types.Set, ctx context.Context, d *diag.Diagnostics) []*console.PolicyBindingAttributes {
	if bindings.IsNull() {
		return nil
	}

	result := make([]*console.PolicyBindingAttributes, 0, len(bindings.Elements()))
	elements := make([]PolicyBinding, len(bindings.Elements()))
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

func (cb *Bindings) From(readBindings []*console.PolicyBindingFragment, writeBindings []*console.PolicyBindingFragment, ctx context.Context, d *diag.Diagnostics) {
	if cb == nil {
		return
	}

	cb.Read = bindingsFrom(readBindings, cb.Read, ctx, d)
	cb.Write = bindingsFrom(writeBindings, cb.Write, ctx, d)
}

func bindingsFrom(bindings []*console.PolicyBindingFragment, config types.Set, ctx context.Context, d *diag.Diagnostics) types.Set {
	if len(bindings) == 0 {
		// Rewriting config to state to avoid inconsistent result errors.
		// This could happen, for example, when sending "nil" to API and "[]" is returned as a result.
		return config
	}

	values := make([]attr.Value, len(bindings))
	for i, binding := range bindings {
		value := PolicyBinding{
			ID: types.StringPointerValue(binding.ID),
		}

		if binding.User != nil {
			value.UserID = types.StringValue(binding.User.ID)
		}

		if binding.Group != nil {
			value.GroupID = types.StringValue(binding.Group.ID)
		}

		objValue, diags := types.ObjectValueFrom(ctx, PolicyBindingAttrTypes, value)
		values[i] = objValue
		d.Append(diags...)
	}

	setValue, diags := types.SetValue(types.ObjectType{AttrTypes: PolicyBindingAttrTypes}, values)
	d.Append(diags...)
	return setValue
}

func BindingsFromReadOnly(bindings []*console.PolicyBindingFragment, planned types.Set, ctx context.Context, d *diag.Diagnostics) types.Set {
	if len(bindings) == 0 {
		// Rewriting planned to state to avoid inconsistent result errors.
		// This could happen, for example, when sending "nil" to API and "[]" is returned as a result.
		return planned
	}

	values := make([]attr.Value, len(bindings))
	for i, binding := range bindings {
		value := PolicyBinding{
			ID: types.StringPointerValue(binding.ID),
		}

		if binding.User != nil {
			value.UserID = types.StringValue(binding.User.ID)
		}

		if binding.Group != nil {
			value.GroupID = types.StringValue(binding.Group.ID)
		}

		objValue, diags := types.ObjectValueFrom(ctx, PolicyBindingAttrTypes, value)
		values[i] = objValue
		d.Append(diags...)
	}

	setValue, diags := types.SetValue(types.ObjectType{AttrTypes: PolicyBindingAttrTypes}, values)
	d.Append(diags...)
	return setValue
}

func SetToPolicyBindingAttributes(set types.Set, ctx context.Context, d *diag.Diagnostics) []*console.PolicyBindingAttributes {
	return policyBindingAttributes(set, ctx, d)
}

type PolicyBinding struct {
	GroupID types.String `tfsdk:"group_id"`
	ID      types.String `tfsdk:"id"`
	UserID  types.String `tfsdk:"user_id"`
}

var PolicyBindingAttrTypes = map[string]attr.Type{
	"group_id": types.StringType,
	"id":       types.StringType,
	"user_id":  types.StringType,
}

func (cpb *PolicyBinding) Attributes() *console.PolicyBindingAttributes {
	return &console.PolicyBindingAttributes{
		ID:      cpb.ID.ValueStringPointer(),
		UserID:  cpb.UserID.ValueStringPointer(),
		GroupID: cpb.GroupID.ValueStringPointer(),
	}
}
