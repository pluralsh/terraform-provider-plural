package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pluralsh/polly/algorithms"

	"terraform-provider-plural/internal/common"

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
	//Bindings       *common.ClusterBindings   `tfsdk:"bindings"`
	NodePools types.List `tfsdk:"node_pools"`
}

func (c *cluster) From(cl *console.ClusterFragment, d diag.Diagnostics) {
	c.Id = types.StringValue(cl.ID)
	c.InsertedAt = types.StringPointerValue(cl.InsertedAt)
	c.Name = types.StringValue(cl.Name)
	c.Handle = types.StringPointerValue(cl.Handle)
	c.DesiredVersion = types.StringPointerValue(cl.Version)
	c.CurrentVersion = types.StringPointerValue(cl.CurrentVersion)
	c.Protect = types.BoolPointerValue(cl.Protect)
	c.fromNodePools(cl.NodePools)
	c.Tags = common.ClusterTagsFrom(cl.Tags, d)
	c.ProviderId = common.ClusterProviderIdFrom(cl.Provider)
}

func (c *cluster) fromNodePools(nodePools []*console.NodePoolFragment) {
	commonNodePools := algorithms.Map(common.ClusterNodePoolsFrom(nodePools), func(nodePool *common.ClusterNodePool) attr.Value {
		return nodePool.Element()
	})

	c.NodePools = types.ListValueMust(basetypes.ObjectType{AttrTypes: map[string]attr.Type{
		"name":          types.StringType,
		"min_size":      types.Int64Type,
		"max_size":      types.Int64Type,
		"instance_type": types.StringType,
		"labels":        types.MapType{ElemType: types.StringType},
		"taints": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"key":    types.StringType,
			"value":  types.StringType,
			"effect": types.StringType,
		}}},
		"cloud_settings": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"aws": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"launch_template_id": types.StringType,
					},
				},
			},
		},
	}}, commonNodePools)
}
