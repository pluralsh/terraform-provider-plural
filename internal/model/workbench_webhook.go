package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
)

type WorkbenchWebhook struct {
	Id             types.String                 `tfsdk:"id"`
	WorkbenchID    types.String                 `tfsdk:"workbench_id"`
	Name           types.String                 `tfsdk:"name"`
	Prompt         types.String                 `tfsdk:"prompt"`
	WebhookID      types.String                 `tfsdk:"webhook_id"`
	IssueWebhookID types.String                 `tfsdk:"issue_webhook_id"`
	Matches        *WorkbenchWebhookMatchConfig `tfsdk:"matches"`
}

type WorkbenchWebhookMatchConfig struct {
	Regex           types.String `tfsdk:"regex"`
	Substring       types.String `tfsdk:"substring"`
	CaseInsensitive types.Bool   `tfsdk:"case_insensitive"`
}

func (in *WorkbenchWebhook) Attributes() gqlclient.WorkbenchWebhookAttributes {
	var matches *gqlclient.WorkbenchWebhookMatchesAttributes
	if in.Matches != nil {
		matches = &gqlclient.WorkbenchWebhookMatchesAttributes{
			Regex:           in.Matches.Regex.ValueStringPointer(),
			Substring:       in.Matches.Substring.ValueStringPointer(),
			CaseInsensitive: in.Matches.CaseInsensitive.ValueBoolPointer(),
		}
	}

	return gqlclient.WorkbenchWebhookAttributes{
		Name:           in.Name.ValueStringPointer(),
		WebhookID:      in.WebhookID.ValueStringPointer(),
		IssueWebhookID: in.IssueWebhookID.ValueStringPointer(),
		Prompt:         in.Prompt.ValueStringPointer(),
		Matches:        matches,
	}
}

func (in *WorkbenchWebhook) From(response *gqlclient.WorkbenchWebhookFragment) {
	if in == nil || response == nil {
		return
	}

	in.Id = types.StringValue(response.ID)
	in.Name = types.StringPointerValue(response.Name)
	in.Prompt = types.StringPointerValue(response.Prompt)

	if response.Webhook != nil {
		in.WebhookID = types.StringValue(response.Webhook.ID)
	} else {
		in.WebhookID = types.StringNull()
	}

	if response.IssueWebhook != nil {
		in.IssueWebhookID = types.StringValue(response.IssueWebhook.ID)
	} else {
		in.IssueWebhookID = types.StringNull()
	}

	if response.Matches != nil {
		if in.Matches == nil {
			in.Matches = &WorkbenchWebhookMatchConfig{}
		}

		in.Matches.Regex = types.StringPointerValue(response.Matches.Regex)
		in.Matches.Substring = types.StringPointerValue(response.Matches.Substring)
		in.Matches.CaseInsensitive = types.BoolPointerValue(response.Matches.CaseInsensitive)
	} else {
		in.Matches = nil
	}
}
