package model

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

type WorkbenchTool struct {
	Id                types.String                `tfsdk:"id"`
	Name              types.String                `tfsdk:"name"`
	Tool              types.String                `tfsdk:"tool"`
	Categories        types.Set                   `tfsdk:"categories"`
	ProjectID         types.String                `tfsdk:"project_id"`
	McpServerID       types.String                `tfsdk:"mcp_server_id"`
	CloudConnectionID types.String                `tfsdk:"cloud_connection_id"`
	ScmConnectionID   types.String                `tfsdk:"scm_connection_id"`
	Configuration     *WorkbenchToolConfiguration `tfsdk:"configuration"`
}

func (in *WorkbenchTool) Attributes(ctx context.Context) (*gqlclient.WorkbenchToolAttributes, error) {
	categories := make([]types.String, len(in.Categories.Elements()))
	in.Categories.ElementsAs(ctx, &categories, false)

	return &gqlclient.WorkbenchToolAttributes{
		Name: in.Name.ValueString(),
		Tool: gqlclient.WorkbenchToolType(in.Tool.ValueString()),
		Categories: lo.Map(categories, func(v types.String, _ int) *gqlclient.WorkbenchToolCategory {
			return lo.ToPtr(gqlclient.WorkbenchToolCategory(v.ValueString()))
		}),
		ProjectID:         in.ProjectID.ValueStringPointer(),
		McpServerID:       in.McpServerID.ValueStringPointer(),
		CloudConnectionID: in.CloudConnectionID.ValueStringPointer(),
		ScmConnectionID:   in.ScmConnectionID.ValueStringPointer(),
		Configuration:     in.Configuration.Attributes(ctx),
	}, nil
}

func (in *WorkbenchTool) From(response *gqlclient.WorkbenchToolFragment, ctx context.Context, d *diag.Diagnostics) {
	if in == nil {
		return
	}
	if response == nil {
		return
	}

	in.Id = types.StringValue(response.ID)
	in.Name = types.StringValue(response.Name)
	in.Tool = types.StringValue(string(response.Tool))
	in.Categories = common.SetFrom(lo.Map(response.Categories, func(v *gqlclient.WorkbenchToolCategory, _ int) *string {
		return lo.Ternary(v == nil, nil, lo.ToPtr(string(*v)))
	}), in.Categories, ctx, d)
	in.ProjectID = common.ProjectFrom(response.Project)
	if response.McpServer != nil {
		in.McpServerID = types.StringValue(response.McpServer.ID)
	} else {
		in.McpServerID = types.StringNull()
	}
	if response.CloudConnection != nil {
		in.CloudConnectionID = types.StringValue(response.CloudConnection.ID)
	} else {
		in.CloudConnectionID = types.StringNull()
	}
	if response.ScmConnection != nil {
		in.ScmConnectionID = types.StringValue(response.ScmConnection.ID)
	} else {
		in.ScmConnectionID = types.StringNull()
	}

	if response.Configuration != nil {
		ensure(&in.Configuration)
		in.Configuration.From(response.Configuration, ctx, d)
	}
}

// FromCreate maps server-managed fields while preserving the exact planned
// configuration. Create responses can normalize optional values and omit
// write-only secrets, neither of which may replace known plan values.
func (in *WorkbenchTool) FromCreate(response *gqlclient.WorkbenchToolFragment, ctx context.Context, d *diag.Diagnostics) {
	if in == nil || response == nil {
		return
	}

	planned := in.Configuration
	in.Configuration = nil
	in.From(response, ctx, d)
	server := in.Configuration
	in.Configuration = planned

	if planned != nil {
		planned.computedFrom(server)
	}
}

type WorkbenchToolConfiguration struct {
	HTTP                *WorkbenchToolHTTPConfig                `tfsdk:"http"`
	Elastic             *WorkbenchToolElasticConfig             `tfsdk:"elastic"`
	Opensearch          *WorkbenchToolOpensearchConfig          `tfsdk:"opensearch"`
	Prometheus          *WorkbenchToolPrometheusConfig          `tfsdk:"prometheus"`
	Loki                *WorkbenchToolTokenAuthConfig           `tfsdk:"loki"`
	Splunk              *WorkbenchToolSplunkConfig              `tfsdk:"splunk"`
	Tempo               *WorkbenchToolTokenAuthConfig           `tfsdk:"tempo"`
	Jaeger              *WorkbenchToolJaegerConfig              `tfsdk:"jaeger"`
	Datadog             *WorkbenchToolDatadogConfig             `tfsdk:"datadog"`
	Dynatrace           *WorkbenchToolDynatraceConfig           `tfsdk:"dynatrace"`
	Cloudwatch          *WorkbenchToolCloudwatchConfig          `tfsdk:"cloudwatch"`
	Azure               *WorkbenchToolAzureConfig               `tfsdk:"azure"`
	Sentry              *WorkbenchToolSentryConfig              `tfsdk:"sentry"`
	Linear              *WorkbenchToolLinearConfig              `tfsdk:"linear"`
	Slack               *WorkbenchToolSlackConfig               `tfsdk:"slack"`
	Pagerduty           *WorkbenchToolPagerdutyConfig           `tfsdk:"pagerduty"`
	Teams               *WorkbenchToolTeamsConfig               `tfsdk:"teams"`
	Atlassian           *WorkbenchToolAtlassianConfig           `tfsdk:"atlassian"`
	Exa                 *WorkbenchToolExaConfig                 `tfsdk:"exa"`
	Github              *WorkbenchToolGithubConfig              `tfsdk:"github"`
	Gitlab              *WorkbenchToolGitlabConfig              `tfsdk:"gitlab"`
	Bitbucket           *WorkbenchToolBitbucketConfig           `tfsdk:"bitbucket"`
	BitbucketDatacenter *WorkbenchToolBitbucketDatacenterConfig `tfsdk:"bitbucket_datacenter"`
	AzureDevops         *WorkbenchToolAzureDevopsConfig         `tfsdk:"azure_devops"`
	Lambda              *WorkbenchToolLambdaConfig              `tfsdk:"lambda"`
	CloudRun            *WorkbenchToolCloudRunConfig            `tfsdk:"cloud_run"`
	AzureFunction       *WorkbenchToolAzureFunctionConfig       `tfsdk:"azure_function"`
	Docker              *WorkbenchToolDockerConfig              `tfsdk:"docker"`
}

