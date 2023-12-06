package resource

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pluralsh/polly/algorithms"
	"github.com/samber/lo"
	"k8s.io/client-go/discovery/cached/disk"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mitchellh/go-homedir"
	"k8s.io/apimachinery/pkg/api/meta"
	apimachineryschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type KubeConfig struct {
	ClientConfig clientcmd.ClientConfig
}

func (k *KubeConfig) ToRESTConfig() (*rest.Config, error) {
	return k.ToRawKubeConfigLoader().ClientConfig()
}

func (k *KubeConfig) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	config, err := k.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	return disk.NewCachedDiscoveryClientForConfig(config, os.TempDir(), os.TempDir(), 1*time.Minute)
}

func (k *KubeConfig) ToRESTMapper() (meta.RESTMapper, error) {
	client, err := k.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	return restmapper.NewShortcutExpander(restmapper.NewDeferredDiscoveryRESTMapper(client), client), nil
}

func (k *KubeConfig) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return k.ClientConfig
}

func newKubeconfig(ctx context.Context, kubeconfig *Kubeconfig, namespace *string) (*KubeConfig, error) {
	overrides := &clientcmd.ConfigOverrides{}
	loader := &clientcmd.ClientConfigLoadingRules{}

	if !lo.IsEmpty(kubeconfig.ConfigPath.ValueString()) {
		tflog.Info(ctx, "using kubeconfig", map[string]interface{}{
			"kubeconfig": kubeconfig.ConfigPath.ValueString(),
		})

		path, err := homedir.Expand(kubeconfig.ConfigPath.ValueString())
		if err != nil {
			return nil, err
		}
		loader.ExplicitPath = path

		if !lo.IsEmpty(kubeconfig.ConfigContext.ValueString()) || !lo.IsEmpty(kubeconfig.ConfigContextAuthInfo.ValueString()) || !lo.IsEmpty(kubeconfig.ConfigContextCluster.ValueString()) {
			if !lo.IsEmpty(kubeconfig.ConfigContext.ValueString()) {
				overrides.CurrentContext = kubeconfig.ConfigContext.ValueString()
				tflog.Info(ctx, "using custom current context", map[string]interface{}{
					"context": overrides.CurrentContext,
				})
			}

			overrides.Context = clientcmdapi.Context{}
			if !lo.IsEmpty(kubeconfig.ConfigContextAuthInfo.ValueString()) {
				overrides.Context.AuthInfo = kubeconfig.ConfigContextAuthInfo.ValueString()
			}
			if !lo.IsEmpty(kubeconfig.ConfigContextCluster.ValueString()) {
				overrides.Context.Cluster = kubeconfig.ConfigContextCluster.ValueString()
			}
			tflog.Info(ctx, "using overridden context", map[string]interface{}{
				"context": overrides.Context,
			})
		}
	}

	// Overriding with static configuration
	if !kubeconfig.Insecure.ValueBool() {
		overrides.ClusterInfo.InsecureSkipTLSVerify = kubeconfig.Insecure.ValueBool()
	}
	if !lo.IsEmpty(kubeconfig.TlsServerName.ValueString()) {
		overrides.ClusterInfo.TLSServerName = kubeconfig.TlsServerName.ValueString()
	}
	if !lo.IsEmpty(kubeconfig.ClusterCACertificate.ValueString()) {
		overrides.ClusterInfo.CertificateAuthorityData = bytes.NewBufferString(kubeconfig.ClusterCACertificate.ValueString()).Bytes()
	}
	if !lo.IsEmpty(kubeconfig.ClientCertificate.ValueString()) {
		overrides.AuthInfo.ClientCertificateData = bytes.NewBufferString(kubeconfig.ClientCertificate.ValueString()).Bytes()
	}
	if !lo.IsEmpty(kubeconfig.Host.ValueString()) {
		hasCA := len(overrides.ClusterInfo.CertificateAuthorityData) != 0
		hasCert := len(overrides.AuthInfo.ClientCertificateData) != 0
		defaultTLS := hasCA || hasCert || overrides.ClusterInfo.InsecureSkipTLSVerify
		host, _, err := rest.DefaultServerURL(kubeconfig.Host.ValueString(), "", apimachineryschema.GroupVersion{}, defaultTLS)
		if err != nil {
			return nil, err
		}

		overrides.ClusterInfo.Server = host.String()
	}
	if !lo.IsEmpty(kubeconfig.Username.ValueString()) {
		overrides.AuthInfo.Username = kubeconfig.Username.ValueString()
	}
	if !lo.IsEmpty(kubeconfig.Password.ValueString()) {
		overrides.AuthInfo.Password = kubeconfig.Password.ValueString()
	}
	if !lo.IsEmpty(kubeconfig.ClientKey.ValueString()) {
		overrides.AuthInfo.ClientKeyData = bytes.NewBufferString(kubeconfig.ClientKey.ValueString()).Bytes()
	}
	if !lo.IsEmpty(kubeconfig.Token.ValueString()) {
		overrides.AuthInfo.Token = kubeconfig.Token.ValueString()
	}
	if !lo.IsEmpty(kubeconfig.ProxyURL.ValueString()) {
		overrides.ClusterDefaults.ProxyURL = kubeconfig.ProxyURL.ValueString()
	}
	if kubeconfig.Exec != nil {
		exec := &clientcmdapi.ExecConfig{
			InteractiveMode: clientcmdapi.IfAvailableExecInteractiveMode,
			APIVersion:      kubeconfig.Exec.APIVersion.ValueString(),
			Command:         kubeconfig.Exec.Command.ValueString(),
		}

		if !kubeconfig.Exec.Env.IsNull() {
			envElements := make(map[string]types.String)
			diags := kubeconfig.Exec.Env.ElementsAs(ctx, &envElements, false)
			if diags.HasError() {
				return nil, fmt.Errorf("error while parsing kubeconfig exec env, got diagnostics: %+v", diags)
			}

			env := make([]clientcmdapi.ExecEnvVar, 0)
			for k, v := range envElements {
				env = append(env, clientcmdapi.ExecEnvVar{
					Name:  k,
					Value: v.ValueString(),
				})
			}
			exec.Env = env
		}

		if !kubeconfig.Exec.Args.IsNull() {
			argsElements := make([]types.String, len(kubeconfig.Exec.Args.Elements()))
			diags := kubeconfig.Exec.Args.ElementsAs(ctx, &argsElements, false)
			if diags.HasError() {
				return nil, fmt.Errorf("error while parsing kubeconfig exec args, got diagnostics: %+v", diags)
			}

			exec.Args = algorithms.Map(argsElements, func(v types.String) string { return v.ValueString() })
		}

		overrides.AuthInfo.Exec = exec
	}

	overrides.Context.Namespace = "default"
	if namespace != nil {
		overrides.Context.Namespace = *namespace
	}

	client := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, overrides)
	if client == nil {
		err := fmt.Errorf("failed to initialize kubernetes config")
		tflog.Error(ctx, err.Error())
		return nil, err
	}

	tflog.Trace(ctx, "successfully initialized kubernetes config")
	return &KubeConfig{ClientConfig: client}, nil
}
