package model

import "github.com/hashicorp/terraform-plugin-framework/types"

// Provider describes the provider resource and data source model.
type Provider struct {
	Id            types.String          `tfsdk:"id"`
	Name          types.String          `tfsdk:"name"`
	Cloud         types.String          `tfsdk:"cloud"`
	CloudSettings ProviderCloudSettings `tfsdk:"cloud_settings"`
}

type ProviderCloudSettings struct {
	AWS   ProviderCloudSettingsAWS   `tfsdk:"aws"`
	Azure ProviderCloudSettingsAzure `tfsdk:"azure"`
	GCP   ProviderCloudSettingsGCP   `tfsdk:"gcp"`
}

type ProviderCloudSettingsAWS struct {
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
}

type ProviderCloudSettingsAzure struct {
	SubscriptionId types.String `tfsdk:"subscription_id"`
	TenantId       types.String `tfsdk:"tenant_id"`
	ClientId       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
}

type ProviderCloudSettingsGCP struct {
	Credentials types.String `tfsdk:"credentials"`
}