func (in *WorkbenchToolConfiguration) computedFrom(server *WorkbenchToolConfiguration) {
	if in == nil || server == nil {
		return
	}
	if in.HTTP != nil && server.HTTP != nil {
		in.HTTP.Function = server.HTTP.Function
	}
	if in.Prometheus != nil && server.Prometheus != nil {
		in.Prometheus.AWSSigv4 = server.Prometheus.AWSSigv4
	}
	if in.Opensearch != nil && server.Opensearch != nil {
		in.Opensearch.UsePodIdentity = server.Opensearch.UsePodIdentity
	}
}

func (in *WorkbenchToolConfiguration) Attributes(ctx context.Context) *gqlclient.WorkbenchToolConfigurationAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolConfigurationAttributes{
		HTTP:                in.HTTP.Attributes(ctx),
		Elastic:             in.Elastic.Attributes(),
		Opensearch:          in.Opensearch.Attributes(),
		Prometheus:          in.Prometheus.Attributes(),
		Loki:                in.Loki.Attributes(),
		Splunk:              in.Splunk.Attributes(),
		Tempo:               in.Tempo.tempoAttributes(),
		Jaeger:              in.Jaeger.Attributes(),
		Datadog:             in.Datadog.Attributes(),
		Dynatrace:           in.Dynatrace.Attributes(),
		Cloudwatch:          in.Cloudwatch.Attributes(ctx),
		Azure:               in.Azure.Attributes(),
		Sentry:              in.Sentry.Attributes(),
		Linear:              in.Linear.Attributes(),
		Slack:               in.Slack.Attributes(),
		Pagerduty:           in.Pagerduty.Attributes(),
		Teams:               in.Teams.Attributes(),
		Atlassian:           in.Atlassian.Attributes(),
		Exa:                 in.Exa.Attributes(),
		Github:              in.Github.Attributes(),
		Gitlab:              in.Gitlab.Attributes(),
		Bitbucket:           in.Bitbucket.Attributes(),
		BitbucketDatacenter: in.BitbucketDatacenter.Attributes(),
		AzureDevops:         in.AzureDevops.Attributes(),
		Lambda:              in.Lambda.Attributes(),
		CloudRun:            in.CloudRun.Attributes(),
		AzureFunction:       in.AzureFunction.Attributes(),
		Docker:              in.Docker.Attributes(),
	}
}

func (in *WorkbenchToolConfiguration) From(configuration *gqlclient.WorkbenchToolFragment_Configuration, ctx context.Context, d *diag.Diagnostics) {
	if in == nil || configuration == nil {
		return
	}

	in.fromObservability(configuration, ctx, d)
	in.fromIntegrations(configuration)
	in.fromCloud(configuration, d)
}

func (in *WorkbenchToolConfiguration) fromObservability(configuration *gqlclient.WorkbenchToolFragment_Configuration, ctx context.Context, d *diag.Diagnostics) {
	if configuration.HTTP != nil {
		ensure(&in.HTTP)
		in.HTTP.From(configuration.HTTP, ctx, d)
	}
	if configuration.Elastic != nil {
		ensure(&in.Elastic)
		in.Elastic.From(configuration.Elastic)
	}
	if configuration.Opensearch != nil {
		ensure(&in.Opensearch)
		in.Opensearch.From(configuration.Opensearch)
	}
	if configuration.Prometheus != nil {
		ensure(&in.Prometheus)
		in.Prometheus.From(configuration.Prometheus)
	}
	if configuration.Loki != nil {
		ensure(&in.Loki)
		in.Loki.FromLoki(configuration.Loki)
	}
	if configuration.Splunk != nil {
		ensure(&in.Splunk)
		in.Splunk.From(configuration.Splunk)
	}
	if configuration.Tempo != nil {
		ensure(&in.Tempo)
		in.Tempo.FromTempo(configuration.Tempo)
	}
	if configuration.Jaeger != nil {
		ensure(&in.Jaeger)
		in.Jaeger.From(configuration.Jaeger)
	}
	if configuration.Datadog != nil {
		ensure(&in.Datadog)
		in.Datadog.From(configuration.Datadog)
	}
	if configuration.Dynatrace != nil {
		ensure(&in.Dynatrace)
		in.Dynatrace.From(configuration.Dynatrace)
	}
	if configuration.Cloudwatch != nil {
		ensure(&in.Cloudwatch)
		in.Cloudwatch.From(configuration.Cloudwatch, ctx, d)
	}
	if configuration.Azure != nil {
		ensure(&in.Azure)
		in.Azure.From(configuration.Azure)
	}
	if configuration.Sentry != nil {
		ensure(&in.Sentry)
		in.Sentry.From(configuration.Sentry)
	}
}

