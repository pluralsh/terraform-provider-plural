package resource

import (
	"context"
	"strings"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console-client-go"
	"github.com/samber/lo"
)

type globalService struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	ServiceId  types.String `tfsdk:"service_id"`
	Distro     types.String `tfsdk:"distro"`
	ProviderId types.String `tfsdk:"provider_id"`
	Tags       types.Map    `tfsdk:"tags"`
}

func (g *globalService) From(response *gqlclient.GlobalServiceFragment, d diag.Diagnostics) {
	g.Id = types.StringValue(response.ID)
	g.Name = types.StringValue(response.Name)
	g.ServiceId = types.StringValue(response.Service.ID)
	if response.Distro != nil {
		g.Distro = types.StringValue(string(*response.Distro))
	}
	if response.Provider != nil {
		g.ProviderId = types.StringValue(response.Provider.ID)
	}
	g.Tags = common.TagsFrom(response.Tags, g.Tags, d)

}

func (g *globalService) Attributes(ctx context.Context, d diag.Diagnostics) gqlclient.GlobalServiceAttributes {
	var distro *gqlclient.ClusterDistro
	if !g.Distro.IsNull() {
		distro = lo.ToPtr(gqlclient.ClusterDistro(strings.ToUpper(g.Distro.ValueString())))
	}
	return gqlclient.GlobalServiceAttributes{
		Name:       g.Name.ValueString(),
		Distro:     distro,
		ProviderID: g.ProviderId.ValueStringPointer(),
		Tags:       g.TagsAttribute(ctx, d),
	}
}

func (g *globalService) TagsAttribute(ctx context.Context, d diag.Diagnostics) []*gqlclient.TagAttributes {
	if g.Tags.IsNull() {
		return nil
	}

	result := make([]*gqlclient.TagAttributes, 0)
	elements := make(map[string]types.String, len(g.Tags.Elements()))
	d.Append(g.Tags.ElementsAs(ctx, &elements, false)...)

	for k, v := range elements {
		result = append(result, &gqlclient.TagAttributes{Name: k, Value: v.ValueString()})
	}

	return result
}
