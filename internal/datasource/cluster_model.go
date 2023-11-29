package datasource

import (
	cluster2 "terraform-provider-plural/internal/resource/cluster"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console-client-go"
)

type cluster struct {
	Id         types.String              `tfsdk:"id"`
	InsertedAt types.String              `tfsdk:"inserted_at"`
	Name       types.String              `tfsdk:"name"`
	Handle     types.String              `tfsdk:"handle"`
	Version    types.String              `tfsdk:"version"`
	ProviderId types.String              `tfsdk:"provider_id"`
	Cloud      types.String              `tfsdk:"cloud"`
	Protect    types.Bool                `tfsdk:"protect"`
	Tags       types.Map                 `tfsdk:"tags"`
	Bindings   *cluster2.ClusterBindings `tfsdk:"bindings"`
	//NodePools     []*ClusterNodePool    `tfsdk:"node_pools"`
}

func (c *cluster) ProviderFrom(provider *console.ClusterProviderFragment) {
	if provider != nil {
		c.ProviderId = types.StringValue(provider.ID)
	}
}

func (c *cluster) TagsFrom(tags []*console.ClusterTags, d diag.Diagnostics) {
	elements := map[string]attr.Value{}
	for _, v := range tags {
		elements[v.Name] = types.StringValue(v.Value)
	}

	tagsValue, tagsDiagnostics := types.MapValue(types.StringType, elements)
	c.Tags = tagsValue
	d.Append(tagsDiagnostics...)
}

func (c *cluster) From(cl *console.ClusterFragment, d diag.Diagnostics) {
	c.Id = types.StringValue(cl.ID)
	c.InsertedAt = types.StringPointerValue(cl.InsertedAt)
	c.Name = types.StringValue(cl.Name)
	c.Handle = types.StringPointerValue(cl.Handle)
	c.Version = types.StringPointerValue(cl.Version)
	c.Protect = types.BoolPointerValue(cl.Protect)
	c.ProviderFrom(cl.Provider)
	// c.NodePoolsFrom(cl.NodePools, d)
	c.TagsFrom(cl.Tags, d)
}