func (in *WorkbenchToolConfiguration) fromIntegrations(configuration *gqlclient.WorkbenchToolFragment_Configuration) {
	if configuration.Linear != nil {
		ensure(&in.Linear)
		in.Linear.From(configuration.Linear)
	}
	if configuration.Slack != nil {
		ensure(&in.Slack)
		in.Slack.From(configuration.Slack)
	}
	if configuration.Pagerduty != nil {
		ensure(&in.Pagerduty)
		in.Pagerduty.From(configuration.Pagerduty)
	}
	if configuration.Teams != nil {
		ensure(&in.Teams)
		in.Teams.From(configuration.Teams)
	}
	if configuration.Atlassian != nil {
		ensure(&in.Atlassian)
		in.Atlassian.From(configuration.Atlassian)
	}
	if configuration.Exa != nil {
		ensure(&in.Exa)
		in.Exa.From(configuration.Exa)
	}
	if configuration.Github != nil {
		ensure(&in.Github)
		in.Github.From(configuration.Github)
	}
	if configuration.Gitlab != nil {
		ensure(&in.Gitlab)
		in.Gitlab.From(configuration.Gitlab)
	}
	if configuration.Bitbucket != nil {
		ensure(&in.Bitbucket)
		in.Bitbucket.From(configuration.Bitbucket)
	}
	if configuration.BitbucketDatacenter != nil {
		ensure(&in.BitbucketDatacenter)
		in.BitbucketDatacenter.From(configuration.BitbucketDatacenter)
	}
	if configuration.AzureDevops != nil {
		ensure(&in.AzureDevops)
		in.AzureDevops.From(configuration.AzureDevops)
	}
	if configuration.Docker != nil {
		ensure(&in.Docker)
		in.Docker.From(configuration.Docker)
	}
}

func (in *WorkbenchToolConfiguration) fromCloud(configuration *gqlclient.WorkbenchToolFragment_Configuration, d *diag.Diagnostics) {
	if configuration.Lambda != nil {
		ensure(&in.Lambda)
		in.Lambda.From(configuration.Lambda, d)
	}
	if configuration.CloudRun != nil {
		ensure(&in.CloudRun)
		in.CloudRun.From(configuration.CloudRun, d)
	}
	if configuration.AzureFunction != nil {
		ensure(&in.AzureFunction)
		in.AzureFunction.From(configuration.AzureFunction, d)
	}
}

type WorkbenchToolHTTPConfig struct {
	URL         types.String `tfsdk:"url"`
	Method      types.String `tfsdk:"method"`
	Function    types.Bool   `tfsdk:"function"`
	Headers     types.Map    `tfsdk:"headers"`
	Body        types.String `tfsdk:"body"`
	InputSchema types.String `tfsdk:"input_schema"`
}

func (in *WorkbenchToolHTTPConfig) Attributes(ctx context.Context) *gqlclient.WorkbenchToolHTTPConfigurationAttributes {
	if in == nil {
		return nil
	}

	headers := make(map[string]types.String, len(in.Headers.Elements()))
	in.Headers.ElementsAs(ctx, &headers, false)

	return &gqlclient.WorkbenchToolHTTPConfigurationAttributes{
		URL:      in.URL.ValueString(),
		Method:   gqlclient.WorkbenchToolHTTPMethod(strings.ToUpper(in.Method.ValueString())),
		Function: in.Function.ValueBoolPointer(),
		Headers: lo.MapToSlice(headers, func(k string, v types.String) *gqlclient.WorkbenchToolHTTPHeaderAttributes {
			return &gqlclient.WorkbenchToolHTTPHeaderAttributes{Name: &k, Value: v.ValueStringPointer()}
		}),
		Body:        in.Body.ValueStringPointer(),
		InputSchema: in.InputSchema.ValueStringPointer(),
	}
}

func (in *WorkbenchToolHTTPConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_HTTP, ctx context.Context, d *diag.Diagnostics) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = stringFromAPI(configuration.URL, in.URL)
	if configuration.Method != nil {
		in.Method = types.StringValue(strings.ToUpper(*configuration.Method))
	} else {
		in.Method = types.StringNull()
	}
	in.Function = types.BoolPointerValue(configuration.Function)

	if configuration.Headers != nil {
		headers := make(map[string]any, len(configuration.Headers))
		for _, v := range configuration.Headers {
			if v.Value != nil {
				headers[*v.Name] = *v.Value
			}
		}

		in.Headers = common.MapFromWithConfig(headers, in.Headers, ctx, d)
	} else {
		in.Headers = types.MapNull(types.StringType)
	}

	if configuration.Body != nil {
		in.Body = types.StringPointerValue(configuration.Body)
	} else {
		in.Body = types.StringNull()
	}

	if configuration.InputSchema != nil {
		inputSchema, err := json.Marshal(configuration.InputSchema)
		if err != nil {
			d.AddError("Provider Error", fmt.Sprintf("Cannot marshal input schema, got error: %s", err))
			return
		}

		in.InputSchema = types.StringValue(string(inputSchema))
	} else {
		in.InputSchema = types.StringNull()
	}
}

