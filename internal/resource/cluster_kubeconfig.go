package resource

import (
	"bytes"
	"context"
	"sync"

	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mitchellh/go-homedir"
	"k8s.io/apimachinery/pkg/api/meta"
	apimachineryschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"terraform-provider-plural/internal/defaults"
)

func kubeconfigAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:            true,
				Default:             defaults.Env("PLURAL_KUBE_HOST", ""),
				MarkdownDescription: "The complete address of the Kubernetes cluster, using scheme://hostname:port format. Can be sourced from `PLURAL_KUBE_HOST`.",
			},
			"username": schema.StringAttribute{
				Optional:            true,
				Default:             defaults.Env("PLURAL_KUBE_HOST", ""),
				MarkdownDescription: "The username for basic authentication to the Kubernetes cluster. Can be sourced from `PLURAL_KUBE_USER`.",
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The password for basic authentication to the Kubernetes cluster. Can be sourced from `PLURAL_KUBE_PASSWORD`.",
			},
			"insecure": schema.BoolAttribute{
				Optional:            true,
				Default:             defaults.Env("PLURAL_KUBE_HOST", false),
				MarkdownDescription: "Skips the validity check for the server's certificate. This will make your HTTPS connections insecure. Can be sourced from `PLURAL_KUBE_INSECURE`.",
			},
			"tls_server_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "TLS server name is used to check server certificate. If it is empty, the hostname used to contact the server is used. Can be sourced from `PLURAL_KUBE_TLS_SERVER_NAME`.",
			},
			"client_certificate": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The path to a client cert file for TLS. Can be sourced from `PLURAL_KUBE_CLIENT_CERT_DATA`.",
			},
			"client_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The path to a client key file for TLS. Can be sourced from `PLURAL_KUBE_CLIENT_KEY_DATA`.",
			},
			"cluster_ca_certificate": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The path to a cert file for the certificate authority. Can be sourced from `PLURAL_KUBE_CLUSTER_CA_CERT_DATA`.",
			},
			"config_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Path to the kubeconfig file. Can be sourced from `PLURAL_KUBE_CONFIG_PATH`.",
			},
			"config_context": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "kubeconfig context to use. Can be sourced from `PLURAL_KUBE_CTX`.",
			},
			"config_context_auth_info": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Can be sourced from `PLURAL_KUBE_CTX_AUTH_INFO`.",
			},
			"config_context_cluster": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Can be sourced from `PLURAL_KUBE_CTX_CLUSTER`.",
			},
			"token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Token is the bearer token for authentication to the Kubernetes cluster. Can be sourced from `PLURAL_KUBE_TOKEN`.",
			},
			"proxy_url": schema.StringAttribute{
				Optional:            true,
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
	}
}

// KubeConfig is a RESTClientGetter interface implementation
type KubeConfig struct {
	ClientConfig clientcmd.ClientConfig

	sync.Mutex
}

// ToRESTConfig implemented interface method
func (k *KubeConfig) ToRESTConfig() (*rest.Config, error) {
	config, err := k.ToRawKubeConfigLoader().ClientConfig()
	return config, err
}

// ToDiscoveryClient implemented interface method
func (k *KubeConfig) ToDiscoveryClient() (discovery.DiscoveryInterface, error) {
	config, err := k.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	return discovery.NewDiscoveryClientForConfigOrDie(config), nil
}

// ToRESTMapper implemented interface method
func (k *KubeConfig) ToRESTMapper() (meta.RESTMapper, error) {
	discoveryClient, err := k.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient)
	return expander, nil
}

// ToRawKubeConfigLoader implemented interface method
func (k *KubeConfig) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return k.ClientConfig
}

