package resource

import (
	"fmt"
	"strings"
	"testing"

	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
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

func TestOperatorHandlerValuesRenderTemplatedAgentHelmValues(t *testing.T) {
	settingsValues := `
clusterName: '{{ cluster.name }}'
region: '{{ cluster.metadata.region }}'
environment: '{{ cluster.tags.env }}'
clusterID: '{{ cluster.id }}'
`
	handler := &OperatorHandler{
		consoleURL:  "https://console.example.com",
		deployToken: "token",
		settings: &gqlclient.DeploymentSettingsFragment{
			AgentHelmValues:             &settingsValues,
			AgentHelmValuesTemplateable: lo.ToPtr(true),
		},
		cluster: &gqlclient.ClusterFragment{
			ID:       "cluster-id",
			Name:     "dev-cluster",
			Metadata: map[string]any{"region": "us-east-1"},
			Tags:     []*gqlclient.ClusterTags{{Name: "env", Value: "staging"}},
		},
		clusterId:        "cluster-id",
		additionalValues: map[string]any{},
	}

	values, err := handler.values()
	if err != nil {
		t.Fatalf("values returned error: %v", err)
	}

	if got := values["clusterName"]; got != "dev-cluster" {
		t.Fatalf("expected clusterName to be rendered from cluster context, got %v", got)
	}
	if got := values["region"]; got != "us-east-1" {
		t.Fatalf("expected region to be rendered from cluster metadata, got %v", got)
	}
	if got := values["environment"]; got != "staging" {
		t.Fatalf("expected environment to be rendered from cluster tags, got %v", got)
	}
	if got := values["clusterID"]; got != "cluster-id" {
		t.Fatalf("expected clusterID to be rendered from cluster id, got %v", got)
	}
}

func TestOperatorHandlerValuesRequireClusterWhenTemplating(t *testing.T) {
	settingsValues := `
clusterName: '{{ cluster.name }}'
`
	handler := &OperatorHandler{
		consoleURL:  "https://console.example.com",
		deployToken: "token",
		settings: &gqlclient.DeploymentSettingsFragment{
			AgentHelmValues:             &settingsValues,
			AgentHelmValuesTemplateable: lo.ToPtr(true),
		},
		clusterId:        "cluster-id",
		additionalValues: map[string]any{},
	}

	_, err := handler.values()
	if err == nil {
		t.Fatal("expected error when cluster context is missing")
	}
	if got := err.Error(); !strings.Contains(got, "cluster context is required") {
		t.Fatalf("expected cluster required error, got %v", err)
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