type WorkbenchToolElasticConfig struct {
	URL      types.String `tfsdk:"url"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Index    types.String `tfsdk:"index"`
}

func (in *WorkbenchToolElasticConfig) Attributes() *gqlclient.WorkbenchToolElasticConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolElasticConnectionAttributes{
		URL:      in.URL.ValueString(),
		Username: in.Username.ValueString(),
		Password: in.Password.ValueStringPointer(),
		Index:    in.Index.ValueString(),
	}
}

func (in *WorkbenchToolElasticConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Elastic) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = types.StringValue(configuration.URL)
	in.Username = types.StringValue(configuration.Username)
	in.Index = types.StringValue(configuration.Index)
}

type WorkbenchToolPrometheusConfig struct {
	URL                types.String `tfsdk:"url"`
	Token              types.String `tfsdk:"token"`
	Username           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	TenantID           types.String `tfsdk:"tenant_id"`
	AWSSigv4           types.Bool   `tfsdk:"aws_sigv4"`
	AWSAccessKeyID     types.String `tfsdk:"aws_access_key_id"`
	AWSSecretAccessKey types.String `tfsdk:"aws_secret_access_key"`
	AWSRegion          types.String `tfsdk:"aws_region"`
}

func (in *WorkbenchToolPrometheusConfig) Attributes() *gqlclient.WorkbenchToolPrometheusConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolPrometheusConnectionAttributes{
		URL:                in.URL.ValueString(),
		Token:              in.Token.ValueStringPointer(),
		Username:           in.Username.ValueStringPointer(),
		Password:           in.Password.ValueStringPointer(),
		TenantID:           in.TenantID.ValueStringPointer(),
		AWSSigv4:           in.AWSSigv4.ValueBoolPointer(),
		AWSAccessKeyID:     in.AWSAccessKeyID.ValueStringPointer(),
		AWSSecretAccessKey: in.AWSSecretAccessKey.ValueStringPointer(),
		AWSRegion:          in.AWSRegion.ValueStringPointer(),
	}
}

func (in *WorkbenchToolPrometheusConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Prometheus) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = stringFromAPI(configuration.URL, in.URL)
	in.Username = stringFromAPI(configuration.Username, in.Username)
	in.TenantID = stringFromAPI(configuration.TenantID, in.TenantID)
	in.AWSSigv4 = types.BoolPointerValue(configuration.AWSSigv4)
	in.AWSRegion = stringFromAPI(configuration.AWSRegion, in.AWSRegion)
	// Token, Password, AWSAccessKeyID, AWSSecretAccessKey are never returned; keep configured values.
}

type WorkbenchToolTokenAuthConfig struct {
	URL      types.String `tfsdk:"url"`
	Token    types.String `tfsdk:"token"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	TenantID types.String `tfsdk:"tenant_id"`
}

func (in *WorkbenchToolTokenAuthConfig) Attributes() *gqlclient.WorkbenchToolLokiConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolLokiConnectionAttributes{
		URL:      in.URL.ValueString(),
		Token:    in.Token.ValueStringPointer(),
		Username: in.Username.ValueStringPointer(),
		Password: in.Password.ValueStringPointer(),
		TenantID: in.TenantID.ValueStringPointer(),
	}
}

func (in *WorkbenchToolTokenAuthConfig) FromLoki(configuration *gqlclient.WorkbenchToolFragment_Configuration_Loki) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = stringFromAPI(configuration.URL, in.URL)
	in.Username = stringFromAPI(configuration.Username, in.Username)
	in.TenantID = stringFromAPI(configuration.TenantID, in.TenantID)
}

type WorkbenchToolSplunkConfig struct {
	URL      types.String `tfsdk:"url"`
	Token    types.String `tfsdk:"token"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (in *WorkbenchToolSplunkConfig) Attributes() *gqlclient.WorkbenchToolSplunkConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolSplunkConnectionAttributes{
		URL:      in.URL.ValueString(),
		Token:    in.Token.ValueStringPointer(),
		Username: in.Username.ValueStringPointer(),
		Password: in.Password.ValueStringPointer(),
	}
}

func (in *WorkbenchToolSplunkConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Splunk) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = stringFromAPI(configuration.URL, in.URL)
	in.Username = stringFromAPI(configuration.Username, in.Username)
}

func (in *WorkbenchToolTokenAuthConfig) tempoAttributes() *gqlclient.WorkbenchToolTempoConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolTempoConnectionAttributes{
		URL:      in.URL.ValueString(),
		Token:    in.Token.ValueStringPointer(),
		Username: in.Username.ValueStringPointer(),
		Password: in.Password.ValueStringPointer(),
		TenantID: in.TenantID.ValueStringPointer(),
	}
}

func (in *WorkbenchToolTokenAuthConfig) FromTempo(configuration *gqlclient.WorkbenchToolFragment_Configuration_Tempo) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = stringFromAPI(configuration.URL, in.URL)
	in.Username = stringFromAPI(configuration.Username, in.Username)
	in.TenantID = stringFromAPI(configuration.TenantID, in.TenantID)
}

type WorkbenchToolJaegerConfig struct {
	URL      types.String `tfsdk:"url"`
	Token    types.String `tfsdk:"token"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (in *WorkbenchToolJaegerConfig) Attributes() *gqlclient.WorkbenchToolJaegerConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolJaegerConnectionAttributes{
		URL:      in.URL.ValueString(),
		Token:    in.Token.ValueStringPointer(),
		Username: in.Username.ValueStringPointer(),
		Password: in.Password.ValueStringPointer(),
	}
}

func (in *WorkbenchToolJaegerConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Jaeger) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = stringFromAPI(configuration.URL, in.URL)
	in.Username = stringFromAPI(configuration.Username, in.Username)
}

type WorkbenchToolDatadogConfig struct {
	Site   types.String `tfsdk:"site"`
	APIKey types.String `tfsdk:"api_key"`
	AppKey types.String `tfsdk:"app_key"`
}

func (in *WorkbenchToolDatadogConfig) Attributes() *gqlclient.WorkbenchToolDatadogConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolDatadogConnectionAttributes{
		Site:   in.Site.ValueStringPointer(),
		APIKey: in.APIKey.ValueStringPointer(),
		AppKey: in.AppKey.ValueStringPointer(),
	}
}

func (in *WorkbenchToolDatadogConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Datadog) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.Site = stringFromAPI(configuration.Site, in.Site)
}

