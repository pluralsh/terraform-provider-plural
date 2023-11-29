package datasource

import (
	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console-client-go"
)

type cluster struct {
	Id         types.String            `tfsdk:"id"`
	InsertedAt types.String            `tfsdk:"inserted_at"`
	Name       types.String            `tfsdk:"name"`
	Handle     types.String            `tfsdk:"handle"`
	Version    types.String            `tfsdk:"version"`
	ProviderId types.String            `tfsdk:"provider_id"`
	Cloud      types.String            `tfsdk:"cloud"`
	Protect    types.Bool              `tfsdk:"protect"`
	Tags       types.Map               `tfsdk:"tags"`
	Bindings   *common.ClusterBindings `tfsdk:"bindings"`
	//NodePools     []*ClusterNodePool    `tfsdk:"node_pools"`
}

func (c *cluster) From(cl *console.ClusterFragment, d diag.Diagnostics) {
	c.Id = types.StringValue(cl.ID)
	c.InsertedAt = types.StringPointerValue(cl.InsertedAt)
	c.Name = types.StringValue(cl.Name)
	c.Handle = types.StringPointerValue(cl.Handle)
	c.Version = types.StringPointerValue(cl.Version)
	c.Protect = types.BoolPointerValue(cl.Protect)

	if cl.Provider != nil {
		c.ProviderId = types.StringValue(cl.Provider.ID)
	}

	tagsValue, tagsDiagnostics := types.MapValue(types.StringType, common.ClusterTagsMap(cl.Tags))
	c.Tags = tagsValue
	d.Append(tagsDiagnostics...)

	//c.NodePoolsFrom(cl.NodePools, d)
}