func newKubeconfig(ctx context.Context, kubeconfig model.Kubeconfig, namespace *string) (*KubeConfig, error) {
	overrides := &clientcmd.ConfigOverrides{}
	loader := &clientcmd.ClientConfigLoadingRules{}

	if !kubeconfig.ConfigPath.IsNull() {
		tflog.Trace(ctx, "using kubeconfig", map[string]interface{}{
			"kubeconfig": kubeconfig.ConfigPath.ValueString(),
		})

		path, err := homedir.Expand(kubeconfig.ConfigPath.ValueString())
		if err != nil {
			return nil, err
		}
		loader.ExplicitPath = path

		if !kubeconfig.ConfigContext.IsNull() || !kubeconfig.ConfigContextAuthInfo.IsNull() || !kubeconfig.ConfigContextCluster.IsNull() {
			if !kubeconfig.ConfigContext.IsNull() {
				overrides.CurrentContext = kubeconfig.ConfigContext.ValueString()
				tflog.Trace(ctx, "using custom current context", map[string]interface{}{
					"context": overrides.CurrentContext,
				})
			}

			overrides.Context = clientcmdapi.Context{}
			if !kubeconfig.ConfigContextAuthInfo.IsNull() {
				overrides.Context.AuthInfo = kubeconfig.ConfigContextAuthInfo.ValueString()
			}
			if !kubeconfig.ConfigContextCluster.IsNull() {
				overrides.Context.Cluster = kubeconfig.ConfigContextCluster.ValueString()
			}
			tflog.Trace(ctx, "using overridden context", map[string]interface{}{
				"context": overrides.Context,
			})
		}
	}

	// Overriding with static configuration
	if !kubeconfig.Inscure.IsNull() {
		overrides.ClusterInfo.InsecureSkipTLSVerify = kubeconfig.Inscure.ValueBool()
	}
	if !kubeconfig.TlsServerName.IsNull() {
		overrides.ClusterInfo.TLSServerName = kubeconfig.TlsServerName.ValueString()
	}
	if !kubeconfig.ClusterCACertificate.IsNull() {
		overrides.ClusterInfo.CertificateAuthorityData = bytes.NewBufferString(kubeconfig.ClusterCACertificate.ValueString()).Bytes()
	}
	if !kubeconfig.ClientCertificate.IsNull() {
		overrides.AuthInfo.ClientCertificateData = bytes.NewBufferString(kubeconfig.ClientCertificate.ValueString()).Bytes()
	}
	if !kubeconfig.Host.IsNull() {
		hasCA := len(overrides.ClusterInfo.CertificateAuthorityData) != 0
		hasCert := len(overrides.AuthInfo.ClientCertificateData) != 0
		defaultTLS := hasCA || hasCert || overrides.ClusterInfo.InsecureSkipTLSVerify
		host, _, err := rest.DefaultServerURL(kubeconfig.Host.ValueString(), "", apimachineryschema.GroupVersion{}, defaultTLS)
		if err != nil {
			return nil, err
		}

		overrides.ClusterInfo.Server = host.String()
	}
	if !kubeconfig.Username.IsNull() {
		overrides.AuthInfo.Username = kubeconfig.Username.ValueString()
	}
	if !kubeconfig.Password.IsNull() {
		overrides.AuthInfo.Password = kubeconfig.Password.ValueString()
	}
	if !kubeconfig.ClientKey.IsNull() {
		overrides.AuthInfo.ClientKeyData = bytes.NewBufferString(kubeconfig.ClientKey.ValueString()).Bytes()
	}
	if !kubeconfig.Token.IsNull() {
		overrides.AuthInfo.Token = kubeconfig.Token.ValueString()
	}
	if !kubeconfig.ProxyURL.IsNull() {
		overrides.ClusterDefaults.ProxyURL = kubeconfig.ProxyURL.ValueString()
	}
	if !kubeconfig.Exec.IsNull() {
		exec := &clientcmdapi.ExecConfig{}

		//if spec, ok := v.([]interface{})[0].(map[string]interface{}); ok {
		//	exec.InteractiveMode = clientcmdapi.IfAvailableExecInteractiveMode
		//	exec.APIVersion = spec["api_version"].(string)
		//	exec.Command = spec["command"].(string)
		//	exec.Args = expandStringSlice(spec["args"].([]interface{}))
		//	for kk, vv := range spec["env"].(map[string]interface{}) {
		//		exec.Env = append(exec.Env, clientcmdapi.ExecEnvVar{Name: kk, Value: vv.(string)})
		//	}
		//} else {
		//	log.Printf("[ERROR] Failed to parse exec")
		//	return nil, fmt.Errorf("failed to parse exec")
		//}

		overrides.AuthInfo.Exec = exec
	}

	overrides.Context.Namespace = "default"
	if namespace != nil {
		overrides.Context.Namespace = *namespace
	}

	client := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, overrides)
	if client == nil {
		tflog.Error(ctx, "failed to initialize kubernetes config")
		return nil, nil
	}
	tflog.Trace(ctx, "successfully initialized kubernetes config")

	return &KubeConfig{ClientConfig: client}, nil
}
