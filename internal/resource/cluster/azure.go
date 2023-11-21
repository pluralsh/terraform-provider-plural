package cluster

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

func AzureCloudSettingsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"resource_group": schema.StringAttribute{
				MarkdownDescription: "Name of the Azure resource group for this cluster.",
				Required:            true,
			},
			"network": schema.StringAttribute{
				MarkdownDescription: "Name of the Azure virtual network for this cluster.",
				Required:            true,
			},
			"subscription_id": schema.StringAttribute{
				MarkdownDescription: "GUID of the Azure subscription to hold this cluster.",
				Required:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "String matching one of the canonical Azure region names, i.e. eastus.",
				Required:            true,
			},
		},
	}
}
