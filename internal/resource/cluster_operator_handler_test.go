package resource

import (
	"fmt"
	"testing"

	gqlclient "github.com/pluralsh/console/go/client"
	"sigs.k8s.io/yaml"
)

func TestOperatorHandlerValuesResourceHelmValuesOverrideDeploymentSettings(t *testing.T) {
	settingsValues := `
image:
  repository: settings-repository
  tag: settings-tag
replicaCount: 1
podAnnotations:
  settings: "true"
`
	resourceValues := `
image:
  tag: resource-tag
replicaCount: 3
podAnnotations:
  resource: "true"
`

	additionalValues := map[string]any{}
	if err := yaml.Unmarshal([]byte(resourceValues), &additionalValues); err != nil {
		t.Fatalf("failed to unmarshal resource values: %v", err)
	}

	handler := &OperatorHandler{
		consoleURL:       "https://console.example.com",
		deployToken:      "token",
		settings:         &gqlclient.DeploymentSettingsFragment{AgentHelmValues: &settingsValues},
		additionalValues: additionalValues,
		clusterId:        "cluster-id",
	}

	values, err := handler.values()
	if err != nil {
		t.Fatalf("values returned error: %v", err)
	}

	image := nestedMap(t, values, "image")
	if got := image["tag"]; got != "resource-tag" {
		t.Fatalf("expected resource image tag to win, got %v", got)
	}
	if got := image["repository"]; got != "settings-repository" {
		t.Fatalf("expected settings image repository to be preserved, got %v", got)
	}

	if got := fmt.Sprint(values["replicaCount"]); got != "3" {
		t.Fatalf("expected resource replicaCount to win, got %v", got)
	}

	podAnnotations := nestedMap(t, values, "podAnnotations")
	if got := podAnnotations["settings"]; got != "true" {
		t.Fatalf("expected settings annotation to be preserved, got %v", got)
	}
	if got := podAnnotations["resource"]; got != "true" {
		t.Fatalf("expected resource annotation to be preserved, got %v", got)
	}
}

func nestedMap(t *testing.T, values map[string]any, key string) map[string]any {
	t.Helper()

	value, ok := values[key].(map[string]any)
	if !ok {
		t.Fatalf("expected %q to be a nested map, got %T", key, values[key])
	}

	return value
}
