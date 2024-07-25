package model

import (
	"context"
	"strings"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

type GlobalService struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	ServiceId  types.String `tfsdk:"service_id"`
	Distro     types.String `tfsdk:"distro"`
	ProviderId types.String `tfsdk:"provider_id"`
	Tags       types.Map    `tfsdk:"tags"`
}

func (gs *GlobalService) From(response *gqlclient.GlobalServiceFragment, d diag.Diagnostics) {
	gs.Id = types.StringValue(response.ID)
	gs.Name = types.StringValue(response.Name)
	gs.ServiceId = types.StringValue(response.Service.ID)
	if response.Distro != nil {
		gs.Distro = types.StringValue(string(*response.Distro))
	}
	if response.Provider != nil {
		gs.ProviderId = types.StringValue(response.Provider.ID)
	}
	gs.Tags = common.TagsFrom(response.Tags, gs.Tags, d)

}

func (gs *GlobalService) Attributes(ctx context.Context, d diag.Diagnostics) gqlclient.GlobalServiceAttributes {
	var distro *gqlclient.ClusterDistro
	if !gs.Distro.IsNull() {
		distro = lo.ToPtr(gqlclient.ClusterDistro(strings.ToUpper(gs.Distro.ValueString())))
	}
	return gqlclient.GlobalServiceAttributes{
		Name:       gs.Name.ValueString(),
		Distro:     distro,
		ProviderID: gs.ProviderId.ValueStringPointer(),
		Tags:       gs.TagsAttribute(ctx, d),
	}
}

func (gs *GlobalService) TagsAttribute(ctx context.Context, d diag.Diagnostics) []*gqlclient.TagAttributes {
	if gs.Tags.IsNull() {
		return nil
	}

	result := make([]*gqlclient.TagAttributes, 0)
	elements := make(map[string]types.String, len(gs.Tags.Elements()))
	d.Append(gs.Tags.ElementsAs(ctx, &elements, false)...)

	for k, v := range elements {
		result = append(result, &gqlclient.TagAttributes{Name: k, Value: v.ValueString()})
	}

	return result
}