type WorkbenchToolDynatraceConfig struct {
	URL           types.String `tfsdk:"url"`
	PlatformToken types.String `tfsdk:"platform_token"`
}

func (in *WorkbenchToolDynatraceConfig) Attributes() *gqlclient.WorkbenchToolDynatraceConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolDynatraceConnectionAttributes{
		URL:           in.URL.ValueString(),
		PlatformToken: in.PlatformToken.ValueString(),
	}
}

func (in *WorkbenchToolDynatraceConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Dynatrace) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = stringFromAPI(configuration.URL, in.URL)
}

type WorkbenchToolCloudwatchConfig struct {
	Region          types.String `tfsdk:"region"`
	LogGroupNames   types.Set    `tfsdk:"log_group_names"`
	AccessKeyID     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	RoleArn         types.String `tfsdk:"role_arn"`
	ExternalID      types.String `tfsdk:"external_id"`
	RoleSessionName types.String `tfsdk:"role_session_name"`
}

func (in *WorkbenchToolCloudwatchConfig) Attributes(ctx context.Context) *gqlclient.WorkbenchToolCloudwatchConnectionAttributes {
	if in == nil {
		return nil
	}
	logGroupNames := make([]types.String, len(in.LogGroupNames.Elements()))
	in.LogGroupNames.ElementsAs(ctx, &logGroupNames, false)

	return &gqlclient.WorkbenchToolCloudwatchConnectionAttributes{
		Region: in.Region.ValueString(),
		LogGroupNames: lo.FilterMap(
			logGroupNames,
			func(value types.String, _ int) (*string, bool) {
				if value.IsNull() || value.IsUnknown() {
					return nil, false
				}
				str := value.ValueString()
				return &str, true
			},
		),
		AccessKeyID:     in.AccessKeyID.ValueStringPointer(),
		SecretAccessKey: in.SecretAccessKey.ValueStringPointer(),
		RoleArn:         in.RoleArn.ValueStringPointer(),
		ExternalID:      in.ExternalID.ValueStringPointer(),
		RoleSessionName: in.RoleSessionName.ValueStringPointer(),
	}
}

func (in *WorkbenchToolCloudwatchConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Cloudwatch, ctx context.Context, d *diag.Diagnostics) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.Region = stringFromAPI(configuration.Region, in.Region)
	in.LogGroupNames = common.SetFrom(configuration.LogGroupNames, in.LogGroupNames, ctx, d)
	in.RoleArn = stringFromAPI(configuration.RoleArn, in.RoleArn)
	in.RoleSessionName = stringFromAPI(configuration.RoleSessionName, in.RoleSessionName)
}

type WorkbenchToolAzureConfig struct {
	SubscriptionID types.String `tfsdk:"subscription_id"`
	TenantID       types.String `tfsdk:"tenant_id"`
	ClientID       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
	PrometheusURL  types.String `tfsdk:"prometheus_url"`
}

func (in *WorkbenchToolAzureConfig) Attributes() *gqlclient.WorkbenchToolAzureConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolAzureConnectionAttributes{
		SubscriptionID: in.SubscriptionID.ValueString(),
		TenantID:       in.TenantID.ValueString(),
		ClientID:       in.ClientID.ValueString(),
		ClientSecret:   in.ClientSecret.ValueString(),
		PrometheusURL:  in.PrometheusURL.ValueStringPointer(),
	}
}

func (in *WorkbenchToolAzureConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Azure) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.SubscriptionID = stringFromAPI(configuration.SubscriptionID, in.SubscriptionID)
	in.TenantID = stringFromAPI(configuration.TenantID, in.TenantID)
	in.ClientID = stringFromAPI(configuration.ClientID, in.ClientID)
	in.PrometheusURL = stringFromAPI(configuration.PrometheusURL, in.PrometheusURL)
}

type WorkbenchToolLinearConfig struct {
	AccessToken types.String `tfsdk:"access_token"`
}

func (in *WorkbenchToolLinearConfig) Attributes() *gqlclient.WorkbenchToolLinearConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolLinearConnectionAttributes{
		AccessToken: in.AccessToken.ValueStringPointer(),
	}
}

func (in *WorkbenchToolLinearConfig) From(_ *gqlclient.WorkbenchToolFragment_Configuration_Linear) {
	// Linear fragment only exposes URL (no credentials); keep configured token value.
}

type WorkbenchToolAtlassianConfig struct {
	ServiceAccount types.String `tfsdk:"service_account"`
	APIToken       types.String `tfsdk:"api_token"`
	Email          types.String `tfsdk:"email"`
}

func (in *WorkbenchToolAtlassianConfig) Attributes() *gqlclient.WorkbenchToolAtlassianConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolAtlassianConnectionAttributes{
		ServiceAccount: in.ServiceAccount.ValueStringPointer(),
		APIToken:       in.APIToken.ValueStringPointer(),
		Email:          in.Email.ValueStringPointer(),
	}
}

func (in *WorkbenchToolAtlassianConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Atlassian) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.Email = stringFromAPI(configuration.Email, in.Email)
}

type WorkbenchToolOpensearchConfig struct {
	Host               types.String `tfsdk:"host"`
	Index              types.String `tfsdk:"index"`
	AWSAccessKeyID     types.String `tfsdk:"aws_access_key_id"`
	AWSSecretAccessKey types.String `tfsdk:"aws_secret_access_key"`
	AWSRegion          types.String `tfsdk:"aws_region"`
	AssumeRoleArn      types.String `tfsdk:"assume_role_arn"`
	UsePodIdentity     types.Bool   `tfsdk:"use_pod_identity"`
}

