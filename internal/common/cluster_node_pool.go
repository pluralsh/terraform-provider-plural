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

func (c *ClusterNodePool) TerraformTypes() map[string]attr.Type {
	return map[string]attr.Type{
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
		"cloud_settings": types.ObjectType{AttrTypes: map[string]attr.Type{
			"aws": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"launch_template_id": types.StringType,
				},
			},
		}},
	}
}

func (c *ClusterNodePool) TerraformAttributes() map[string]attr.Value {
	return map[string]attr.Value{
		"name":          c.Name,
		"min_size":      c.MinSize,
		"max_size":      c.MaxSize,
		"instance_type": c.InstanceType,
		"labels":        types.MapNull(types.StringType),
		"taints": types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{
			"key":    types.StringType,
			"value":  types.StringType,
			"effect": types.StringType,
		}}),
		"cloud_settings": types.ObjectNull(map[string]attr.Type{
			"aws": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"launch_template_id": types.StringType,
				},
			},
		}),
	}
}

func (c *ClusterNodePool) Element() attr.Value {
	return types.ObjectValueMust(c.TerraformTypes(), c.TerraformAttributes())
}

type NodePoolTaint struct {
	Key    types.String `tfsdk:"key"`
	Value  types.String `tfsdk:"value"`
	Effect types.String `tfsdk:"effect"`
}

type NodePoolCloudSettings struct {
	AWS *NodePoolCloudSettingsAWS `tfsdk:"aws"`
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

func (c *NodePoolCloudSettingsAWS) Attributes() *console.AwsNodeCloudAttributes {
	return &console.AwsNodeCloudAttributes{
		LaunchTemplateID: c.LaunchTemplateId.ValueStringPointer(),
	}
}

func ClusterNodePoolsFrom(nodepools []*console.NodePoolFragment) []*ClusterNodePool {
	result := make([]*ClusterNodePool, 0, len(nodepools))
	for _, np := range nodepools {
		result = append(result, &ClusterNodePool{
			Name:         types.StringValue(np.Name),
			MinSize:      types.Int64Value(np.MinSize),
			MaxSize:      types.Int64Value(np.MaxSize),
			InstanceType: types.StringValue(np.InstanceType),
		})
	}

	return result
}
