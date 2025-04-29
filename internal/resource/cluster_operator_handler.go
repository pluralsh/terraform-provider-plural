package resource

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"terraform-provider-plural/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/pluralsh/plural-cli/pkg/console"
	"github.com/pluralsh/plural-cli/pkg/helm"
	"github.com/pluralsh/plural-cli/pkg/utils"
	"github.com/pluralsh/polly/algorithms"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func InstallOrUpgradeAgent(ctx context.Context, client *client.Client, kubeconfig *Kubeconfig, repoUrl string,
	values *string, consoleUrl string, token string, d *diag.Diagnostics) error {
	workingDir, chartPath, err := fetchVendoredAgentChart(consoleUrl)
	if err != nil {
		d.AddWarning("Client Warning", fmt.Sprintf("Could not fetch vendored agent chart, using chart from the registry: %s", err))
	}
	if workingDir != "" {
		defer func(path string) {
			if err := os.RemoveAll(path); err != nil {
				d.AddError("Provider Error", fmt.Sprintf("Cannot remove temporary working directory, got error: %s", err))
			}
		}(workingDir)
	}

	handler, err := NewOperatorHandler(ctx, client, kubeconfig, repoUrl, chartPath, values, consoleUrl, token)
	if err != nil {
		return err
	}

	return handler.Apply()
}

func fetchVendoredAgentChart(consoleURL string) (string, string, error) {
	parsedConsoleURL, err := url.Parse(consoleURL)
	if err != nil {
		return "", "", fmt.Errorf("cannot parse console URL: %s", err.Error())
	}

	directory, err := os.MkdirTemp("", "agent-chart-")
	if err != nil {
		return directory, "", fmt.Errorf("cannot create directory: %s", err.Error())
	}

	agentChartURL := fmt.Sprintf("https://%s/ext/v1/agent/chart", parsedConsoleURL.Host)
	agentChartPath := filepath.Join(directory, "agent-chart.tgz")
	if err = utils.DownloadFile(agentChartPath, agentChartURL); err != nil {
		return directory, "", fmt.Errorf("cannot download agent chart: %s", err.Error())
	}

	return directory, agentChartPath, nil
}

func NewOperatorHandler(ctx context.Context, client *client.Client, kubeconfig *Kubeconfig, repoUrl, chartPath string,
	values *string, consoleUrl, token string) (*OperatorHandler, error) {
	settings, err := client.GetDeploymentSettings(ctx)
	if err != nil {
		return nil, err
	}

	k, err := newKubeconfig(ctx, kubeconfig, lo.ToPtr(console.OperatorNamespace))
	if err != nil {
		return nil, err
	}

	clientSet, err := k.ToClientSet()
	if err != nil {
		return nil, err
	}

	additionalValues := map[string]any{}
	if values != nil {
		if err = yaml.Unmarshal([]byte(*values), &additionalValues); err != nil {
			return nil, err
		}
	}

	handler := &OperatorHandler{
		ctx:               ctx,
		consoleURL:        consoleUrl,
		deployToken:       token,
		settings:          settings.DeploymentSettings,
		clientSet:         clientSet,
		vendoredChartPath: chartPath,
		additionalValues:  additionalValues,
	}

	if err := handler.init(k, repoUrl); err != nil {
		return nil, err
	}

	return handler, nil
}

type OperatorHandler struct {
	ctx         context.Context
	consoleURL  string
	deployToken string
	settings    *gqlclient.DeploymentSettingsFragment
	clientSet   *kubernetes.Clientset

	// vendoredChartPath contains a local path to vendored agent chart if it was downloadable, it is empty otherwise.
	vendoredChartPath string

	chart            *chart.Chart
	configuration    *action.Configuration
	additionalValues map[string]any
}

func (oh *OperatorHandler) init(kubeconfig *KubeConfig, repoUrl string) error {
	if oh.configuration != nil {
		return fmt.Errorf("operator handler is already initialized")
	}

	oh.configuration = new(action.Configuration)
	if err := oh.configuration.Init(kubeconfig, console.OperatorNamespace, "", logrus.Debugf); err != nil {
		return err
	}

	var path string
	var err error
	if oh.vendoredChartPath != "" {
		path = oh.vendoredChartPath
	} else {
		if err := helm.AddRepo(console.ReleaseName, repoUrl); err != nil {
			return err
		}

		install := action.NewInstall(oh.configuration)
		if oh.settings != nil {
			install.Version = strings.TrimPrefix(oh.settings.AgentVsn, "v")
		}

		chartName := fmt.Sprintf("%s/%s", console.ReleaseName, console.ChartName)
		if path, err = install.LocateChart(chartName, cli.New()); err != nil {
			return err
		}
	}

	oh.chart, err = loader.Load(path)
	return err
}

func (oh *OperatorHandler) Apply() error {
	if err := oh.ensureNamespace(); err != nil {
		return err
	}

	isChartInstalled, err := oh.isChartInstalled()
	if err != nil {
		return err
	}

	if isChartInstalled {
		return oh.upgrade()
	}

	return oh.install()
}

func (oh *OperatorHandler) ensureNamespace() error {
	if _, err := oh.clientSet.CoreV1().Namespaces().Get(oh.ctx, console.OperatorNamespace, metav1.GetOptions{}); err == nil {
		return nil
	}

	_, err := oh.clientSet.CoreV1().Namespaces().Create(oh.ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: console.OperatorNamespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "plural",
				"app.plural.sh/name":           console.OperatorNamespace,
			},
		},
	}, metav1.CreateOptions{})
	return err
}

// isChartInstalled checks whether a chart is already installed in a namespace based on the provided chart spec.
// Note that this function only considers the contained chart name and namespace.
func (oh *OperatorHandler) isChartInstalled() (bool, error) {
	releases, err := oh.listReleases(action.ListAll)
	if err != nil {
		return false, err
	}

	for _, r := range releases {
		if r.Name == console.ReleaseName && r.Namespace == console.OperatorNamespace {
			return true, nil
		}
	}

	return false, nil
}

// listReleases lists all releases that match the given state.
func (oh *OperatorHandler) listReleases(state action.ListStates) ([]*release.Release, error) {
	list := action.NewList(oh.configuration)
	list.StateMask = state
	return list.Run()
}

func (oh *OperatorHandler) upgrade() error {
	upgrade := action.NewUpgrade(oh.configuration)
	upgrade.Namespace = console.OperatorNamespace
	upgrade.Timeout = 5 * time.Minute
	upgrade.Wait = false

	values, err := oh.values()
	if err != nil {
		return err
	}

	_, err = upgrade.Run(console.ReleaseName, oh.chart, values)
	return err
}

func (oh *OperatorHandler) install() error {
	install := action.NewInstall(oh.configuration)
	install.Namespace = console.OperatorNamespace
	install.ReleaseName = console.ReleaseName
	install.Timeout = 5 * time.Minute
	install.Wait = false
	install.CreateNamespace = true

	values, err := oh.values()
	if err != nil {
		return err
	}

	_, err = install.Run(oh.chart, values)
	return err
}

func (oh *OperatorHandler) values() (map[string]any, error) {
	settingsValues := map[string]any{}
	if oh.settings != nil && oh.settings.AgentHelmValues != nil {
		if err := yaml.Unmarshal([]byte(*oh.settings.AgentHelmValues), &settingsValues); err != nil {
			return nil, err
		}
	}

	return algorithms.Merge(map[string]any{
		"secrets":    map[string]string{"deployToken": oh.deployToken},
		"consoleUrl": fmt.Sprintf("%s/ext/gql", oh.consoleURL),
	}, oh.additionalValues, settingsValues), nil
}
