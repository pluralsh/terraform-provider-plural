package resource

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mitchellh/go-homedir"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"

	"terraform-provider-plural/internal/defaults"
	"terraform-provider-plural/internal/model"
)

func kubeconfigAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:            true,
				Default:             defaults.Env("PLURAL_KUBE_HOST", ""),
				MarkdownDescription: "The address of the Kubernetes clusters. Can be sourced from `PLURAL_KUBE_HOST`.",
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

	Burst int

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

	// The more groups you have, the more discovery requests you need to make.
	// given 25 groups (our groups + a few custom resources) with one-ish version each, discovery needs to make 50 requests
	// double it just so we don't end up here again for a while.  This config is only used for discovery.
	config.Burst = k.Burst

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

func newKubeConfig(configData *model.Kubeconfig, namespace *string) (*KubeConfig, error) {
	overrides := &clientcmd.ConfigOverrides{}
	loader := &clientcmd.ClientConfigLoadingRules{}

	if v := os.Getenv("PLURAL_KUBE_CONFIG_PATH"); v != "" {
		configData.ConfigPath = types.StringValue(v)
	}

	if len(configPaths) > 0 {
		expandedPaths := []string{}
		for _, p := range configPaths {
			path, err := homedir.Expand(p)
			if err != nil {
				return nil, err
			}

			log.Printf("[DEBUG] Using kubeconfig: %s", path)
			expandedPaths = append(expandedPaths, path)
		}

		if len(expandedPaths) == 1 {
			loader.ExplicitPath = expandedPaths[0]
		} else {
			loader.Precedence = expandedPaths
		}

		ctx, ctxOk := k8sGetOk(configData, "config_context")
		authInfo, authInfoOk := k8sGetOk(configData, "config_context_auth_info")
		cluster, clusterOk := k8sGetOk(configData, "config_context_cluster")
		if ctxOk || authInfoOk || clusterOk {
			if ctxOk {
				overrides.CurrentContext = ctx.(string)
				log.Printf("[DEBUG] Using custom current context: %q", overrides.CurrentContext)
			}

			overrides.Context = clientcmdapi.Context{}
			if authInfoOk {
				overrides.Context.AuthInfo = authInfo.(string)
			}
			if clusterOk {
				overrides.Context.Cluster = cluster.(string)
			}
			log.Printf("[DEBUG] Using overidden context: %#v", overrides.Context)
		}
	}

	// Overriding with static configuration
	if v, ok := k8sGetOk(configData, "insecure"); ok {
		overrides.ClusterInfo.InsecureSkipTLSVerify = v.(bool)
	}
	if v, ok := k8sGetOk(configData, "tls_server_name"); ok {
		overrides.ClusterInfo.TLSServerName = v.(string)
	}
	if v, ok := k8sGetOk(configData, "cluster_ca_certificate"); ok {
		overrides.ClusterInfo.CertificateAuthorityData = bytes.NewBufferString(v.(string)).Bytes()
	}
	if v, ok := k8sGetOk(configData, "client_certificate"); ok {
		overrides.AuthInfo.ClientCertificateData = bytes.NewBufferString(v.(string)).Bytes()
	}
	if v, ok := k8sGetOk(configData, "host"); ok {
		// Server has to be the complete address of the kubernetes cluster (scheme://hostname:port), not just the hostname,
		// because `overrides` are processed too late to be taken into account by `defaultServerUrlFor()`.
		// This basically replicates what defaultServerUrlFor() does with config but for overrides,
		// see https://github.com/kubernetes/client-go/blob/v12.0.0/rest/url_utils.go#L85-L87
		hasCA := len(overrides.ClusterInfo.CertificateAuthorityData) != 0
		hasCert := len(overrides.AuthInfo.ClientCertificateData) != 0
		defaultTLS := hasCA || hasCert || overrides.ClusterInfo.InsecureSkipTLSVerify
		host, _, err := rest.DefaultServerURL(v.(string), "", apimachineryschema.GroupVersion{}, defaultTLS)
		if err != nil {
			return nil, err
		}

		overrides.ClusterInfo.Server = host.String()
	}
	if v, ok := k8sGetOk(configData, "username"); ok {
		overrides.AuthInfo.Username = v.(string)
	}
	if v, ok := k8sGetOk(configData, "password"); ok {
		overrides.AuthInfo.Password = v.(string)
	}
	if v, ok := k8sGetOk(configData, "client_key"); ok {
		overrides.AuthInfo.ClientKeyData = bytes.NewBufferString(v.(string)).Bytes()
	}
	if v, ok := k8sGetOk(configData, "token"); ok {
		overrides.AuthInfo.Token = v.(string)
	}

	if v, ok := k8sGetOk(configData, "proxy_url"); ok {
		overrides.ClusterDefaults.ProxyURL = v.(string)
	}

	if v, ok := k8sGetOk(configData, "exec"); ok {
		exec := &clientcmdapi.ExecConfig{}
		if spec, ok := v.([]interface{})[0].(map[string]interface{}); ok {
			exec.InteractiveMode = clientcmdapi.IfAvailableExecInteractiveMode
			exec.APIVersion = spec["api_version"].(string)
			exec.Command = spec["command"].(string)
			exec.Args = expandStringSlice(spec["args"].([]interface{}))
			for kk, vv := range spec["env"].(map[string]interface{}) {
				exec.Env = append(exec.Env, clientcmdapi.ExecEnvVar{Name: kk, Value: vv.(string)})
			}
		} else {
			log.Printf("[ERROR] Failed to parse exec")
			return nil, fmt.Errorf("failed to parse exec")
		}
		overrides.AuthInfo.Exec = exec
	}

	overrides.Context.Namespace = "default"

	if namespace != nil {
		overrides.Context.Namespace = *namespace
	}
	burstLimit := configData.Get("burst_limit").(int)

	client := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, overrides)
	if client == nil {
		log.Printf("[ERROR] Failed to initialize kubernetes config")
		return nil, nil
	}
	log.Printf("[INFO] Successfully initialized kubernetes config")

	return &KubeConfig{ClientConfig: client, Burst: burstLimit}, nil
}
