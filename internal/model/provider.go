package model

import "github.com/hashicorp/terraform-plugin-framework/types"

// Provider describes the Provider resource and data source model.
type Provider struct {
	Id            types.String          `tfsdk:"id"`
	Name          types.String          `tfsdk:"name"`
	Cloud         types.String          `tfsdk:"cloud"`
	CloudSettings ProviderCloudSettings `tfsdk:"cloud_settings"`
}

type ProviderCloudSettings struct {
	AWS ProviderCloudSettingsAWS `tfsdk:"aws"`
}

type ProviderCloudSettingsAWS struct {
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
}
