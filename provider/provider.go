package provider

import "github.com/valyala/fastjson"

// Provider combines task and secret retrieval for cloud providers.
type Provider interface {
	GetTaskDefinitions(taskDefinitionName string) ([]*ContainerDefinition, error)
	GetSecretValue(secretArn, secretName string) (*fastjson.Value, error)
}

// ContainerDefinition represents the container details.
type ContainerDefinition struct {
	Name        string
	Environment map[string]string
	Secrets     map[string]string
}

var providers = make(map[string]Provider)

// RegisterProvider registers a provider.
func RegisterProvider(name string, provider Provider) {
	providers[name] = provider
}

// GetProvider gets a specific provider.
func GetProvider(name string) (Provider, bool) {
	provider, exists := providers[name]
	return provider, exists
}

func Providers() []string {
	var names []string
	for name := range providers {
		names = append(names, name)
	}
	return names
}
