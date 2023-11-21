package model

import "github.com/hashicorp/terraform-plugin-framework/types"

// Cluster describes the cluster resource and data source model.
type Cluster struct {
	Id            types.String         `tfsdk:"id"`
	InseredAt     types.String         `tfsdk:"inserted_at"`
	Name          types.String         `tfsdk:"name"`
	Handle        types.String         `tfsdk:"handle"`
	Cloud         types.String         `tfsdk:"cloud"`
	CloudSettings ClusterCloudSettings `tfsdk:"cloud_settings"`
	Protect       types.Bool           `tfsdk:"protect"`
	Tags          types.Map            `tfsdk:"tags"`
	Kubeconfig    Kubeconfig           `tfsdk:"kubeconfig"`
}

type ClusterCloudSettings struct {
	AWS   ClusterCloudSettingsAWS   `tfsdk:"aws"`
	Azure ClusterCloudSettingsAzure `tfsdk:"azure"`
	GCP   ClusterCloudSettingsGCP   `tfsdk:"gcp"`
}

type ClusterCloudSettingsAWS struct {
	Region types.String `tfsdk:"region"`
}

type ClusterCloudSettingsAzure struct {
	ResourceGroup  types.String `tfsdk:"resource_group"`
	Network        types.String `tfsdk:"network"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
	Location       types.String `tfsdk:"location"`
}

type ClusterCloudSettingsGCP struct {
	Region  types.String `tfsdk:"region"`
	Network types.String `tfsdk:"network"`
	Project types.String `tfsdk:"project"`
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
