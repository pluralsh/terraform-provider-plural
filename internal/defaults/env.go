package defaults

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Env[T any, R any](envVar string, defaultValue T) R {
	return envDefaultValue{
		envVar: envVar,
		defaultValue: defaultValue,
	}
}

type envDefaultValue struct {
	envVar       string
	defaultValue string
}

func (d envDefaultValue) Description(_ context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to a string representation of the provided env variable")
}

func (d envDefaultValue) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to a string representation of the provided env variable")
}

func (d envDefaultValue) DefaultString(_ context.Context, req defaults.StringRequest, resp *defaults.StringResponse) {
	resp.PlanValue = types.StringValue(os.Getenv(d.envVar))
}
