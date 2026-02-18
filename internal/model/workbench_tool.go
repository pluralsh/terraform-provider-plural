package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

type WorkbenchTool struct {
	Id            types.String                `tfsdk:"id"`
	Name          types.String                `tfsdk:"name"`
	Tool          types.String                `tfsdk:"tool"`
	Categories    types.Set                   `tfsdk:"categories"`
	ProjectID     types.String                `tfsdk:"project_id"`
	Configuration *WorkbenchToolConfiguration `tfsdk:"configuration"`
}

func (wt *WorkbenchTool) Attributes(ctx context.Context, d *diag.Diagnostics) (*gqlclient.WorkbenchToolAttributes, error) {
	categories := make([]types.String, len(wt.Categories.Elements()))
	wt.Categories.ElementsAs(ctx, &categories, false)

	return &gqlclient.WorkbenchToolAttributes{
		Name: wt.Name.ValueString(),
		Tool: gqlclient.WorkbenchToolType(wt.Tool.ValueString()),
		Categories: lo.Map(categories, func(v types.String, _ int) *gqlclient.WorkbenchToolCategory {
			return lo.ToPtr(gqlclient.WorkbenchToolCategory(v.ValueString()))
		}),
		ProjectID:     wt.ProjectID.ValueStringPointer(),
		Configuration: wt.Configuration.Attributes(ctx),
	}, nil
}

type WorkbenchToolConfiguration struct {
	HTTP *WorkbenchToolHTTPConfig `tfsdk:"http"`
}

func (in *WorkbenchToolConfiguration) Attributes(ctx context.Context) *gqlclient.WorkbenchToolConfigurationAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.WorkbenchToolConfigurationAttributes{
		HTTP: in.HTTP.Attributes(ctx),
	}
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
		Method: gqlclient.WorkbenchToolHTTPMethod(in.Method.ValueString()),
		Headers: lo.MapToSlice(headers, func(k string, v types.String) *gqlclient.WorkbenchToolHTTPHeaderAttributes {
			return &gqlclient.WorkbenchToolHTTPHeaderAttributes{Name: &k, Value: v.ValueStringPointer()}
		}),
		Body:        in.Body.ValueStringPointer(),
		InputSchema: in.InputSchema.ValueStringPointer(),
	}
}
