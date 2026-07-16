package resource

import (
	"fmt"
	"strings"

	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/pluralsh/polly/template"
	"github.com/samber/lo"
	"sigs.k8s.io/yaml"
)

func resolveAgentHelmValues(
	settings *gqlclient.DeploymentSettingsFragment,
	cluster *gqlclient.ClusterFragment,
) (map[string]any, error) {
	if settings == nil || settings.AgentHelmValues == nil {
		return map[string]any{}, nil
	}

	raw := *settings.AgentHelmValues
	rendered := raw
	if lo.FromPtr(settings.AgentHelmValuesTemplateable) {
		if cluster == nil {
			return nil, fmt.Errorf("cluster context is required to render agent helm values")
		}
		bindings := agentHelmValuesBindings(cluster)
		out, err := template.RenderLiquid([]byte(raw), bindings)
		if err != nil {
			return nil, fmt.Errorf("rendering agent helm values: %w", err)
		}
		rendered = string(out)
	}

	globalVals := map[string]any{}
	if err := yaml.Unmarshal([]byte(rendered), &globalVals); err != nil {
		return nil, err
	}

	return globalVals, nil
}

func agentHelmValuesBindings(cluster *gqlclient.ClusterFragment) map[string]any {
	return map[string]any{
		"cluster": clusterBindings(cluster),
	}
}

func clusterBindings(cluster *gqlclient.ClusterFragment) map[string]any {
	if cluster == nil {
		return map[string]any{}
	}

	res := map[string]any{
		"ID":             cluster.ID,
		"Self":           cluster.Self,
		"Handle":         cluster.Handle,
		"Name":           cluster.Name,
		"Version":        cluster.Version,
		"CurrentVersion": cluster.CurrentVersion,
		"KasUrl":         cluster.KasURL,
		"Metadata":       cluster.Metadata,
		"Distro":         cluster.Distro,
		"Tags":           clusterTagsMap(cluster.Tags),
	}

	lowercase := make(map[string]any, len(res))
	for k, v := range res {
		lowercase[strings.ToLower(k)] = v
	}
	for k, v := range lowercase {
		res[k] = v
	}
	res["kasUrl"] = cluster.KasURL
	res["currentVersion"] = cluster.CurrentVersion

	return res
}

func clusterTagsMap(tags []*gqlclient.ClusterTags) map[string]string {
	res := map[string]string{}
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		res[tag.Name] = tag.Value
	}
	return res
}
