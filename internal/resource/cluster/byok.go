package cluster

import (
	"terraform-provider-plural/internal/defaults"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func BYOKCloudSettingsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"kubeconfig": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"host": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             defaults.Env("PLURAL_KUBE_HOST", ""),
						MarkdownDescription: "The complete address of the Kubernetes cluster, using scheme://hostname:port format. Can be sourced from `PLURAL_KUBE_HOST`.",
					},
					"username": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             defaults.Env("PLURAL_KUBE_USER", ""),
						MarkdownDescription: "The username for basic authentication to the Kubernetes cluster. Can be sourced from `PLURAL_KUBE_USER`.",
					},
					"password": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Sensitive:           true,
						Default:             defaults.Env("PLURAL_KUBE_PASSWORD", ""),
						MarkdownDescription: "The password for basic authentication to the Kubernetes cluster. Can be sourced from `PLURAL_KUBE_PASSWORD`.",
					},
					"insecure": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             defaults.Env("PLURAL_KUBE_INSECURE", false),
						MarkdownDescription: "Skips the validity check for the server's certificate. This will make your HTTPS connections insecure. Can be sourced from `PLURAL_KUBE_INSECURE`.",
					},
					"tls_server_name": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             defaults.Env("PLURAL_KUBE_TLS_SERVER_NAME", ""),
						MarkdownDescription: "TLS server name is used to check server certificate. If it is empty, the hostname used to contact the server is used. Can be sourced from `PLURAL_KUBE_TLS_SERVER_NAME`.",
					},
					"client_certificate": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             defaults.Env("PLURAL_KUBE_CLIENT_CERT_DATA", ""),
						MarkdownDescription: "The path to a client cert file for TLS. Can be sourced from `PLURAL_KUBE_CLIENT_CERT_DATA`.",
					},
					"client_key": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Sensitive:           true,
						Default:             defaults.Env("PLURAL_KUBE_CLIENT_KEY_DATA", ""),
						MarkdownDescription: "The path to a client key file for TLS. Can be sourced from `PLURAL_KUBE_CLIENT_KEY_DATA`.",
					},
					"cluster_ca_certificate": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             defaults.Env("PLURAL_KUBE_CLUSTER_CA_CERT_DATA", ""),
						MarkdownDescription: "The path to a cert file for the certificate authority. Can be sourced from `PLURAL_KUBE_CLUSTER_CA_CERT_DATA`.",
					},
					"config_path": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             defaults.Env("PLURAL_KUBE_CONFIG_PATH", ""),
						MarkdownDescription: "Path to the kubeconfig file. Can be sourced from `PLURAL_KUBE_CONFIG_PATH`.",
					},
					"config_context": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             defaults.Env("PLURAL_KUBE_CTX", ""),
						MarkdownDescription: "kubeconfig context to use. Can be sourced from `PLURAL_KUBE_CTX`.",
					},
					"config_context_auth_info": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             defaults.Env("PLURAL_KUBE_CTX_AUTH_INFO", ""),
						MarkdownDescription: "Can be sourced from `PLURAL_KUBE_CTX_AUTH_INFO`.",
					},
					"config_context_cluster": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             defaults.Env("PLURAL_KUBE_CTX_CLUSTER", ""),
						MarkdownDescription: "Can be sourced from `PLURAL_KUBE_CTX_CLUSTER`.",
					},
					"token": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Sensitive:           true,
						Default:             defaults.Env("PLURAL_KUBE_TOKEN", ""),
						MarkdownDescription: "Token is the bearer token for authentication to the Kubernetes cluster. Can be sourced from `PLURAL_KUBE_TOKEN`.",
					},
					"proxy_url": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             defaults.Env("PLURAL_KUBE_PROXY_URL", ""),
						MarkdownDescription: "The URL to the proxy to be used for all requests made by this client. Can be sourced from `PLURAL_KUBE_PROXY_URL`.",
					},
					"exec": schema.ListNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Specifies a command to provide client credentials",
						Validators:          []validator.List{listvalidator.SizeAtMost(1)},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"command": schema.StringAttribute{
									MarkdownDescription: "Command to execute.",
									Required:            true,
								},
								"args": schema.ListAttribute{
									MarkdownDescription: "Arguments to pass to the command when executing it.",
									Optional:            true,
									ElementType:         types.StringType,
								},
								"env": schema.MapAttribute{
									MarkdownDescription: "Defines  environment variables to expose to the process.",
									Optional:            true,
									ElementType:         types.StringType,
								},
								"api_version": schema.StringAttribute{
									MarkdownDescription: "Preferred input version.",
									Required:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}
