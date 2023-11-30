package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console-client-go"
)

type ClusterNodePool struct {
	Name          types.String           `tfsdk:"name"`
	MinSize       types.Int64            `tfsdk:"min_size"`
	MaxSize       types.Int64            `tfsdk:"max_size"`
	InstanceType  types.String           `tfsdk:"instance_type"`
	Labels        types.Map              `tfsdk:"labels"`
	Taints        []NodePoolTaint        `tfsdk:"taints"`
	CloudSettings *NodePoolCloudSettings `tfsdk:"cloud_settings"`
}

var ClusterNodePoolAttrTypes = map[string]attr.Type{
	"name":           types.StringType,
	"min_size":       types.Int64Type,
	"max_size":       types.Int64Type,
	"instance_type":  types.StringType,
	"labels":         types.MapType{ElemType: types.StringType},
	"taints":         types.ListType{ElemType: types.ObjectType{AttrTypes: NodePoolTaintAttrTypes}},
	"cloud_settings": types.ObjectType{AttrTypes: NodePoolCloudSettingsAttrTypes},
}

func (c *ClusterNodePool) LabelsAttribute(ctx context.Context, d diag.Diagnostics) map[string]interface{} {
	elements := make(map[string]types.String, len(c.Labels.Elements()))
	d.Append(c.Labels.ElementsAs(ctx, &elements, false)...)

	return ToAttributesMap(elements)
}

func (c *ClusterNodePool) TaintsAttribute() []*console.TaintAttributes {
	result := make([]*console.TaintAttributes, 0, len(c.Taints))
	for _, np := range c.Taints {
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
		"name":           c.Name,
		"min_size":       c.MinSize,
		"max_size":       c.MaxSize,
		"instance_type":  c.InstanceType,
		"labels":         c.TerraformAttributesLabels(),
		"taints":         c.TerraformAttributesTaints(),
		"cloud_settings": c.TerraformAttributesCloudSettings(),
	}
}

func (c *ClusterNodePool) TerraformAttributesLabels() attr.Value {
	if c.Labels.IsNull() {
		return types.MapNull(types.StringType)
	}

	return types.MapValueMust(types.StringType, c.Labels.Elements())
}

func (c *ClusterNodePool) TerraformAttributesTaints() attr.Value {
	if len(c.Taints) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: NodePoolTaintAttrTypes})
	}

	taints := make([]attr.Value, len(c.Taints))
	for i, taint := range c.Taints {
		taints[i] = types.ObjectValueMust(NodePoolTaintAttrTypes,
			map[string]attr.Value{"key": taint.Key, "value": taint.Value, "effect": taint.Effect})
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: NodePoolTaintAttrTypes}, taints)
}

func (c *ClusterNodePool) TerraformAttributesCloudSettings() attr.Value {
	if c.CloudSettings == nil {
		return types.ObjectNull(NodePoolCloudSettingsAttrTypes)
	}

	if c.CloudSettings.AWS != nil {
		return types.ObjectValueMust(NodePoolCloudSettingsAttrTypes,
			map[string]attr.Value{"aws": types.ObjectValueMust(NodePoolCloudSettingsAWSAttrTypes,
				map[string]attr.Value{"launch_template_id": c.CloudSettings.AWS.LaunchTemplateId})})
	}

	return types.ObjectNull(NodePoolCloudSettingsAttrTypes)
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

func ClusterNodePoolsFrom(nodePools []*console.NodePoolFragment) []*ClusterNodePool {
	result := make([]*ClusterNodePool, len(nodePools))
	for i, nodePool := range nodePools {
		result[i] = &ClusterNodePool{
			Name:          types.StringValue(nodePool.Name),
			MinSize:       types.Int64Value(nodePool.MinSize),
			MaxSize:       types.Int64Value(nodePool.MaxSize),
			InstanceType:  types.StringValue(nodePool.InstanceType),
			Labels:        clusterNodePoolLabelsFrom(nodePool),
			Taints:        clusterNodePoolTaintsFrom(nodePool),
			CloudSettings: nil,
		}
	}

	return result
}

func clusterNodePoolLabelsFrom(nodePool *console.NodePoolFragment) types.Map {
	labels := make(map[string]attr.Value)
	for k, v := range nodePool.Labels {
		labels[k] = types.StringValue(v.(string))
	}

	return types.MapValueMust(types.StringType, labels)
}

func clusterNodePoolTaintsFrom(nodePool *console.NodePoolFragment) []NodePoolTaint {
	taints := make([]NodePoolTaint, 0)
	for _, taint := range nodePool.Taints {
		taints = append(taints, NodePoolTaint{
			Key:    types.StringValue(taint.Key),
			Value:  types.StringValue(taint.Value),
			Effect: types.StringValue(taint.Effect),
		})
	}
	return taints
}
