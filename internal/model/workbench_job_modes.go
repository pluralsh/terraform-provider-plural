package model

import (
	"context"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

type WorkbenchJobModes struct {
	Plan         types.Bool                   `tfsdk:"plan"`
	Verification types.Bool                   `tfsdk:"verification"`
	Model        *WorkbenchJobModel           `tfsdk:"model"`
	Coding       *WorkbenchJobCodingModes     `tfsdk:"coding"`
	Budget       *WorkbenchJobBudget          `tfsdk:"budget"`
	Kubernetes   *WorkbenchJobKubernetesModes `tfsdk:"kubernetes"`
}

func (in *WorkbenchJobModes) Attributes(ctx context.Context, d *diag.Diagnostics) *gqlclient.WorkbenchJobModesAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchJobModesAttributes{
		Plan:         in.Plan.ValueBoolPointer(),
		Verification: in.Verification.ValueBoolPointer(),
		Model:        in.Model.Attributes(),
		Coding:       in.Coding.Attributes(),
		Budget:       in.Budget.Attributes(),
		Kubernetes:   in.Kubernetes.Attributes(ctx, d),
	}
}

func (in *WorkbenchJobModes) From(response *gqlclient.WorkbenchJobModesFragment, ctx context.Context, d *diag.Diagnostics) {
	in.Plan = types.BoolPointerValue(response.Plan)
	in.Verification = types.BoolPointerValue(response.Verification)

	if response.Model != nil {
		ensure(&in.Model)
		in.Model.From(response.Model)
	} else {
		in.Model = nil
	}

	if response.Coding != nil {
		in.Coding = &WorkbenchJobCodingModes{
			Babysit:  types.BoolPointerValue(response.Coding.Babysit),
			Approval: types.BoolPointerValue(response.Coding.Approval),
			Review:   types.BoolPointerValue(response.Coding.Review),
		}
	} else {
		in.Coding = nil
	}

	if response.Budget != nil {
		in.Budget = &WorkbenchJobBudget{
			Cost:   types.Float64PointerValue(response.Budget.Cost),
			Tokens: types.Int64PointerValue(response.Budget.Tokens),
		}
	} else {
		in.Budget = nil
	}

	if response.Kubernetes != nil {
		ensure(&in.Kubernetes)
		in.Kubernetes.From(response.Kubernetes, ctx, d)
	} else {
		in.Kubernetes = nil
	}
}

type WorkbenchJobModel struct {
	Provider types.String `tfsdk:"provider"`
	Model    types.String `tfsdk:"model"`
}

func (in *WorkbenchJobModel) Attributes() *gqlclient.WorkbenchJobModelAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchJobModelAttributes{
		Provider: gqlclient.AiProvider(in.Provider.ValueString()),
		Model:    in.Model.ValueString(),
	}
}

func (in *WorkbenchJobModel) From(response *gqlclient.WorkbenchJobModesFragment_Model) {
	if response.Provider != nil {
		in.Provider = types.StringValue(string(*response.Provider))
	} else {
		in.Provider = types.StringNull()
	}
	in.Model = types.StringPointerValue(response.Model)
}

type WorkbenchJobCodingModes struct {
	Babysit  types.Bool `tfsdk:"babysit"`
	Approval types.Bool `tfsdk:"approval"`
	Review   types.Bool `tfsdk:"review"`
}

func (in *WorkbenchJobCodingModes) Attributes() *gqlclient.WorkbenchJobCodingModesAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchJobCodingModesAttributes{
		Babysit:  in.Babysit.ValueBoolPointer(),
		Approval: in.Approval.ValueBoolPointer(),
		Review:   in.Review.ValueBoolPointer(),
	}
}

type WorkbenchJobBudget struct {
	Cost   types.Float64 `tfsdk:"cost"`
	Tokens types.Int64   `tfsdk:"tokens"`
}

func (in *WorkbenchJobBudget) Attributes() *gqlclient.WorkbenchJobBudgetAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchJobBudgetAttributes{
		Cost:   in.Cost.ValueFloat64Pointer(),
		Tokens: in.Tokens.ValueInt64Pointer(),
	}
}

type WorkbenchJobKubernetesModes struct {
	Update            types.Bool `tfsdk:"update"`
	Delete            types.Bool `tfsdk:"delete"`
	Exec              types.Bool `tfsdk:"exec"`
	Drain             types.Bool `tfsdk:"drain"`
	ExcludeNamespaces types.Set  `tfsdk:"exclude_namespaces"`
	RequireNamespaces types.Set  `tfsdk:"require_namespaces"`
}

func (in *WorkbenchJobKubernetesModes) Attributes(ctx context.Context, d *diag.Diagnostics) *gqlclient.WorkbenchJobKubernetesModesAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchJobKubernetesModesAttributes{
		Update:            in.Update.ValueBoolPointer(),
		Delete:            in.Delete.ValueBoolPointer(),
		Exec:              in.Exec.ValueBoolPointer(),
		Drain:             in.Drain.ValueBoolPointer(),
		ExcludeNamespaces: stringPointersFrom(in.ExcludeNamespaces, ctx, d),
		RequireNamespaces: stringPointersFrom(in.RequireNamespaces, ctx, d),
	}
}

func (in *WorkbenchJobKubernetesModes) From(response *gqlclient.WorkbenchJobModesFragment_Kubernetes, ctx context.Context, d *diag.Diagnostics) {
	in.Update = types.BoolPointerValue(response.Update)
	in.Delete = types.BoolPointerValue(response.Delete)
	in.Exec = types.BoolPointerValue(response.Exec)
	in.Drain = types.BoolPointerValue(response.Drain)
	in.ExcludeNamespaces = common.SetFrom(response.ExcludeNamespaces, in.ExcludeNamespaces, ctx, d)
	in.RequireNamespaces = common.SetFrom(response.RequireNamespaces, in.RequireNamespaces, ctx, d)
}

type stringCollection interface {
	IsNull() bool
	IsUnknown() bool
	Elements() []attr.Value
	ElementsAs(ctx context.Context, target interface{}, allowUnhandled bool) diag.Diagnostics
}

// stringPointersFrom converts a set or list of strings to a slice of string pointers,
// returning nil for a null or unknown collection.
func stringPointersFrom(collection stringCollection, ctx context.Context, d *diag.Diagnostics) []*string {
	if collection.IsNull() || collection.IsUnknown() {
		return nil
	}

	values := make([]types.String, len(collection.Elements()))
	d.Append(collection.ElementsAs(ctx, &values, false)...)
	return lo.Map(values, func(v types.String, _ int) *string { return v.ValueStringPointer() })
}
