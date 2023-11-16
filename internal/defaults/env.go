package defaults

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Env[T string | bool](envVar string, defaultValue T) EnvDefaultValue[T] {
	return EnvDefaultValue[T]{
		envVar:       envVar,
		defaultValue: defaultValue,
	}
}

type EnvDefaultValue[T string | bool] struct {
	envVar       string
	defaultValue T
}

func (d EnvDefaultValue[_]) Description(_ context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to a representation of the provided env variable")
}

func (d EnvDefaultValue[_]) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to a representation of the provided env variable")
}

func (d EnvDefaultValue[T]) DefaultString(_ context.Context, _ defaults.StringRequest, resp *defaults.StringResponse) {
	value := interface{}(d.defaultValue)
	if v := os.Getenv(d.envVar); len(v) > 0 {
		value = v
	}

	resp.PlanValue = types.StringValue(value.(string))
}

func (d EnvDefaultValue[T]) DefaultBool(_ context.Context, _ defaults.BoolRequest, resp *defaults.BoolResponse) {
	value := interface{}(d.defaultValue)
	if v := os.Getenv(d.envVar); len(v) > 0 {
		value = v == "true"
	}

	resp.PlanValue = types.BoolValue(value.(bool))
}
