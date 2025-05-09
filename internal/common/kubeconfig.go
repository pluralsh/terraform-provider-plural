package common

import (
	"terraform-provider-plural/internal/defaults"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

func KubeconfigProviderSchema() providerschema.SingleNestedAttribute {
	return providerschema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]providerschema.Attribute{
			"host": providerschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_HOST", ""),
				Description:         "The complete address of the Kubernetes cluster, using scheme://hostname:port format. Can be sourced from PLURAL_KUBE_HOST.",
				MarkdownDescription: "The complete address of the Kubernetes cluster, using scheme://hostname:port format. Can be sourced from `PLURAL_KUBE_HOST`.",
			},
			"username": providerschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_USER", ""),
				Description:         "The username for basic authentication to the Kubernetes cluster. Can be sourced from PLURAL_KUBE_USER.",
				MarkdownDescription: "The username for basic authentication to the Kubernetes cluster. Can be sourced from `PLURAL_KUBE_USER`.",
			},
			"password": providerschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Default:             defaults.Env("PLURAL_KUBE_PASSWORD", ""),
				Description:         "The password for basic authentication to the Kubernetes cluster. Can be sourced from PLURAL_KUBE_PASSWORD.",
				MarkdownDescription: "The password for basic authentication to the Kubernetes cluster. Can be sourced from `PLURAL_KUBE_PASSWORD`.",
			},
			"insecure": providerschema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_INSECURE", false),
				Description:         "Skips the validity check for the server's certificate. This will make your HTTPS connections insecure. Can be sourced from PLURAL_KUBE_INSECURE.",
				MarkdownDescription: "Skips the validity check for the server's certificate. This will make your HTTPS connections insecure. Can be sourced from `PLURAL_KUBE_INSECURE`.",
			},
			"tls_server_name": providerschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_TLS_SERVER_NAME", ""),
				Description:         "TLS server name is used to check server certificate. If it is empty, the hostname used to contact the server is used. Can be sourced from PLURAL_KUBE_TLS_SERVER_NAME.",
				MarkdownDescription: "TLS server name is used to check server certificate. If it is empty, the hostname used to contact the server is used. Can be sourced from `PLURAL_KUBE_TLS_SERVER_NAME`.",
			},
			"client_certificate": providerschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_CLIENT_CERT_DATA", ""),
				Description:         "The path to a client cert file for TLS. Can be sourced from PLURAL_KUBE_CLIENT_CERT_DATA.",
				MarkdownDescription: "The path to a client cert file for TLS. Can be sourced from `PLURAL_KUBE_CLIENT_CERT_DATA`.",
			},
			"client_key": providerschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Default:             defaults.Env("PLURAL_KUBE_CLIENT_KEY_DATA", ""),
				Description:         "The path to a client key file for TLS. Can be sourced from PLURAL_KUBE_CLIENT_KEY_DATA.",
				MarkdownDescription: "The path to a client key file for TLS. Can be sourced from `PLURAL_KUBE_CLIENT_KEY_DATA`.",
			},
			"cluster_ca_certificate": providerschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_CLUSTER_CA_CERT_DATA", ""),
				Description:         "The path to a cert file for the certificate authority. Can be sourced from PLURAL_KUBE_CLUSTER_CA_CERT_DATA.",
				MarkdownDescription: "The path to a cert file for the certificate authority. Can be sourced from `PLURAL_KUBE_CLUSTER_CA_CERT_DATA`.",
			},
			"config_path": providerschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_CONFIG_PATH", ""),
				Description:         "Path to the kubeconfig file. Can be sourced from PLURAL_KUBE_CONFIG_PATH.",
				MarkdownDescription: "Path to the kubeconfig file. Can be sourced from `PLURAL_KUBE_CONFIG_PATH`.",
			},
			"config_context": providerschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_CTX", ""),
				Description:         "kubeconfig context to use. Can be sourced from PLURAL_KUBE_CTX.",
				MarkdownDescription: "kubeconfig context to use. Can be sourced from `PLURAL_KUBE_CTX`.",
			},
			"config_context_auth_info": providerschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_CTX_AUTH_INFO", ""),
				Description:         "Can be sourced from PLURAL_KUBE_CTX_AUTH_INFO.",
				MarkdownDescription: "Can be sourced from `PLURAL_KUBE_CTX_AUTH_INFO`.",
			},
			"config_context_cluster": providerschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_CTX_CLUSTER", ""),
				Description:         "Can be sourced from PLURAL_KUBE_CTX_CLUSTER.",
				MarkdownDescription: "Can be sourced from `PLURAL_KUBE_CTX_CLUSTER`.",
			},
			"token": providerschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Default:             defaults.Env("PLURAL_KUBE_TOKEN", ""),
				Description:         "Token is the bearer token for authentication to the Kubernetes cluster. Can be sourced from PLURAL_KUBE_TOKEN.",
				MarkdownDescription: "Token is the bearer token for authentication to the Kubernetes cluster. Can be sourced from `PLURAL_KUBE_TOKEN`.",
			},
			"proxy_url": providerschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_PROXY_URL", ""),
				Description:         "The URL to the proxy to be used for all requests made by this client. Can be sourced from PLURAL_KUBE_PROXY_URL.",
				MarkdownDescription: "The URL to the proxy to be used for all requests made by this client. Can be sourced from `PLURAL_KUBE_PROXY_URL`.",
			},
			"exec": providerschema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies a command to provide client credentials",
				Validators:          []validator.List{listvalidator.SizeAtMost(1)},
				NestedObject: providerschema.NestedAttributeObject{
					Attributes: map[string]providerschema.Attribute{
						"command": providerschema.StringAttribute{
							Description:         "Command to execute.",
							MarkdownDescription: "Command to execute.",
							Required:            true,
						},
						"args": providerschema.ListAttribute{
							Description:         "Arguments to pass to the command when executing it.",
							MarkdownDescription: "Arguments to pass to the command when executing it.",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"env": providerschema.MapAttribute{
							Description:         "Defines environment variables to expose to the process.",
							MarkdownDescription: "Defines environment variables to expose to the process.",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"api_version": providerschema.StringAttribute{
							Description:         "Preferred input version.",
							MarkdownDescription: "Preferred input version.",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func KubeconfigResourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		DeprecationMessage: "kubeconfig configuration has been moved to the provider.",
		Optional:           true,
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_HOST", ""),
				Description:         "The complete address of the Kubernetes cluster, using scheme://hostname:port format. Can be sourced from PLURAL_KUBE_HOST.",
				MarkdownDescription: "The complete address of the Kubernetes cluster, using scheme://hostname:port format. Can be sourced from `PLURAL_KUBE_HOST`.",
			},
			"username": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_USER", ""),
				Description:         "The username for basic authentication to the Kubernetes cluster. Can be sourced from PLURAL_KUBE_USER.",
				MarkdownDescription: "The username for basic authentication to the Kubernetes cluster. Can be sourced from `PLURAL_KUBE_USER`.",
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Default:             defaults.Env("PLURAL_KUBE_PASSWORD", ""),
				Description:         "The password for basic authentication to the Kubernetes cluster. Can be sourced from PLURAL_KUBE_PASSWORD.",
				MarkdownDescription: "The password for basic authentication to the Kubernetes cluster. Can be sourced from `PLURAL_KUBE_PASSWORD`.",
			},
			"insecure": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_INSECURE", false),
				Description:         "Skips the validity check for the server's certificate. This will make your HTTPS connections insecure. Can be sourced from PLURAL_KUBE_INSECURE.",
				MarkdownDescription: "Skips the validity check for the server's certificate. This will make your HTTPS connections insecure. Can be sourced from `PLURAL_KUBE_INSECURE`.",
			},
			"tls_server_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_TLS_SERVER_NAME", ""),
				Description:         "TLS server name is used to check server certificate. If it is empty, the hostname used to contact the server is used. Can be sourced from PLURAL_KUBE_TLS_SERVER_NAME.",
				MarkdownDescription: "TLS server name is used to check server certificate. If it is empty, the hostname used to contact the server is used. Can be sourced from `PLURAL_KUBE_TLS_SERVER_NAME`.",
			},
			"client_certificate": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_CLIENT_CERT_DATA", ""),
				Description:         "The path to a client cert file for TLS. Can be sourced from PLURAL_KUBE_CLIENT_CERT_DATA.",
				MarkdownDescription: "The path to a client cert file for TLS. Can be sourced from `PLURAL_KUBE_CLIENT_CERT_DATA`.",
			},
			"client_key": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Default:             defaults.Env("PLURAL_KUBE_CLIENT_KEY_DATA", ""),
				Description:         "The path to a client key file for TLS. Can be sourced from PLURAL_KUBE_CLIENT_KEY_DATA.",
				MarkdownDescription: "The path to a client key file for TLS. Can be sourced from `PLURAL_KUBE_CLIENT_KEY_DATA`.",
			},
			"cluster_ca_certificate": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_CLUSTER_CA_CERT_DATA", ""),
				Description:         "The path to a cert file for the certificate authority. Can be sourced from PLURAL_KUBE_CLUSTER_CA_CERT_DATA.",
				MarkdownDescription: "The path to a cert file for the certificate authority. Can be sourced from `PLURAL_KUBE_CLUSTER_CA_CERT_DATA`.",
			},
			"config_path": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_CONFIG_PATH", ""),
				Description:         "Path to the kubeconfig file. Can be sourced from PLURAL_KUBE_CONFIG_PATH.",
				MarkdownDescription: "Path to the kubeconfig file. Can be sourced from `PLURAL_KUBE_CONFIG_PATH`.",
			},
			"config_context": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_CTX", ""),
				Description:         "kubeconfig context to use. Can be sourced from PLURAL_KUBE_CTX.",
				MarkdownDescription: "kubeconfig context to use. Can be sourced from `PLURAL_KUBE_CTX`.",
			},
			"config_context_auth_info": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_CTX_AUTH_INFO", ""),
				Description:         "Can be sourced from PLURAL_KUBE_CTX_AUTH_INFO.",
				MarkdownDescription: "Can be sourced from `PLURAL_KUBE_CTX_AUTH_INFO`.",
			},
			"config_context_cluster": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_CTX_CLUSTER", ""),
				Description:         "Can be sourced from PLURAL_KUBE_CTX_CLUSTER.",
				MarkdownDescription: "Can be sourced from `PLURAL_KUBE_CTX_CLUSTER`.",
			},
			"token": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Default:             defaults.Env("PLURAL_KUBE_TOKEN", ""),
				Description:         "Token is the bearer token for authentication to the Kubernetes cluster. Can be sourced from PLURAL_KUBE_TOKEN.",
				MarkdownDescription: "Token is the bearer token for authentication to the Kubernetes cluster. Can be sourced from `PLURAL_KUBE_TOKEN`.",
			},
			"proxy_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             defaults.Env("PLURAL_KUBE_PROXY_URL", ""),
				Description:         "The URL to the proxy to be used for all requests made by this client. Can be sourced from PLURAL_KUBE_PROXY_URL.",
				MarkdownDescription: "The URL to the proxy to be used for all requests made by this client. Can be sourced from `PLURAL_KUBE_PROXY_URL`.",
			},
			"exec": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies a command to provide client credentials",
				Validators:          []validator.List{listvalidator.SizeAtMost(1)},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"command": schema.StringAttribute{
							Description:         "Command to execute.",
							MarkdownDescription: "Command to execute.",
							Required:            true,
						},
						"args": schema.ListAttribute{
							Description:         "Arguments to pass to the command when executing it.",
							MarkdownDescription: "Arguments to pass to the command when executing it.",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"env": schema.MapAttribute{
							Description:         "Defines environment variables to expose to the process.",
							MarkdownDescription: "Defines environment variables to expose to the process.",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"api_version": schema.StringAttribute{
							Description:         "Preferred input version.",
							MarkdownDescription: "Preferred input version.",
							Required:            true,
						},
					},
				},
			},
		},
	}
}
