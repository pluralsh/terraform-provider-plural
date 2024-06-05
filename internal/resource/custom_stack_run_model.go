package resource

import (
	"context"

	"terraform-provider-plural/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	gqlclient "github.com/pluralsh/console-client-go"
	"github.com/pluralsh/polly/algorithms"
)

type customStackRun struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Documentation types.String `tfsdk:"documentation"`
	StackId       types.String `tfsdk:"stack_id"`
	Commands      types.Set    `tfsdk:"commands"`
	Configuration types.Set    `tfsdk:"configuration"`
}

func (csr *customStackRun) Attributes(ctx context.Context, d diag.Diagnostics, client *client.Client) (*gqlclient.CustomStackRunAttributes, error) {
	return &gqlclient.CustomStackRunAttributes{
		Name:          csr.Name.ValueString(),
		Documentation: csr.Documentation.ValueStringPointer(),
		StackID:       csr.StackId.ValueStringPointer(),
		Commands:      csr.commandsAttribute(ctx, d),
		Configuration: csr.configurationAttribute(ctx, d),
	}, nil
}

func (csr *customStackRun) commandsAttribute(ctx context.Context, d diag.Diagnostics) []*gqlclient.CommandAttributes {
	if csr.Commands.IsNull() {
		return nil
	}

	result := make([]*gqlclient.CommandAttributes, 0, len(csr.Commands.Elements()))
	elements := make([]CustomStackRunCommand, len(csr.Commands.Elements()))
	d.Append(csr.Commands.ElementsAs(ctx, &elements, false)...)

	for _, cmd := range elements {
		args := make([]types.String, len(cmd.Args.Elements()))
		d.Append(cmd.Args.ElementsAs(ctx, &args, false)...)

		result = append(result, &gqlclient.CommandAttributes{
			Cmd:  cmd.Cmd.ValueString(),
			Args: algorithms.Map(args, func(v types.String) *string { return v.ValueStringPointer() }),
			Dir:  cmd.Dir.ValueStringPointer(),
		})
	}

	return result
}

func (csr *customStackRun) configurationAttribute(ctx context.Context, d diag.Diagnostics) []*gqlclient.PrConfigurationAttributes {
	if csr.Configuration.IsNull() {
		return nil
	}

	result := make([]*gqlclient.PrConfigurationAttributes, 0, len(csr.Commands.Elements()))
	elements := make([]CustomStackRunConfiguration, len(csr.Commands.Elements()))
	d.Append(csr.Commands.ElementsAs(ctx, &elements, false)...)

	for _, cfg := range elements {
		result = append(result, &gqlclient.PrConfigurationAttributes{
			Type:          gqlclient.ConfigurationType(cfg.Type.ValueString()),
			Name:          cfg.Name.ValueString(),
			Default:       cfg.Default.ValueStringPointer(),
			Documentation: cfg.Documentation.ValueStringPointer(),
			Longform:      cfg.Longform.ValueStringPointer(),
			Placeholder:   cfg.Placeholder.ValueStringPointer(),
			Optional:      cfg.Optional.ValueBoolPointer(),
			Condition:     cfg.Condition.Attributes(),
		})
	}

	return result
}

func (csr *customStackRun) From(customStackRun *gqlclient.CustomStackRunFragment, ctx context.Context, d diag.Diagnostics) {
	csr.Id = types.StringValue(customStackRun.ID)
	csr.Name = types.StringValue(customStackRun.Name)
	csr.Documentation = types.StringPointerValue(customStackRun.Documentation)
	csr.StackId = types.StringPointerValue(customStackRun.Stack.ID)
	csr.Commands = commandsFrom(customStackRun.Commands, csr.Commands, ctx, d)
	csr.Configuration = configurationFrom(customStackRun.Configuration, csr.Configuration, ctx, d)
}

type CustomStackRunCommand struct {
	Cmd  types.String `tfsdk:"cmd"`
	Args types.Set    `tfsdk:"args"`
	Dir  types.String `tfsdk:"dir"`
}

var CustomStackRunCommandAttrTypes = map[string]attr.Type{
	"cmd":  types.StringType,
	"args": types.SetType{ElemType: types.StringType},
	"dir":  types.StringType,
}

