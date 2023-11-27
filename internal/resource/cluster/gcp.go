package cluster

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

func GCPCloudSettingsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"project": schema.StringAttribute{
				Required:            true,
				Description:         "GCP project id to deploy cluster to.",
				MarkdownDescription: "GCP project id to deploy cluster to.",
			},
			"network": schema.StringAttribute{
				Required:            true,
				Description:         "GCP network id to use when creating the cluster.",
				MarkdownDescription: "GCP network id to use when creating the cluster.",
			},
			"region": schema.StringAttribute{
				Required:            true,
				Description:         "GCP region to deploy cluster to.",
				MarkdownDescription: "GCP region to deploy cluster to.",
			},
		},
	}
}
