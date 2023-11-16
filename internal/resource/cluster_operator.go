package resource

import (
	"context"
	"fmt"
	"time"

	"github.com/pluralsh/plural-cli/pkg/helm"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"

	"terraform-provider-plural/internal/model"
)

const (
	// TODO: read from CLI pkg?
	operatorNamespace = "plrl-deploy-operator"
	releaseName       = "deploy-operator"
	chartName         = "deployment-operator"
	repoUrl           = "https://pluralsh.github.io/deployment-operator"
)

func doInstallOperator(ctx context.Context, kubeconfig model.Kubeconfig, consoleUrl, deployToken string) error {
	actionConfig := new(action.Configuration)
	vals := map[string]interface{}{
		"secrets": map[string]string{
			"deployToken": deployToken,
		},
		"consoleUrl": fmt.Sprintf("%s/ext/gql", consoleUrl),
	}

	kc, err := newKubeconfig(ctx, kubeconfig, lo.ToPtr(operatorNamespace))
	if err != nil {
		return err
	}

	if err = helm.AddRepo(releaseName, repoUrl); err != nil {
		return err
	}

	err = actionConfig.Init(kc, operatorNamespace, "", logrus.Debugf)
	if err != nil {
		return err
	}

	client := action.NewInstall(actionConfig)

	locateName := fmt.Sprintf("%s/%s", releaseName, chartName)
	path, err := client.ChartPathOptions.LocateChart(locateName, cli.New())
	if err != nil {
		return err
	}

	c, err := loader.Load(path)
	if err != nil {
		return err
	}

	client.Namespace = operatorNamespace
	client.ReleaseName = releaseName
	client.Timeout = 5 * time.Minute
	client.Wait = true
	client.CreateNamespace = true

	_, err = client.Run(c, vals)
	return err

	//histClient := action.NewHistory(actionConfig)
	//histClient.Max = 5
	//
	//if _, err = histClient.Run(releaseName); errors.Is(err, driver.ErrReleaseNotFound) {
	//	tflog.Info(ctx, "installing deployment operator...")
	//	instClient := action.NewInstall(actionConfig)
	//	instClient.Namespace = operatorNamespace
	//	instClient.ReleaseName = releaseName
	//	instClient.Timeout = time.Minute * 1
	//	instClient.Wait = true
	//	instClient.CreateNamespace = true
	//	_, err = instClient.Run(chart, vals)
	//	return err
	//}
	//
	//tflog.Info(ctx, "upgrading deployment operator...")
	//client := action.NewUpgrade(actionConfig)
	//client.Namespace = operatorNamespace
	//client.Timeout = time.Minute * 5
	//_, err = client.Run(releaseName, chart, vals)
	//return err
}