func commandsFrom(commands []*gqlclient.StackCommandFragment, config types.Set, ctx context.Context, d diag.Diagnostics) types.Set {
	if len(commands) == 0 {
		// Rewriting config to state to avoid inconsistent result errors.
		// This could happen, for example, when sending "nil" to API and "[]" is returned as a result.
		return config
	}

	values := make([]attr.Value, len(commands))
	for i, command := range commands {
		objValue, diags := types.ObjectValueFrom(ctx, CustomStackRunCommandAttrTypes, CustomStackRunCommand{
			Cmd:  types.StringValue(command.Cmd),
			Args: commandArgsFrom(command.Args, ctx, d),
			Dir:  types.StringPointerValue(command.Dir),
		})
		values[i] = objValue
		d.Append(diags...)
	}

	setValue, diags := types.SetValue(basetypes.ObjectType{AttrTypes: CustomStackRunCommandAttrTypes}, values)
	d.Append(diags...)
	return setValue
}

func commandArgsFrom(values []*string, ctx context.Context, d diag.Diagnostics) types.Set {
	if values == nil {
		return types.SetNull(types.StringType)
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, values)
	d.Append(diags...)
	return setValue
}

type CustomStackRunConfiguration struct {
	Type          types.String                          `tfsdk:"type"`
	Name          types.String                          `tfsdk:"name"`
	Default       types.String                          `tfsdk:"default"`
	Documentation types.String                          `tfsdk:"documentation"`
	Longform      types.String                          `tfsdk:"longform"`
	Placeholder   types.String                          `tfsdk:"placeholder"`
	Optional      types.Bool                            `tfsdk:"optional"`
	Condition     *CustomStackRunConfigurationCondition `tfsdk:"condition"`
}

var CustomStackRunConfigurationAttrTypes = map[string]attr.Type{
	"type":          types.StringType,
	"name":          types.StringType,
	"default":       types.StringType,
	"documentation": types.StringType,
	"longform":      types.StringType,
	"placeholder":   types.StringType,
	"optional":      types.BoolType,
	"condition":     types.ObjectType{AttrTypes: CustomStackRunCommandConditionAttrTypes},
}

func configurationFrom(configs []*gqlclient.PrConfigurationFragment, config types.Set, ctx context.Context, d diag.Diagnostics) types.Set {
	if len(configs) == 0 {
		// Rewriting config to state to avoid inconsistent result errors.
		// This could happen, for example, when sending "nil" to API and "[]" is returned as a result.
		return config
	}

	values := make([]attr.Value, len(configs))
	for i, cfg := range configs {
		value := CustomStackRunConfiguration{
			Type:          types.StringValue(string(cfg.Type)),
			Name:          types.StringValue(cfg.Name),
			Default:       types.StringPointerValue(cfg.Default),
			Documentation: types.StringPointerValue(cfg.Documentation),
			Longform:      types.StringPointerValue(cfg.Longform),
			Placeholder:   types.StringPointerValue(cfg.Placeholder),
			Optional:      types.BoolPointerValue(cfg.Optional),
		}

		if cfg.Condition != nil {
			value.Condition = &CustomStackRunConfigurationCondition{
				Operation: types.StringValue(string(cfg.Condition.Operation)),
				Field:     types.StringValue(cfg.Condition.Field),
				Value:     types.StringPointerValue(cfg.Condition.Value),
			}
		}

		objValue, diags := types.ObjectValueFrom(ctx, CustomStackRunConfigurationAttrTypes, value)
		values[i] = objValue
		d.Append(diags...)
	}

	setValue, diags := types.SetValue(basetypes.ObjectType{AttrTypes: CustomStackRunConfigurationAttrTypes}, values)
	d.Append(diags...)
	return setValue
}

type CustomStackRunConfigurationCondition struct {
	Operation types.String `tfsdk:"operation"`
	Field     types.String `tfsdk:"field"`
	Value     types.String `tfsdk:"value"`
}

var CustomStackRunCommandConditionAttrTypes = map[string]attr.Type{
	"operation": types.StringType,
	"field":     types.StringType,
	"value":     types.StringType,
}

func (csrcc *CustomStackRunConfigurationCondition) Attributes() *gqlclient.ConditionAttributes {
	if csrcc == nil {
		return nil
	}

	return &gqlclient.ConditionAttributes{
		Operation: gqlclient.Operation(csrcc.Operation.ValueString()),
		Field:     csrcc.Field.ValueString(),
		Value:     csrcc.Value.ValueStringPointer(),
	}
}
