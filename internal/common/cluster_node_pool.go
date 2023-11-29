package common

import (
	"context"

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
	Taints        types.List             `tfsdk:"taints"`
	CloudSettings *NodePoolCloudSettings `tfsdk:"cloud_settings"`
}

func (c *ClusterNodePool) LabelsAttribute(ctx context.Context, d diag.Diagnostics) map[string]interface{} {
	elements := make(map[string]types.String, len(c.Labels.Elements()))
	d.Append(c.Labels.ElementsAs(ctx, &elements, false)...)

	return ToAttributesMap(elements)
}

type NodePoolCloudSettings struct {
	AWS *NodePoolCloudSettingsAWS `tfsdk:"aws"`
}

type NodePoolCloudSettingsAWS struct {
	LaunchTemplateId types.String `tfsdk:"launch_template_id"`
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
