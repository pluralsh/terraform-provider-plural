package cluster

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

func AWSCloudSettingsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				MarkdownDescription: "AWS region to deploy the cluster to.",
				Required:            true,
			},
		},
	}
}
