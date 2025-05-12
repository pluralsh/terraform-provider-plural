package common

import internalclient "terraform-provider-plural/internal/client"

type ProviderData struct {
	Client     *internalclient.Client
	ConsoleUrl string
	KubeClient *KubeClient
}

func NewProviderData(client *internalclient.Client, consoleUrl string, kubeClient *KubeClient) *ProviderData {
	return &ProviderData{
		Client:     client,
		ConsoleUrl: consoleUrl,
		KubeClient: kubeClient,
	}
}
