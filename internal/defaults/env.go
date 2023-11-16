package defaults

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnvDefaultValue interface {
	defaults.String
	defaults.Bool
}

type defaultable interface {
	string | bool
}

func Env[T defaultable](envVar string, defaultValue T) EnvDefaultValue {
	return envDefaultValue[T]{
		envVar:       envVar,
		defaultValue: defaultValue,
	}
}

type envDefaultValue[T defaultable] struct {
	envVar       string
	defaultValue T
}

func (d envDefaultValue[_]) Description(_ context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to a representation of the provided env variable")
}

func (d envDefaultValue[_]) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to a representation of the provided env variable")
}

func (d envDefaultValue[T]) DefaultString(_ context.Context, _ defaults.StringRequest, resp *defaults.StringResponse) {
	value := interface{}(d.defaultValue)
	if v := os.Getenv(d.envVar); len(v) > 0 {
		value = v
	}

	resp.PlanValue = types.StringValue(value.(string))
}

func (d envDefaultValue[T]) DefaultBool(_ context.Context, _ defaults.BoolRequest, resp *defaults.BoolResponse) {
	value := interface{}(d.defaultValue)
	if v := os.Getenv(d.envVar); len(v) > 0 {
		value = v == "true"
	}

	resp.PlanValue = types.BoolValue(value.(bool))
}
