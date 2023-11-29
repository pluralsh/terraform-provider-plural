package cluster

import (
	"context"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console-client-go"
)

type cluster struct {
	Id             types.String              `tfsdk:"id"`
	InsertedAt     types.String              `tfsdk:"inserted_at"`
	Name           types.String              `tfsdk:"name"`
	Handle         types.String              `tfsdk:"handle"`
	Version        types.String              `tfsdk:"version"`
	CurrentVersion types.String              `tfsdk:"current_version"`
	DesiredVersion types.String              `tfsdk:"desired_version"`
	ProviderId     types.String              `tfsdk:"provider_id"`
	Cloud          types.String              `tfsdk:"cloud"`
	Protect        types.Bool                `tfsdk:"protect"`
	Tags           types.Map                 `tfsdk:"tags"`
	Bindings       *common.ClusterBindings   `tfsdk:"bindings"`
	NodePools      []*common.ClusterNodePool `tfsdk:"node_pools"`
	CloudSettings  *ClusterCloudSettings     `tfsdk:"cloud_settings"`
}

func (c *cluster) NodePoolsAttribute(ctx context.Context, d diag.Diagnostics) []*console.NodePoolAttributes {
	result := make([]*console.NodePoolAttributes, 0, len(c.NodePools))
	for _, np := range c.NodePools {
		result = append(result, &console.NodePoolAttributes{
			Name:         np.Name.ValueString(),
			MinSize:      np.MinSize.ValueInt64(),
			MaxSize:      np.MaxSize.ValueInt64(),
			InstanceType: np.InstanceType.ValueString(),
			Labels:       np.LabelsAttribute(ctx, d),
			Taints:       np.TaintsAttribute(),
			//CloudSettings: np.CloudSettings.Attributes(), TODO
		})
	}

	return result
}

func (c *cluster) TagsAttribute(ctx context.Context, d diag.Diagnostics) (result []*console.TagAttributes) {
	elements := make(map[string]types.String, len(c.Tags.Elements()))
	d.Append(c.Tags.ElementsAs(ctx, &elements, false)...)
	for k, v := range elements {
		result = append(result, &console.TagAttributes{Name: k, Value: v.ValueString()})
	}

	return
}

func (c *cluster) Attributes(ctx context.Context, d diag.Diagnostics) console.ClusterAttributes {
	return console.ClusterAttributes{
		Name:          c.Name.ValueString(),
		Handle:        c.Handle.ValueStringPointer(),
		ProviderID:    c.ProviderId.ValueStringPointer(),
		Version:       c.Version.ValueStringPointer(),
		Protect:       c.Protect.ValueBoolPointer(),
		CloudSettings: c.CloudSettings.Attributes(),
		ReadBindings:  c.Bindings.ReadAttributes(),
		WriteBindings: c.Bindings.WriteAttributes(),
		Tags:          c.TagsAttribute(ctx, d),
		NodePools:     c.NodePoolsAttribute(ctx, d),
	}
}

func (c *cluster) UpdateAttributes(ctx context.Context, d diag.Diagnostics) console.ClusterUpdateAttributes {
	return console.ClusterUpdateAttributes{
		Version:   c.Version.ValueStringPointer(),
		Handle:    c.Handle.ValueStringPointer(),
		Protect:   c.Protect.ValueBoolPointer(),
		NodePools: c.NodePoolsAttribute(ctx, d),
	}
}

func (c *cluster) From(cl *console.ClusterFragment, d diag.Diagnostics) {
	c.Id = types.StringValue(cl.ID)
	c.InsertedAt = types.StringPointerValue(cl.InsertedAt)
	c.Name = types.StringValue(cl.Name)
	c.Handle = types.StringPointerValue(cl.Handle)
	c.DesiredVersion = types.StringPointerValue(cl.Version)
	c.CurrentVersion = types.StringPointerValue(cl.CurrentVersion)
	c.Protect = types.BoolPointerValue(cl.Protect)
	c.NodePools = common.ClusterNodePoolsFrom(cl.NodePools)
	c.Tags = common.ClusterTagsFrom(cl.Tags, d)
	c.ProviderId = common.ClusterProviderIdFrom(cl.Provider)
}

func (c *cluster) FromCreate(cc *console.CreateCluster, d diag.Diagnostics) {
	c.Id = types.StringValue(cc.CreateCluster.ID)
	c.InsertedAt = types.StringPointerValue(cc.CreateCluster.InsertedAt)
	c.Name = types.StringValue(cc.CreateCluster.Name)
	c.Handle = types.StringPointerValue(cc.CreateCluster.Handle)
	c.DesiredVersion = types.StringPointerValue(cc.CreateCluster.Version)
	c.CurrentVersion = types.StringPointerValue(cc.CreateCluster.CurrentVersion)
	c.Protect = types.BoolPointerValue(cc.CreateCluster.Protect)
	c.NodePools = common.ClusterNodePoolsFrom(cc.CreateCluster.NodePools)
	c.Tags = common.ClusterTagsFrom(cc.CreateCluster.Tags, d)
	c.ProviderId = common.ClusterProviderIdFrom(cc.CreateCluster.Provider)
}
