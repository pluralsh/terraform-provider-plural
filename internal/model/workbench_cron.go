package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
)

type WorkbenchCron struct {
	Id          types.String `tfsdk:"id"`
	WorkbenchID types.String `tfsdk:"workbench_id"`
	Crontab     types.String `tfsdk:"crontab"`
	Prompt      types.String `tfsdk:"prompt"`
}

func (in *WorkbenchCron) Attributes() gqlclient.WorkbenchCronAttributes {
	return gqlclient.WorkbenchCronAttributes{
		Crontab: in.Crontab.ValueStringPointer(),
		Prompt:  in.Prompt.ValueStringPointer(),
	}
}

func (in *WorkbenchCron) From(response *gqlclient.WorkbenchCronFragment) {
	if in == nil || response == nil {
		return
	}

	in.Id = types.StringValue(response.ID)
	in.Crontab = types.StringPointerValue(response.Crontab)
	in.Prompt = types.StringPointerValue(response.Prompt)
}
