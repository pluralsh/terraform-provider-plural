package resource

import (
	"context"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console-client-go"
)

type project struct {
	Id          types.String     `tfsdk:"id"`
	Name        types.String     `tfsdk:"name"`
	Default     types.Bool       `tfsdk:"default"`
	Description types.String     `tfsdk:"description"`
	Bindings    *common.Bindings `tfsdk:"bindings"`
}

func (p *project) Attributes(ctx context.Context, d diag.Diagnostics, client *client.Client) (*gqlclient.ProjectAttributes, error) {
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
	p.Default = types.BoolPointerValue(project.Default)
	p.Description = types.StringPointerValue(project.Description)
	p.Bindings.From(project.ReadBindings, project.WriteBindings, ctx, d)
}
