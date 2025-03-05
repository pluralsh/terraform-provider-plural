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

type OperatorHandler struct {
	client *client.Client

	kube *kubernetes.Clientset

	ctx context.Context

	// kubeconfig is a model.Kubeconfig data model read from terraform
	kubeconfig *Kubeconfig

	settings *gqlclient.DeploymentSettingsFragment

	// url is an url to the Console API, i.e. https://console.mycluster.onplural.sh
	url string

	// repoUrl is an URL of the deployment agent chart.
	repoUrl string

	// agentChartPath contains a local path to vendored agent chart if it was downloadable, it is empty otherwise.
	agentChartPath string

	// additional values used on install
	vals map[string]any

	// Preconfigured helm actions and chart
	chart         *chart.Chart
	configuration *action.Configuration
	install       *action.Install
	upgrade       *action.Upgrade
}

func (oh *OperatorHandler) init() error {
	oh.configuration = new(action.Configuration)

	kubeconfig, err := newKubeconfig(oh.ctx, oh.kubeconfig, lo.ToPtr(console.OperatorNamespace))
	if err != nil {
		return err
	}
	kube, err := kubeconfig.ToClientSet()
	if err != nil {
		return err
	}
	oh.kube = kube

	err = oh.configuration.Init(kubeconfig, console.OperatorNamespace, "", logrus.Debugf)
	if err != nil {
		return err
	}

	oh.initSettings()

	repo := lo.Ternary(oh.agentChartPath != "", oh.agentChartPath, oh.repoUrl)
	if err = helm.AddRepo(console.ReleaseName, repo); err != nil {
		return err
	}

	err = oh.initChart()
	if err != nil {
		return err
	}

	oh.initInstallAction()
	oh.initUpgradeAction()

	return nil
}

func (oh *OperatorHandler) initSettings() {
	settings, err := oh.client.GetDeploymentSettings(oh.ctx)
	if err != nil {
		return
	}
	oh.settings = settings.DeploymentSettings
}

func (oh *OperatorHandler) initChart() error {
	vsn := ""
	if oh.settings != nil {
		vsn = oh.settings.AgentVsn
	}

	client := action.NewInstall(oh.configuration)
	client.ChartPathOptions.Version = strings.TrimPrefix(vsn, "v") // TODO ?
	locateName := fmt.Sprintf("%s/%s", console.ReleaseName, console.ChartName)
	path, err := client.ChartPathOptions.LocateChart(locateName, cli.New())
	if err != nil {
		return err
	}

	oh.chart, err = loader.Load(path)
	return err
}

func (oh *OperatorHandler) initInstallAction() {
	oh.install = action.NewInstall(oh.configuration)
	oh.install.Namespace = console.OperatorNamespace
	oh.install.ReleaseName = console.ReleaseName
	oh.install.Timeout = 5 * time.Minute
	oh.install.Wait = false
	oh.install.CreateNamespace = true
}

func (oh *OperatorHandler) initUpgradeAction() {
	oh.upgrade = action.NewUpgrade(oh.configuration)
	oh.upgrade.Namespace = console.OperatorNamespace
	oh.upgrade.Timeout = 5 * time.Minute
	oh.upgrade.Wait = false
}

// chartExists checks whether a chart is already installed
// in a namespace or not based on the provided chart spec.
// Note that this function only considers the contained chart name and namespace.
func (oh *OperatorHandler) chartExists() (bool, error) {
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
	c := action.NewList(oh.configuration)
	c.StateMask = state
	return c.Run()
}

func (oh *OperatorHandler) values(token string) (map[string]any, error) {
	globalVals := map[string]any{}
	vals := map[string]any{
		"secrets": map[string]string{
			"deployToken": token,
		},
		"consoleUrl": fmt.Sprintf("%s/ext/gql", oh.url),
	}

	if oh.settings != nil && oh.settings.AgentHelmValues != nil {
		if err := yaml.Unmarshal([]byte(*oh.settings.AgentHelmValues), &globalVals); err != nil {
			return nil, err
		}
	}
	return algorithms.Merge(vals, oh.vals, globalVals), nil
}

func (oh *OperatorHandler) UpsertNamespace() error {
	_, err := oh.kube.CoreV1().Namespaces().Get(oh.ctx, console.OperatorNamespace, metav1.GetOptions{})
	if err == nil {
		return nil
	}

	_, err = oh.kube.CoreV1().Namespaces().Create(oh.ctx, &v1.Namespace{
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

func (oh *OperatorHandler) Install(token string) error {
	values, err := oh.values(token)
	if err != nil {
		return err
	}
	_, err = oh.install.Run(oh.chart, values)
	return err
}

func (oh *OperatorHandler) Upgrade(token string) error {
	values, err := oh.values(token)
	if err != nil {
		return err
	}
	_, err = oh.upgrade.Run(console.ReleaseName, oh.chart, values)
	return err
}

func InstallOrUpgradeAgent(ctx context.Context, client *client.Client, kubeconfig *Kubeconfig, repoUrl string, values *string, consoleUrl string, token string, d diag.Diagnostics) error {
	workingDir, agentChartPath, err := fetchVendoredAgentChart(consoleUrl)
	if err != nil {
		d.AddWarning("Client Warning", fmt.Sprintf("Could not fetch vendored agent chart, using chart from the registry: %s", err))
	}
	if workingDir != "" {
		defer func(path string) {
			if err := os.RemoveAll(path); err != nil {
				d.AddError("Provider Error", fmt.Sprintf("Cannot remove working directory, got error: %s", err))
			}
		}(workingDir)
	}

	vals := map[string]any{}
	if values != nil {
		if err := yaml.Unmarshal([]byte(*values), &vals); err != nil {
			return err
		}
	}

	handler := &OperatorHandler{
		client:         client,
		ctx:            ctx,
		kubeconfig:     kubeconfig,
		repoUrl:        repoUrl,
		agentChartPath: agentChartPath,
		url:            consoleUrl,
		vals:           vals,
	}

	if err := handler.init(); err != nil {
		return err
	}

	if err = handler.UpsertNamespace(); err != nil {
		return err
	}

	exists, err := handler.chartExists()
	if err != nil {
		return err
	}
	if exists {
		return handler.Upgrade(token)
	}
	return handler.Install(token)
}

// fetchVendoredAgentChart downloads vendored agent chart.
// It returns the working directory name, chart path and any error encountered.
func fetchVendoredAgentChart(consoleURL string) (string, string, error) {
	parsedConsoleURL, err := url.Parse(consoleURL)
	if err != nil {
		return "", "", fmt.Errorf("cannot parse console URL: %s", err.Error())
	}

	directory, err := os.MkdirTemp("", "agent-chart-*")
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