func (in *WorkbenchToolOpensearchConfig) Attributes() *gqlclient.WorkbenchToolOpensearchConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolOpensearchConnectionAttributes{
		Host:               in.Host.ValueString(),
		Index:              in.Index.ValueString(),
		AWSAccessKeyID:     in.AWSAccessKeyID.ValueStringPointer(),
		AWSSecretAccessKey: in.AWSSecretAccessKey.ValueStringPointer(),
		AWSRegion:          in.AWSRegion.ValueStringPointer(),
		AssumeRoleArn:      in.AssumeRoleArn.ValueStringPointer(),
		UsePodIdentity:     in.UsePodIdentity.ValueBoolPointer(),
	}
}

func (in *WorkbenchToolOpensearchConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Opensearch) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.Host = types.StringValue(configuration.Host)
	in.Index = types.StringValue(configuration.Index)
	in.AWSRegion = stringFromAPI(configuration.AWSRegion, in.AWSRegion)
	in.AssumeRoleArn = stringFromAPI(configuration.AssumeRoleArn, in.AssumeRoleArn)
	in.UsePodIdentity = types.BoolPointerValue(configuration.UsePodIdentity)
	// AWSAccessKeyID and AWSSecretAccessKey are never returned; keep configured values.
}

type WorkbenchToolSentryConfig struct {
	URL         types.String `tfsdk:"url"`
	AccessToken types.String `tfsdk:"access_token"`
}

func (in *WorkbenchToolSentryConfig) Attributes() *gqlclient.WorkbenchToolSentryConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolSentryConnectionAttributes{
		URL:         in.URL.ValueStringPointer(),
		AccessToken: in.AccessToken.ValueStringPointer(),
	}
}

func (in *WorkbenchToolSentryConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Sentry) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = stringFromAPI(configuration.URL, in.URL)
	// AccessToken is never returned; keep configured value.
}

type WorkbenchToolSlackConfig struct {
	BotToken types.String `tfsdk:"bot_token"`
}

func (in *WorkbenchToolSlackConfig) Attributes() *gqlclient.WorkbenchToolSlackConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolSlackConnectionAttributes{
		BotToken: in.BotToken.ValueStringPointer(),
	}
}

func (in *WorkbenchToolSlackConfig) From(_ *gqlclient.WorkbenchToolFragment_Configuration_Slack) {
	// Slack fragment only exposes URL (no credentials); keep configured token value.
}

type WorkbenchToolPagerdutyConfig struct {
	APIToken types.String `tfsdk:"api_token"`
}

func (in *WorkbenchToolPagerdutyConfig) Attributes() *gqlclient.WorkbenchToolPagerdutyConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolPagerdutyConnectionAttributes{
		APIToken: in.APIToken.ValueStringPointer(),
	}
}

func (in *WorkbenchToolPagerdutyConfig) From(_ *gqlclient.WorkbenchToolFragment_Configuration_Pagerduty) {
	// PagerDuty fragment only exposes URL (no credentials); keep configured token value.
}

type WorkbenchToolTeamsConfig struct {
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	TenantID     types.String `tfsdk:"tenant_id"`
}

func (in *WorkbenchToolTeamsConfig) Attributes() *gqlclient.WorkbenchToolTeamsConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolTeamsConnectionAttributes{
		ClientID:     in.ClientID.ValueString(),
		ClientSecret: in.ClientSecret.ValueString(),
		TenantID:     in.TenantID.ValueString(),
	}
}

func (in *WorkbenchToolTeamsConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Teams) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.ClientID = stringFromAPI(configuration.ClientID, in.ClientID)
	in.TenantID = stringFromAPI(configuration.TenantID, in.TenantID)
	// ClientSecret is never returned; keep configured value.
}

type WorkbenchToolExaConfig struct {
	APIKey types.String `tfsdk:"api_key"`
}

func (in *WorkbenchToolExaConfig) Attributes() *gqlclient.WorkbenchToolExaConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolExaConnectionAttributes{
		APIKey: in.APIKey.ValueStringPointer(),
	}
}

func (in *WorkbenchToolExaConfig) From(_ *gqlclient.WorkbenchToolFragment_Configuration_Exa) {
	// Exa fragment only exposes URL (no credentials); keep configured API key.
}

type WorkbenchToolGithubConfig struct {
	URL            types.String `tfsdk:"url"`
	AccessToken    types.String `tfsdk:"access_token"`
	Toolset        types.String `tfsdk:"toolset"`
	AppID          types.String `tfsdk:"app_id"`
	InstallationID types.String `tfsdk:"installation_id"`
	PrivateKey     types.String `tfsdk:"private_key"`
}

func (in *WorkbenchToolGithubConfig) Attributes() *gqlclient.WorkbenchToolGithubConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolGithubConnectionAttributes{
		URL:            in.URL.ValueStringPointer(),
		AccessToken:    in.AccessToken.ValueStringPointer(),
		Toolset:        in.Toolset.ValueStringPointer(),
		AppID:          in.AppID.ValueStringPointer(),
		InstallationID: in.InstallationID.ValueStringPointer(),
		PrivateKey:     in.PrivateKey.ValueStringPointer(),
	}
}

func (in *WorkbenchToolGithubConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Github) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = stringFromAPI(lo.EmptyableToPtr(configuration.URL), in.URL)
	in.Toolset = stringFromAPI(configuration.Toolset, in.Toolset)
	in.AppID = stringFromAPI(configuration.AppID, in.AppID)
	in.InstallationID = stringFromAPI(configuration.InstallationID, in.InstallationID)
	// AccessToken and PrivateKey are never returned; keep configured values.
}

