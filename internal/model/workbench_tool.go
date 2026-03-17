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
	Id            types.String                `tfsdk:"id"`
	Name          types.String                `tfsdk:"name"`
	Tool          types.String                `tfsdk:"tool"`
	Categories    types.Set                   `tfsdk:"categories"`
	ProjectID     types.String                `tfsdk:"project_id"`
	Configuration *WorkbenchToolConfiguration `tfsdk:"configuration"`
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
		ProjectID:     in.ProjectID.ValueStringPointer(),
		Configuration: in.Configuration.Attributes(ctx),
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

	in.Configuration.From(response.Configuration, ctx, d)
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

func (in *WorkbenchToolConfiguration) From(configuration *gqlclient.WorkbenchToolFragment_Configuration, ctx context.Context, d *diag.Diagnostics) {
	if in == nil {
		return
	}
	if configuration == nil {
		return
	}

	in.HTTP.From(configuration.HTTP, ctx, d)
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
	}

	if configuration.Headers != nil {
		headers := make(map[string]any, len(configuration.Headers))
		for _, v := range configuration.Headers {
			if v.Value != nil {
				headers[*v.Name] = *v.Value
			}
		}

		in.Headers = common.MapFromWithConfig(headers, in.Headers, ctx, d)
	}

	if configuration.Body != nil {
		in.Body = types.StringPointerValue(configuration.Body)
	}

	if configuration.InputSchema != nil {
		inputSchema, err := json.Marshal(configuration.InputSchema)
		if err != nil {
			d.AddError("Provider Error", fmt.Sprintf("Cannot marshall input schema, got error: %s", err))
			return
		}

		in.InputSchema = types.StringValue(string(inputSchema))
	}
}
