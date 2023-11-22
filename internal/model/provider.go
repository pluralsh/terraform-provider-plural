package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console-client-go"
)

// Provider describes the provider resource and data source model.
type Provider struct {
	Id            types.String          `tfsdk:"id"`
	Name          types.String          `tfsdk:"name"`
	Namespace     types.String          `tfsdk:"namespace"`
	Editable      types.Bool            `tfsdk:"editable"`
	Cloud         types.String          `tfsdk:"cloud"`
	CloudSettings ProviderCloudSettings `tfsdk:"cloud_settings"`
}

type ProviderCloudSettings struct {
	AWS   *ProviderCloudSettingsAWS   `tfsdk:"aws"`
	Azure *ProviderCloudSettingsAzure `tfsdk:"azure"`
	GCP   *ProviderCloudSettingsGCP   `tfsdk:"gcp"`
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

func (p *Provider) CloudProviderSettingsAttributes() *console.CloudProviderSettingsAttributes {
	if IsCloud(p.Cloud.ValueString(), CloudAWS) {
		return &console.CloudProviderSettingsAttributes{
			Aws: &console.AwsSettingsAttributes{
				AccessKeyID:     p.CloudSettings.AWS.AccessKeyId.ValueString(),
				SecretAccessKey: p.CloudSettings.AWS.SecretAccessKey.ValueString(),
			},
		}
	}

	if IsCloud(p.Cloud.ValueString(), CloudAzure) {
		return &console.CloudProviderSettingsAttributes{
			Azure: &console.AzureSettingsAttributes{
				SubscriptionID: p.CloudSettings.Azure.SubscriptionId.ValueString(),
				TenantID:       p.CloudSettings.Azure.TenantId.ValueString(),
				ClientID:       p.CloudSettings.Azure.ClientId.ValueString(),
				ClientSecret:   p.CloudSettings.Azure.ClientSecret.ValueString(),
			},
		}
	}

	if IsCloud(p.Cloud.ValueString(), CloudGCP) {
		return &console.CloudProviderSettingsAttributes{
			Gcp: &console.GcpSettingsAttributes{
				ApplicationCredentials: p.CloudSettings.GCP.Credentials.ValueString(),
			},
		}
	}

	return nil
}

func (p *Provider) CreateAttributes() console.ClusterProviderAttributes {
	return console.ClusterProviderAttributes{
		Name:          p.Name.ValueString(),
		Namespace:     p.Namespace.ValueStringPointer(),
		Cloud:         p.Cloud.ValueStringPointer(),
		CloudSettings: p.CloudProviderSettingsAttributes(),
	}
}

func (p *Provider) UpdateAttributes() console.ClusterProviderUpdateAttributes {
	return console.ClusterProviderUpdateAttributes{
		CloudSettings: p.CloudProviderSettingsAttributes(),
	}
}

func (p *Provider) From(cp *console.ClusterProviderFragment) {
	p.Id = types.StringValue(cp.ID)
	p.Name = types.StringValue(cp.Name)
	p.Namespace = types.StringValue(cp.Namespace)
	p.Editable = types.BoolPointerValue(cp.Editable)
	p.Cloud = types.StringValue(cp.Cloud)
}