type WorkbenchToolGitlabConfig struct {
	URL   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

func (in *WorkbenchToolGitlabConfig) Attributes() *gqlclient.WorkbenchToolGitlabConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolGitlabConnectionAttributes{
		URL:   in.URL.ValueStringPointer(),
		Token: in.Token.ValueStringPointer(),
	}
}

func (in *WorkbenchToolGitlabConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Gitlab) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = stringFromAPI(configuration.URL, in.URL)
	// Token is never returned; keep configured value.
}

type WorkbenchToolBitbucketConfig struct {
	URL   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

func (in *WorkbenchToolBitbucketConfig) Attributes() *gqlclient.WorkbenchToolBitbucketConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolBitbucketConnectionAttributes{
		URL:   in.URL.ValueStringPointer(),
		Token: in.Token.ValueStringPointer(),
	}
}

func (in *WorkbenchToolBitbucketConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Bitbucket) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = stringFromAPI(configuration.URL, in.URL)
	// Token is never returned; keep configured value.
}

type WorkbenchToolBitbucketDatacenterConfig struct {
	URL   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

func (in *WorkbenchToolBitbucketDatacenterConfig) Attributes() *gqlclient.WorkbenchToolBitbucketDatacenterConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolBitbucketDatacenterConnectionAttributes{
		URL:   in.URL.ValueString(),
		Token: in.Token.ValueStringPointer(),
	}
}

func (in *WorkbenchToolBitbucketDatacenterConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_BitbucketDatacenter) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = stringFromAPI(configuration.URL, in.URL)
	// Token is never returned; keep configured value.
}

type WorkbenchToolAzureDevopsConfig struct {
	URL   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

func (in *WorkbenchToolAzureDevopsConfig) Attributes() *gqlclient.WorkbenchToolAzureDevopsConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolAzureDevopsConnectionAttributes{
		URL:   in.URL.ValueStringPointer(),
		Token: in.Token.ValueStringPointer(),
	}
}

func (in *WorkbenchToolAzureDevopsConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_AzureDevops) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = stringFromAPI(configuration.URL, in.URL)
	// Token is never returned; keep configured value.
}

type WorkbenchToolLambdaConfig struct {
	LambdaArn   types.String `tfsdk:"lambda_arn"`
	Description types.String `tfsdk:"description"`
	InputSchema types.String `tfsdk:"input_schema"`
}

func (in *WorkbenchToolLambdaConfig) Attributes() *gqlclient.WorkbenchToolLambdaConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolLambdaConnectionAttributes{
		LambdaArn:   in.LambdaArn.ValueString(),
		Description: in.Description.ValueString(),
		InputSchema: in.InputSchema.ValueStringPointer(),
	}
}

func (in *WorkbenchToolLambdaConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Lambda, d *diag.Diagnostics) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.LambdaArn = types.StringPointerValue(configuration.LambdaArn)
	in.Description = types.StringPointerValue(configuration.Description)
	if configuration.InputSchema != nil {
		inputSchema, err := json.Marshal(configuration.InputSchema)
		if err != nil {
			d.AddError("Provider Error", fmt.Sprintf("Cannot marshal lambda input schema, got error: %s", err))
			return
		}
		in.InputSchema = types.StringValue(string(inputSchema))
	} else {
		in.InputSchema = types.StringNull()
	}
}

type WorkbenchToolCloudRunConfig struct {
	Identifier  types.String `tfsdk:"identifier"`
	Description types.String `tfsdk:"description"`
	InputSchema types.String `tfsdk:"input_schema"`
}

func (in *WorkbenchToolCloudRunConfig) Attributes() *gqlclient.WorkbenchToolCloudRunConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolCloudRunConnectionAttributes{
		Identifier:  in.Identifier.ValueString(),
		Description: in.Description.ValueString(),
		InputSchema: in.InputSchema.ValueStringPointer(),
	}
}

func (in *WorkbenchToolCloudRunConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_CloudRun, d *diag.Diagnostics) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.Identifier = types.StringPointerValue(configuration.Identifier)
	in.Description = types.StringPointerValue(configuration.Description)
	if configuration.InputSchema != nil {
		inputSchema, err := json.Marshal(configuration.InputSchema)
		if err != nil {
			d.AddError("Provider Error", fmt.Sprintf("Cannot marshal cloud_run input schema, got error: %s", err))
			return
		}
		in.InputSchema = types.StringValue(string(inputSchema))
	} else {
		in.InputSchema = types.StringNull()
	}
}

type WorkbenchToolAzureFunctionConfig struct {
	Identifier  types.String `tfsdk:"identifier"`
	Description types.String `tfsdk:"description"`
	InputSchema types.String `tfsdk:"input_schema"`
}

func (in *WorkbenchToolAzureFunctionConfig) Attributes() *gqlclient.WorkbenchToolAzureFunctionConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolAzureFunctionConnectionAttributes{
		Identifier:  in.Identifier.ValueString(),
		Description: in.Description.ValueString(),
		InputSchema: in.InputSchema.ValueStringPointer(),
	}
}

func (in *WorkbenchToolAzureFunctionConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_AzureFunction, d *diag.Diagnostics) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.Identifier = types.StringPointerValue(configuration.Identifier)
	in.Description = types.StringPointerValue(configuration.Description)
	if configuration.InputSchema != nil {
		inputSchema, err := json.Marshal(configuration.InputSchema)
		if err != nil {
			d.AddError("Provider Error", fmt.Sprintf("Cannot marshal azure_function input schema, got error: %s", err))
			return
		}
		in.InputSchema = types.StringValue(string(inputSchema))
	} else {
		in.InputSchema = types.StringNull()
	}
}

