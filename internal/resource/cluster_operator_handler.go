package resource

import (
	"context"
	"fmt"
	"time"

	"github.com/pluralsh/plural-cli/pkg/console"
	"github.com/pluralsh/plural-cli/pkg/helm"
	"github.com/pluralsh/polly/algorithms"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

type OperatorHandler struct {
	ctx context.Context
	// kubeconfig is a model.Kubeconfig data model read from terraform
	kubeconfig *Kubeconfig
	// url is an url to the Console API, i.e. https://console.mycluster.onplural.sh
	url string

	// additional values used on install
	vals map[string]interface{}

	// Preconfigured helm actions and chart
	chart         *chart.Chart
	configuration *action.Configuration
	install       *action.Install
	upgrade       *action.Upgrade
	uninstall     *action.Uninstall
}

func (oh *OperatorHandler) init() error {
	oh.configuration = new(action.Configuration)

	kubeconfig, err := newKubeconfig(oh.ctx, oh.kubeconfig, lo.ToPtr(console.OperatorNamespace))
	if err != nil {
		return err
	}

	err = oh.configuration.Init(kubeconfig, console.OperatorNamespace, "", logrus.Debugf)
	if err != nil {
		return err
	}

	err = oh.initRepo()
	if err != nil {
		return err
	}

	err = oh.initChart()
	if err != nil {
		return err
	}

	oh.initInstallAction()
	oh.initUpgradeAction()
	oh.initUninstallAction()

	return nil
}

func (oh *OperatorHandler) initRepo() error {
	return helm.AddRepo(console.ReleaseName, console.RepoUrl)
}

func (oh *OperatorHandler) initChart() error {
	client := action.NewInstall(oh.configuration)
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
	oh.install.Wait = true
	oh.install.CreateNamespace = true
}

func (oh *OperatorHandler) initUpgradeAction() {
	oh.upgrade = action.NewUpgrade(oh.configuration)

	oh.upgrade.Namespace = console.OperatorNamespace
	oh.upgrade.Timeout = 5 * time.Minute
	oh.upgrade.Wait = true
}

func (oh *OperatorHandler) initUninstallAction() {
	oh.uninstall = action.NewUninstall(oh.configuration)

	oh.uninstall.Timeout = 5 * time.Minute
	oh.uninstall.Wait = true
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
	client := action.NewList(oh.configuration)
	client.StateMask = state

	return client.Run()
}

func (oh *OperatorHandler) values(token string) map[string]interface{} {
	vals := map[string]interface{}{
		"secrets": map[string]string{
			"deployToken": token,
		},
		"consoleUrl": fmt.Sprintf("%s/ext/gql", oh.url),
	}
	return algorithms.Merge(vals, oh.vals)
}

func (oh *OperatorHandler) InstallOrUpgrade(token string) error {
	exists, err := oh.chartExists()
	if err != nil {
		return err
	}

	if exists {
		return oh.Upgrade(token)
	}

	return oh.Install(token)
}

func (oh *OperatorHandler) Install(token string) error {
	_, err := oh.install.Run(oh.chart, oh.values(token))
	return err
}

func (oh *OperatorHandler) Upgrade(token string) error {
	_, err := oh.upgrade.Run(console.ReleaseName, oh.chart, oh.values(token))
	return err
}

func (oh *OperatorHandler) Uninstall() error {
	_, err := oh.uninstall.Run(console.ReleaseName)
	return err
}

func NewOperatorHandler(ctx context.Context, kubeconfig *Kubeconfig, values *string, consoleUrl string) (*OperatorHandler, error) {
	vals := map[string]interface{}{}
	if values != nil {
		if err := yaml.Unmarshal([]byte(*values), &vals); err != nil {
			return nil, err
		}
	}

	handler := &OperatorHandler{
		ctx:        ctx,
		kubeconfig: kubeconfig,
		url:        consoleUrl,
		vals:       vals,
	}

	err := handler.init()
	if err != nil {
		return nil, err
	}

	return handler, nil
}
