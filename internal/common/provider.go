package common

import console "terraform-provider-plural/internal/client"

type ProviderData struct {
	Client     *console.Client
	ConsoleUrl string
	KubeClient *KubeClient
}

func NewProviderData(client *console.Client, consoleUrl string, kubeClient *KubeClient) *ProviderData {
	return &ProviderData{
		Client:     client,
		ConsoleUrl: consoleUrl,
		KubeClient: kubeClient,
	}
}