type WorkbenchToolDockerConfig struct {
	URL      types.String             `tfsdk:"url"`
	Provider types.String             `tfsdk:"provider"`
	Auth     *WorkbenchToolDockerAuth `tfsdk:"auth"`
}

type WorkbenchToolDockerAuth struct {
	Proxy  *WorkbenchToolHTTPProxyConfig  `tfsdk:"proxy"`
	Basic  *WorkbenchToolDockerBasicAuth  `tfsdk:"basic"`
	Bearer *WorkbenchToolDockerBearerAuth `tfsdk:"bearer"`
	AWS    *WorkbenchToolDockerAWSAuth    `tfsdk:"aws"`
	Azure  *WorkbenchToolDockerAzureAuth  `tfsdk:"azure"`
	GCP    *WorkbenchToolDockerGCPAuth    `tfsdk:"gcp"`
}

type WorkbenchToolHTTPProxyConfig struct {
	URL     types.String `tfsdk:"url"`
	Noproxy types.String `tfsdk:"noproxy"`
}

type WorkbenchToolDockerBasicAuth struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type WorkbenchToolDockerBearerAuth struct {
	Token types.String `tfsdk:"token"`
}

type WorkbenchToolDockerAWSAuth struct {
	AccessKey       types.String `tfsdk:"access_key"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	AssumeRoleArn   types.String `tfsdk:"assume_role_arn"`
}

type WorkbenchToolDockerAzureAuth struct {
	ClientID       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
	TenantID       types.String `tfsdk:"tenant_id"`
	SubscriptionID types.String `tfsdk:"subscription_id"`
}

type WorkbenchToolDockerGCPAuth struct {
	ApplicationCredentials types.String `tfsdk:"application_credentials"`
}

func (in *WorkbenchToolDockerConfig) Attributes() *gqlclient.WorkbenchToolDockerConnectionAttributes {
	if in == nil {
		return nil
	}

	attrs := &gqlclient.WorkbenchToolDockerConnectionAttributes{
		URL: in.URL.ValueStringPointer(),
	}
	if !in.Provider.IsNull() && !in.Provider.IsUnknown() {
		attrs.Provider = lo.ToPtr(gqlclient.HelmAuthProvider(in.Provider.ValueString()))
	}
	attrs.Auth = in.Auth.Attributes()
	return attrs
}

func (in *WorkbenchToolDockerAuth) Attributes() *gqlclient.HelmAuthAttributes {
	if in == nil {
		return nil
	}

	if lo.EveryBy([]any{in.Proxy, in.Basic, in.Bearer, in.AWS, in.Azure, in.GCP}, lo.IsNil) {
		return nil
	}

	attrs := &gqlclient.HelmAuthAttributes{}
	if in.Proxy != nil {
		attrs.Proxy = &gqlclient.HTTPProxyAttributes{
			URL:     in.Proxy.URL.ValueString(),
			Noproxy: in.Proxy.Noproxy.ValueStringPointer(),
		}
	}
	if in.Basic != nil {
		attrs.Basic = &gqlclient.HelmBasicAuthAttributes{
			Username: in.Basic.Username.ValueString(),
			Password: in.Basic.Password.ValueString(),
		}
	}
	if in.Bearer != nil {
		attrs.Bearer = &gqlclient.HelmBearerAuthAttributes{
			Token: in.Bearer.Token.ValueString(),
		}
	}
	if in.AWS != nil {
		attrs.AWS = &gqlclient.HelmAWSAuthAttributes{
			AccessKey:       in.AWS.AccessKey.ValueStringPointer(),
			SecretAccessKey: in.AWS.SecretAccessKey.ValueStringPointer(),
			AssumeRoleArn:   in.AWS.AssumeRoleArn.ValueStringPointer(),
		}
	}
	if in.Azure != nil {
		attrs.Azure = &gqlclient.HelmAzureAuthAttributes{
			ClientID:       in.Azure.ClientID.ValueStringPointer(),
			ClientSecret:   in.Azure.ClientSecret.ValueStringPointer(),
			TenantID:       in.Azure.TenantID.ValueStringPointer(),
			SubscriptionID: in.Azure.SubscriptionID.ValueStringPointer(),
		}
	}
	if in.GCP != nil {
		attrs.GCP = &gqlclient.HelmGCPAuthAttributes{
			ApplicationCredentials: in.GCP.ApplicationCredentials.ValueStringPointer(),
		}
	}
	return attrs
}

func (in *WorkbenchToolDockerConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Docker) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = stringFromAPI(configuration.URL, in.URL)
	if configuration.Provider != nil {
		in.Provider = types.StringValue(string(*configuration.Provider))
	}
	// Proxy is returned at the docker connection level; credentials stay under auth.
	// Only refresh proxy when it was already configured so we don't reshape auth
	// objects that hold write-only secrets.
	if configuration.Proxy != nil && in.Auth != nil && in.Auth.Proxy != nil && !lo.IsEmpty(configuration.Proxy.URL) {
		in.Auth.Proxy.URL = types.StringValue(configuration.Proxy.URL)
		in.Auth.Proxy.Noproxy = stringFromAPI(configuration.Proxy.Noproxy, in.Auth.Proxy.Noproxy)
	}
	// Auth credentials are never returned; keep configured values.
}

// stringFromAPI keeps the planned/state value when the API omits or returns an empty string.
func stringFromAPI(api *string, current types.String) types.String {
	return lo.Ternary(lo.IsEmpty(lo.FromPtr(api)), current, types.StringValue(lo.FromPtr(api)))
}

// ensure allocates dst when nil so From can populate imported/nested config blocks.
func ensure[T any](dst **T) {
	if *dst == nil {
		*dst = new(T)
	}
}
