package common

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

func (k *Kubeconfig) Unchanged(other *Kubeconfig) bool {
	if k == nil {
		return other == nil
	}

	if other == nil {
		return false
	}

	return k.Host == other.Host
}

type KubeconfigExec struct {
	Command    types.String `tfsdk:"command"`
	Args       types.List   `tfsdk:"args"`
	Env        types.Map    `tfsdk:"env"`
	APIVersion types.String `tfsdk:"api_version"`
}
