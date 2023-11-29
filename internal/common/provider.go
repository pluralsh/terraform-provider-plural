package common

import internalclient "terraform-provider-plural/internal/client"

type ProviderData struct {
	Client     *internalclient.Client
	ConsoleUrl string
}

func NewProviderData(client *internalclient.Client, consoleUrl string) *ProviderData {
	return &ProviderData{
		Client:     client,
		ConsoleUrl: consoleUrl,
	}
}
