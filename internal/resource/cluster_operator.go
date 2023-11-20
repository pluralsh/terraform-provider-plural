package resource

import (
	"context"
	"fmt"
	"time"

	"github.com/pluralsh/plural-cli/pkg/helm"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"

	"terraform-provider-plural/internal/model"
)

const (
	// TODO: read from CLI pkg?
	operatorNamespace = "plrl-deploy-operator"
	releaseName       = "deploy-operator"
	chartName         = "deployment-operator"
	repoUrl           = "https://pluralsh.github.io/deployment-operator"
)

type OperatorHandler struct {
	ctx context.Context
	// kubeconfig is a model.Kubeconfig data model read from terraform
	kubeconfig *model.Kubeconfig
	// url is an url to the Console API, i.e. https://console.mycluster.onplural.sh
	url string

	// Preconfigured helm actions and chart
	chart         *chart.Chart
	configuration *action.Configuration
	install       *action.Install
	upgrade       *action.Upgrade
	uninstall     *action.Uninstall
}

func (this *OperatorHandler) init() error {
	this.configuration = new(action.Configuration)

	kubeconfig, err := newKubeconfig(this.ctx, this.kubeconfig, lo.ToPtr(operatorNamespace))
	if err != nil {
		return err
	}

	err = this.configuration.Init(kubeconfig, operatorNamespace, "", logrus.Debugf)
	if err != nil {
		return err
	}

	err = this.initRepo()
	if err != nil {
		return err
	}

	err = this.initChart()
	if err != nil {
		return err
	}

	this.initInstallAction()
	this.initUpgradeAction()
	this.initUninstallAction()

	return nil
}

func (this *OperatorHandler) initRepo() error {
	return helm.AddRepo(releaseName, repoUrl)
}

func (this *OperatorHandler) initChart() error {
	client := action.NewInstall(this.configuration)
	locateName := fmt.Sprintf("%s/%s", releaseName, chartName)
	path, err := client.ChartPathOptions.LocateChart(locateName, cli.New())
	if err != nil {
		return err
	}

	this.chart, err = loader.Load(path)
	return err
}

func (this *OperatorHandler) initInstallAction() {
	this.install = action.NewInstall(this.configuration)

	this.install.Namespace = operatorNamespace
	this.install.ReleaseName = releaseName
	this.install.Timeout = 5 * time.Minute
	this.install.Wait = true
	this.install.CreateNamespace = true
}

func (this *OperatorHandler) initUpgradeAction() {
	this.upgrade = action.NewUpgrade(this.configuration)

	this.upgrade.Namespace = operatorNamespace
	this.upgrade.Timeout = 5 * time.Minute
	this.upgrade.Wait = true
}

func (this *OperatorHandler) initUninstallAction() {
	this.uninstall = action.NewUninstall(this.configuration)

	this.uninstall.Timeout = 5 * time.Minute
	this.uninstall.Wait = true
}

// chartExists checks whether a chart is already installed
// in a namespace or not based on the provided chart spec.
// Note that this function only considers the contained chart name and namespace.
func (this *OperatorHandler) chartExists() (bool, error) {
	releases, err := this.listReleases(action.ListAll)
	if err != nil {
		return false, err
	}

	for _, r := range releases {
		if r.Name == releaseName && r.Namespace == operatorNamespace {
			return true, nil
		}
	}

	return false, nil
}

// listReleases lists all releases that match the given state.
func (this *OperatorHandler) listReleases(state action.ListStates) ([]*release.Release, error) {
	client := action.NewList(this.configuration)
	client.StateMask = state

	return client.Run()
}

func (this *OperatorHandler) values(token string) map[string]interface{} {
	return map[string]interface{}{
		"secrets": map[string]string{
			"deployToken": token,
		},
		"consoleUrl": fmt.Sprintf("%s/ext/gql", this.url),
	}
}

func (this *OperatorHandler) InstallOrUpgrade(token string) error {
	exists, err := this.chartExists()
	if err != nil {
		return err
	}

	if exists {
		return this.Upgrade(token)
	}

	return this.Install(token)
}

func (this *OperatorHandler) Install(token string) error {
	_, err := this.install.Run(this.chart, this.values(token))
	return err
}

func (this *OperatorHandler) Upgrade(token string) error {
	_, err := this.upgrade.Run(releaseName, this.chart, this.values(token))
	return err
}

func (this *OperatorHandler) Uninstall() error {
	_, err := this.uninstall.Run(releaseName)
	return err
}

func NewOperatorHandler(ctx context.Context, kubeconfig *model.Kubeconfig, consoleUrl string) (*OperatorHandler, error) {
	handler := &OperatorHandler{
		ctx:        ctx,
		kubeconfig: kubeconfig,
		url:        consoleUrl,
	}

	err := handler.init()
	if err != nil {
		return nil, err
	}

	return handler, nil
}
