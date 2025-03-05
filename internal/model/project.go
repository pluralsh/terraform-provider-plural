package model

import (
	"context"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
)

type Project struct {
	Id          types.String     `tfsdk:"id"`
	Name        types.String     `tfsdk:"name"`
	Description types.String     `tfsdk:"description"`
	Default     types.Bool       `tfsdk:"default"`
	Bindings    *common.Bindings `tfsdk:"bindings"`
}

func (p *Project) Attributes(ctx context.Context, d *diag.Diagnostics) (*gqlclient.ProjectAttributes, error) {
	return &gqlclient.ProjectAttributes{
		Name:          p.Name.ValueString(),
		Description:   p.Description.ValueStringPointer(),
		ReadBindings:  p.Bindings.ReadAttributes(ctx, d),
		WriteBindings: p.Bindings.WriteAttributes(ctx, d),
	}, nil
}

func (p *Project) From(response *gqlclient.ProjectFragment, ctx context.Context, d *diag.Diagnostics) {
	p.Id = types.StringValue(response.ID)
	p.Name = types.StringValue(response.Name)
	p.Description = types.StringPointerValue(response.Description)
	p.Default = defaultFrom(response.Default)
	p.Bindings.From(response.ReadBindings, response.WriteBindings, ctx, d)
}

func defaultFrom(def *bool) types.Bool {
	if def == nil {
		return types.BoolValue(false)
	}

	return types.BoolValue(*def)
}
