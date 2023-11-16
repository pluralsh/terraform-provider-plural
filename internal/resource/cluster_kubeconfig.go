package resource

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func kubeconfigAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The address of the Kubernetes clusters", // KUBE_HOST
			},
			"username": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The username for basic authentication to the Kubernetes cluster.", // KUBE_USER
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The password for basic authentication to the Kubernetes cluster.", // KUBE_PASSWORD
			},
			"insecure": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Skips the validity check for the server's certificate. This will make your HTTPS connections insecure.", // KUBE_INSECURE
			},
			"tls_server_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "TLS server name is used to check server certificate. If it is empty, the hostname used to contact the server is used.", // KUBE_TLS_SERVER_NAME
			},
			"client_certificate": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The path to a client cert file for TLS.", // KUBE_CLIENT_CERT_DATA
			},
			"client_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The path to a client key file for TLS.", // KUBE_CLIENT_KEY_DATA
			},
			"cluster_ca_certificate": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The path to a cert file for the certificate authority.", // KUBE_CLUSTER_CA_CERT_DATA
			},
			"config_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Path to the kubeconfig file.", // KUBE_CONFIG_PATH
			},
			"config_context": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "kubeconfig context to use.", // KUBE_CTX
			},
			"config_context_auth_info": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "", // KUBE_CTX_AUTH_INFO
			},
			"config_context_cluster": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "", // KUBE_CTX_CLUSTER
			},
			"token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Token is the bearer token for authentication to the Kubernetes cluster.", // KUBE_TOKEN
			},
			"proxy_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The URL to the proxy to be used for all requests made by this client.", // KUBE_PROXY_URL
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
	}
}
