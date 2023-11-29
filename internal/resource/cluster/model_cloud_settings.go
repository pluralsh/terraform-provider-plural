package cluster

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console-client-go"
)

type ClusterCloudSettings struct {
	AWS   *ClusterCloudSettingsAWS   `tfsdk:"aws"`
	Azure *ClusterCloudSettingsAzure `tfsdk:"azure"`
	GCP   *ClusterCloudSettingsGCP   `tfsdk:"gcp"`
	BYOK  *ClusterCloudSettingsBYOK  `tfsdk:"byok"`
}

func (c *ClusterCloudSettings) Attributes() *console.CloudSettingsAttributes {
	if c == nil {
		return nil
	}

	if c.AWS != nil {
		return &console.CloudSettingsAttributes{Aws: c.AWS.Attributes()}
	}

	if c.Azure != nil {
		return &console.CloudSettingsAttributes{Azure: c.Azure.Attributes()}
	}

	if c.GCP != nil {
		return &console.CloudSettingsAttributes{Gcp: c.GCP.Attributes()}
	}

	return nil
}

type ClusterCloudSettingsAWS struct {
	Region types.String `tfsdk:"region"`
}

func (c *ClusterCloudSettingsAWS) Attributes() *console.AwsCloudAttributes {
	return &console.AwsCloudAttributes{
		Region: c.Region.ValueStringPointer(),
	}
}

type ClusterCloudSettingsAzure struct {
	ResourceGroup  types.String `tfsdk:"resource_group"`
	Network        types.String `tfsdk:"network"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
	Location       types.String `tfsdk:"location"`
}

func (c *ClusterCloudSettingsAzure) Attributes() *console.AzureCloudAttributes {
	return &console.AzureCloudAttributes{
		Location:       c.Location.ValueStringPointer(),
		SubscriptionID: c.SubscriptionId.ValueStringPointer(),
		ResourceGroup:  c.ResourceGroup.ValueStringPointer(),
		Network:        c.Network.ValueStringPointer(),
	}
}

type ClusterCloudSettingsGCP struct {
	Region  types.String `tfsdk:"region"`
	Network types.String `tfsdk:"network"`
	Project types.String `tfsdk:"project"`
}

func (c *ClusterCloudSettingsGCP) Attributes() *console.GcpCloudAttributes {
	return &console.GcpCloudAttributes{
		Project: c.Project.ValueStringPointer(),
		Network: c.Network.ValueStringPointer(),
		Region:  c.Region.ValueStringPointer(),
	}
}

type ClusterCloudSettingsBYOK struct {
	Kubeconfig Kubeconfig `tfsdk:"kubeconfig"`
}

type Kubeconfig struct {
	Host                  types.String    `tfsdk:"host"`
	Username              types.String    `tfsdk:"username"`
	Password              types.String    `tfsdk:"password"`
	Insecure              types.Bool      `tfsdk:"insecure"`
	TlsServerName         types.String    `tfsdk:"tls_server_name"`
	ClientCertificate     types.String    `tfsdk:"client_certificate"`
	ClientKey             types.String    `tfsdk:"client_key"`
	ClusterCACertificate  types.String    `tfsdk:"cluster_ca_certificate"`
	ConfigPath            types.String    `tfsdk:"config_path"`
	ConfigContext         types.String    `tfsdk:"config_context"`
	ConfigContextAuthInfo types.String    `tfsdk:"config_context_auth_info"`
	ConfigContextCluster  types.String    `tfsdk:"config_context_cluster"`
	Token                 types.String    `tfsdk:"token"`
	ProxyURL              types.String    `tfsdk:"proxy_url"`
	Exec                  *KubeconfigExec `tfsdk:"exec"`
}

type KubeconfigExec struct {
	Command    types.String `tfsdk:"command"`
	Args       types.String `tfsdk:"args"`
	Env        types.String `tfsdk:"env"`
	APIVersion types.String `tfsdk:"api_version"`
}
