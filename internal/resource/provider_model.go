package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
)

type provider struct {
	Id            types.String          `tfsdk:"id"`
	Name          types.String          `tfsdk:"name"`
	Namespace     types.String          `tfsdk:"namespace"`
	Editable      types.Bool            `tfsdk:"editable"`
	Cloud         types.String          `tfsdk:"cloud"`
	CloudSettings ProviderCloudSettings `tfsdk:"cloud_settings"`
}

func (p *provider) From(cp *console.ClusterProviderFragment) {
	p.Id = types.StringValue(cp.ID)
	p.Name = types.StringValue(cp.Name)
	p.Namespace = types.StringValue(cp.Namespace)
	p.Editable = types.BoolPointerValue(cp.Editable)
	p.Cloud = types.StringValue(cp.Cloud)
}

func (p *provider) Attributes() console.ClusterProviderAttributes {
	return console.ClusterProviderAttributes{
		Name:          p.Name.ValueString(),
		Namespace:     p.Namespace.ValueStringPointer(),
		Cloud:         p.Cloud.ValueStringPointer(),
		CloudSettings: p.CloudSettings.Attributes(),
	}
}

func (p *provider) UpdateAttributes() console.ClusterProviderUpdateAttributes {
	return console.ClusterProviderUpdateAttributes{
		CloudSettings: p.CloudSettings.Attributes(),
	}
}

type ProviderCloudSettings struct {
	AWS   *ProviderCloudSettingsAWS   `tfsdk:"aws"`
	Azure *ProviderCloudSettingsAzure `tfsdk:"azure"`
	GCP   *ProviderCloudSettingsGCP   `tfsdk:"gcp"`
}

func (p *ProviderCloudSettings) Attributes() *console.CloudProviderSettingsAttributes {
	if p == nil {
		return nil
	}

	if p.AWS != nil {
		return &console.CloudProviderSettingsAttributes{Aws: p.AWS.Attributes()}
	}

	if p.Azure != nil {
		return &console.CloudProviderSettingsAttributes{Azure: p.Azure.Attributes()}
	}

	if p.GCP != nil {
		return &console.CloudProviderSettingsAttributes{Gcp: p.GCP.Attributes()}
	}

	return nil
}

type ProviderCloudSettingsAWS struct {
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
}

func (p *ProviderCloudSettingsAWS) Attributes() *console.AwsSettingsAttributes {
	return &console.AwsSettingsAttributes{
		AccessKeyID:     p.AccessKeyId.ValueString(),
		SecretAccessKey: p.SecretAccessKey.ValueString(),
	}
}

type ProviderCloudSettingsAzure struct {
	SubscriptionId types.String `tfsdk:"subscription_id"`
	TenantId       types.String `tfsdk:"tenant_id"`
	ClientId       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
}

func (p *ProviderCloudSettingsAzure) Attributes() *console.AzureSettingsAttributes {
	return &console.AzureSettingsAttributes{
		SubscriptionID: p.SubscriptionId.ValueString(),
		TenantID:       p.TenantId.ValueString(),
		ClientID:       p.ClientId.ValueString(),
		ClientSecret:   p.ClientSecret.ValueString(),
	}
}

type ProviderCloudSettingsGCP struct {
	Credentials types.String `tfsdk:"credentials"`
}

func (p *ProviderCloudSettingsGCP) Attributes() *console.GcpSettingsAttributes {
	return &console.GcpSettingsAttributes{
		ApplicationCredentials: p.Credentials.ValueString(),
	}
}
