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

	in.Configuration.From(response.Configuration, ctx, d)
}

type WorkbenchToolConfiguration struct {
	HTTP       *WorkbenchToolHTTPConfig       `tfsdk:"http"`
	Elastic    *WorkbenchToolElasticConfig    `tfsdk:"elastic"`
	Prometheus *WorkbenchToolTokenAuthConfig  `tfsdk:"prometheus"`
	Loki       *WorkbenchToolTokenAuthConfig  `tfsdk:"loki"`
	Splunk     *WorkbenchToolTokenAuthConfig  `tfsdk:"splunk"`
	Tempo      *WorkbenchToolTokenAuthConfig  `tfsdk:"tempo"`
	Jaeger     *WorkbenchToolJaegerConfig     `tfsdk:"jaeger"`
	Datadog    *WorkbenchToolDatadogConfig    `tfsdk:"datadog"`
	Dynatrace  *WorkbenchToolDynatraceConfig  `tfsdk:"dynatrace"`
	Cloudwatch *WorkbenchToolCloudwatchConfig `tfsdk:"cloudwatch"`
	Azure      *WorkbenchToolAzureConfig      `tfsdk:"azure"`
	Linear     *WorkbenchToolLinearConfig     `tfsdk:"linear"`
	Atlassian  *WorkbenchToolAtlassianConfig  `tfsdk:"atlassian"`
}

func (in *WorkbenchToolConfiguration) Attributes(ctx context.Context) *gqlclient.WorkbenchToolConfigurationAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolConfigurationAttributes{
		HTTP:       in.HTTP.Attributes(ctx),
		Elastic:    in.Elastic.Attributes(),
		Prometheus: in.Prometheus.Attributes(),
		Loki:       in.Loki.lokiAttributes(),
		Splunk:     in.Splunk.splunkAttributes(),
		Tempo:      in.Tempo.tempoAttributes(),
		Jaeger:     in.Jaeger.Attributes(),
		Datadog:    in.Datadog.Attributes(),
		Dynatrace:  in.Dynatrace.Attributes(),
		Cloudwatch: in.Cloudwatch.Attributes(ctx),
		Azure:      in.Azure.Attributes(),
		Linear:     in.Linear.Attributes(),
		Atlassian:  in.Atlassian.Attributes(),
	}
}

func (in *WorkbenchToolConfiguration) From(configuration *gqlclient.WorkbenchToolFragment_Configuration, ctx context.Context, d *diag.Diagnostics) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.HTTP.From(configuration.HTTP, ctx, d)
	in.Elastic.From(configuration.Elastic)
	in.Prometheus.FromPrometheus(configuration.Prometheus)
	in.Loki.FromLoki(configuration.Loki)
	in.Splunk.FromSplunk(configuration.Splunk)
	in.Tempo.FromTempo(configuration.Tempo)
	in.Jaeger.From(configuration.Jaeger)
	in.Datadog.From(configuration.Datadog)
	in.Dynatrace.From(configuration.Dynatrace)
	in.Cloudwatch.From(configuration.Cloudwatch, ctx, d)
	in.Azure.From(configuration.Azure)
	in.Linear.From(configuration.Linear)
	in.Atlassian.From(configuration.Atlassian)
}

type WorkbenchToolHTTPConfig struct {
	URL         types.String `tfsdk:"url"`
	Method      types.String `tfsdk:"method"`
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
		URL:    in.URL.ValueString(),
		Method: gqlclient.WorkbenchToolHTTPMethod(strings.ToUpper(in.Method.ValueString())),
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

	in.URL = types.StringPointerValue(configuration.URL)
	if configuration.Method != nil {
		in.Method = types.StringValue(strings.ToUpper(*configuration.Method))
	} else {
		in.Method = types.StringNull()
	}

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

type WorkbenchToolTokenAuthConfig struct {
	URL      types.String `tfsdk:"url"`
	Token    types.String `tfsdk:"token"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	TenantID types.String `tfsdk:"tenant_id"`
}

func (in *WorkbenchToolTokenAuthConfig) Attributes() *gqlclient.WorkbenchToolPrometheusConnectionAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolPrometheusConnectionAttributes{
		URL:      in.URL.ValueString(),
		Token:    in.Token.ValueStringPointer(),
		Username: in.Username.ValueStringPointer(),
		Password: in.Password.ValueStringPointer(),
		TenantID: in.TenantID.ValueStringPointer(),
	}
}

func (in *WorkbenchToolTokenAuthConfig) FromPrometheus(configuration *gqlclient.WorkbenchToolFragment_Configuration_Prometheus) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = types.StringPointerValue(configuration.URL)
	in.Username = types.StringPointerValue(configuration.Username)
	in.TenantID = types.StringPointerValue(configuration.TenantID)
}

func (in *WorkbenchToolTokenAuthConfig) lokiAttributes() *gqlclient.WorkbenchToolLokiConnectionAttributes {
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

	in.URL = types.StringPointerValue(configuration.URL)
	in.Username = types.StringPointerValue(configuration.Username)
	in.TenantID = types.StringPointerValue(configuration.TenantID)
}

func (in *WorkbenchToolTokenAuthConfig) splunkAttributes() *gqlclient.WorkbenchToolSplunkConnectionAttributes {
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

func (in *WorkbenchToolTokenAuthConfig) FromSplunk(configuration *gqlclient.WorkbenchToolFragment_Configuration_Splunk) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.URL = types.StringPointerValue(configuration.URL)
	in.Username = types.StringPointerValue(configuration.Username)
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

	in.URL = types.StringPointerValue(configuration.URL)
	in.Username = types.StringPointerValue(configuration.Username)
	in.TenantID = types.StringPointerValue(configuration.TenantID)
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

	in.URL = types.StringPointerValue(configuration.URL)
	in.Username = types.StringPointerValue(configuration.Username)
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

	in.Site = types.StringPointerValue(configuration.Site)
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

	in.URL = types.StringPointerValue(configuration.URL)
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

	in.Region = types.StringPointerValue(configuration.Region)
	in.LogGroupNames = common.SetFrom(configuration.LogGroupNames, in.LogGroupNames, ctx, d)
	in.RoleArn = types.StringPointerValue(configuration.RoleArn)
	in.RoleSessionName = types.StringPointerValue(configuration.RoleSessionName)
}

type WorkbenchToolAzureConfig struct {
	SubscriptionID types.String `tfsdk:"subscription_id"`
	TenantID       types.String `tfsdk:"tenant_id"`
	ClientID       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
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
	}
}

func (in *WorkbenchToolAzureConfig) From(configuration *gqlclient.WorkbenchToolFragment_Configuration_Azure) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.SubscriptionID = types.StringPointerValue(configuration.SubscriptionID)
	in.TenantID = types.StringPointerValue(configuration.TenantID)
	in.ClientID = types.StringPointerValue(configuration.ClientID)
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

	in.Email = types.StringPointerValue(configuration.Email)
}
