package cluster

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

func GCPCloudSettingsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"project": schema.StringAttribute{
				MarkdownDescription: "",
				Required:            true,
			},
			"network": schema.StringAttribute{
				MarkdownDescription: "",
				Required:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "",
				Required:            true,
			},
		},
	}
}
