package datasource

import (
	"context"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console-client-go"
)

type cluster struct {
	Id             types.String `tfsdk:"id"`
	InsertedAt     types.String `tfsdk:"inserted_at"`
	Name           types.String `tfsdk:"name"`
	Handle         types.String `tfsdk:"handle"`
	CurrentVersion types.String `tfsdk:"current_version"`
	DesiredVersion types.String `tfsdk:"desired_version"`
	ProviderId     types.String `tfsdk:"provider_id"`
	Cloud          types.String `tfsdk:"cloud"`
	Protect        types.Bool   `tfsdk:"protect"`
	Tags           types.Map    `tfsdk:"tags"`
	NodePools      types.Map    `tfsdk:"node_pools"`
}

func (c *cluster) From(cl *console.ClusterFragment, ctx context.Context, d diag.Diagnostics) {
	c.Id = types.StringValue(cl.ID)
	c.InsertedAt = types.StringPointerValue(cl.InsertedAt)
	c.Name = types.StringValue(cl.Name)
	c.Handle = types.StringPointerValue(cl.Handle)
	c.DesiredVersion = types.StringPointerValue(cl.Version)
	c.CurrentVersion = types.StringPointerValue(cl.CurrentVersion)
	c.Protect = types.BoolPointerValue(cl.Protect)
	c.fromNodePools(cl.NodePools, ctx, d)
	c.Tags = common.ClusterTagsFrom(cl.Tags, d)
	c.ProviderId = common.ClusterProviderIdFrom(cl.Provider)
}

func (c *cluster) fromNodePools(nodePools []*console.NodePoolFragment, ctx context.Context, d diag.Diagnostics) {
	mapValue, diags := types.MapValue(basetypes.ObjectType{AttrTypes: common.ClusterNodePoolAttrTypes},
		common.ClusterNodePoolsFrom(nodePools, ctx, d))
	d.Append(diags...)
	c.NodePools = mapValue
}
