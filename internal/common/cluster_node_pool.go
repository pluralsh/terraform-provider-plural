package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	console "github.com/pluralsh/console-client-go"
)

type ClusterNodePool struct {
	Name          types.String `tfsdk:"name"`
	MinSize       types.Int64  `tfsdk:"min_size"`
	MaxSize       types.Int64  `tfsdk:"max_size"`
	InstanceType  types.String `tfsdk:"instance_type"`
	Labels        types.Map    `tfsdk:"labels"`
	Taints        types.Set    `tfsdk:"taints"`
	CloudSettings types.Object `tfsdk:"cloud_settings"`
}

var ClusterNodePoolAttrTypes = map[string]attr.Type{
	"name":          types.StringType,
	"min_size":      types.Int64Type,
	"max_size":      types.Int64Type,
	"instance_type": types.StringType,
	"labels":        types.MapType{ElemType: types.StringType},
	"taints":        types.SetType{ElemType: types.ObjectType{AttrTypes: NodePoolTaintAttrTypes}},
	//"cloud_settings": types.ObjectType{AttrTypes: NodePoolCloudSettingsAttrTypes},
}

func (c *ClusterNodePool) LabelsAttribute(ctx context.Context, d diag.Diagnostics) *string {
	if c.Labels.IsNull() {
		return nil
	}

	elements := make(map[string]types.String, len(c.Labels.Elements()))
	d.Append(c.Labels.ElementsAs(ctx, &elements, false)...)
	return AttributesJson(elements, d)
}

func (c *ClusterNodePool) TaintsAttribute(ctx context.Context, d diag.Diagnostics) []*console.TaintAttributes {
	if c.Taints.IsNull() {
		return nil
	}

	result := make([]*console.TaintAttributes, 0, len(c.Taints.Elements()))
	elements := make([]NodePoolTaint, len(c.Taints.Elements()))
	d.Append(c.Taints.ElementsAs(ctx, &elements, false)...)

	for _, np := range elements {
		result = append(result, &console.TaintAttributes{
			Key:    np.Key.ValueString(),
			Value:  np.Value.ValueString(),
			Effect: np.Effect.ValueString(),
		})
	}

	return result
}

func (c *ClusterNodePool) terraformAttributes() map[string]attr.Value {
	return map[string]attr.Value{
		"name":          c.Name,
		"min_size":      c.MinSize,
		"max_size":      c.MaxSize,
		"instance_type": c.InstanceType,
		"labels":        c.TerraformAttributesLabels(),
		"taints":        c.Taints,
		//"cloud_settings": c.CloudSettings,
	}
}

func (c *ClusterNodePool) TerraformAttributesLabels() attr.Value {
	if c.Labels.IsNull() {
		return types.MapNull(types.StringType)
	}

	return types.MapValueMust(types.StringType, c.Labels.Elements())
}

func (c *ClusterNodePool) Element() attr.Value {
	return types.ObjectValueMust(ClusterNodePoolAttrTypes, c.terraformAttributes())
}

type NodePoolTaint struct {
	Key    types.String `tfsdk:"key"`
	Value  types.String `tfsdk:"value"`
	Effect types.String `tfsdk:"effect"`
}

var NodePoolTaintAttrTypes = map[string]attr.Type{
	"key":    types.StringType,
	"value":  types.StringType,
	"effect": types.StringType,
}

type NodePoolCloudSettings struct {
	AWS *NodePoolCloudSettingsAWS `tfsdk:"aws"`
}

var NodePoolCloudSettingsAttrTypes = map[string]attr.Type{
	"aws": types.ObjectType{AttrTypes: NodePoolCloudSettingsAWSAttrTypes},
}

func (c *NodePoolCloudSettings) Attributes() *console.NodePoolCloudAttributes {
	if c == nil {
		return nil
	}

	if c.AWS != nil {
		return &console.NodePoolCloudAttributes{Aws: c.AWS.Attributes()}
	}

	return nil
}

type NodePoolCloudSettingsAWS struct {
	LaunchTemplateId types.String `tfsdk:"launch_template_id"`
}

var NodePoolCloudSettingsAWSAttrTypes = map[string]attr.Type{
	"launch_template_id": types.StringType,
}

func (c *NodePoolCloudSettingsAWS) Attributes() *console.AwsNodeCloudAttributes {
	return &console.AwsNodeCloudAttributes{
		LaunchTemplateID: c.LaunchTemplateId.ValueStringPointer(),
	}
}

func ClusterNodePoolsFrom(nodePools []*console.NodePoolFragment, ctx context.Context, d diag.Diagnostics) map[string]attr.Value {
	result := make(map[string]attr.Value)
	for _, nodePool := range nodePools {
		clusterNodePoolLabels, diags := types.MapValueFrom(ctx, types.StringType, nodePool.Labels)
		d.Append(diags...)

		clusterNodePoolCloudSettings, diags := types.ObjectValueFrom(ctx, NodePoolCloudSettingsAttrTypes,
			NodePoolCloudSettings{AWS: &NodePoolCloudSettingsAWS{LaunchTemplateId: types.StringNull()}})
		d.Append(diags...)

		result[nodePool.Name] = (&ClusterNodePool{
			Name:          types.StringValue(nodePool.Name),
			MinSize:       types.Int64Value(nodePool.MinSize),
			MaxSize:       types.Int64Value(nodePool.MaxSize),
			InstanceType:  types.StringValue(nodePool.InstanceType),
			Labels:        clusterNodePoolLabels,
			Taints:        clusterNodePoolTaintsFrom(nodePool, ctx, d),
			CloudSettings: clusterNodePoolCloudSettings,
		}).Element()
	}

	return result
}

func clusterNodePoolTaintsFrom(nodePool *console.NodePoolFragment, ctx context.Context, d diag.Diagnostics) types.Set {
	taints := make([]attr.Value, len(nodePool.Taints))
	for i, taint := range nodePool.Taints {
		objValue, diags := types.ObjectValueFrom(ctx, NodePoolTaintAttrTypes, NodePoolTaint{
			Key:    types.StringValue(taint.Key),
			Value:  types.StringValue(taint.Value),
			Effect: types.StringValue(taint.Effect),
		})
		taints[i] = objValue
		d.Append(diags...)
	}

	setValue, diags := types.SetValue(basetypes.ObjectType{AttrTypes: NodePoolTaintAttrTypes}, taints)
	d.Append(diags...)
	return setValue
}
