package resource

import (
	"context"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console-client-go"
)

type project struct {
	Id          types.String     `tfsdk:"id"`
	Name        types.String     `tfsdk:"name"`
	Description types.String     `tfsdk:"description"`
	Default     types.Bool       `tfsdk:"default"`
	Bindings    *common.Bindings `tfsdk:"bindings"`
}

func (p *project) Attributes(ctx context.Context, d diag.Diagnostics) (*gqlclient.ProjectAttributes, error) {
	return &gqlclient.ProjectAttributes{
		Name:          p.Name.ValueString(),
		Description:   p.Description.ValueStringPointer(),
		ReadBindings:  p.Bindings.ReadAttributes(ctx, d),
		WriteBindings: p.Bindings.WriteAttributes(ctx, d),
	}, nil
}

func (p *project) From(project *gqlclient.ProjectFragment, ctx context.Context, d diag.Diagnostics) {
	p.Id = types.StringValue(project.ID)
	p.Name = types.StringValue(project.Name)
	p.Description = types.StringPointerValue(project.Description)
	p.Default = types.BoolPointerValue(project.Default)
	p.Bindings.From(project.ReadBindings, project.WriteBindings, ctx, d)
}
