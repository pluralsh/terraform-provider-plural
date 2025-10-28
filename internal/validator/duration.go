package validator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = durationValidator{}

// durationValidator validates that a string is a valid time.Duration.
type durationValidator struct{}

func (v durationValidator) Description(_ context.Context) string {
	return "Value must be a valid duration string (e.g., '5m', '1h30m', '500ms')."
}

func (v durationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v durationValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if _, err := time.ParseDuration(value); err != nil {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Duration",
			fmt.Sprintf("Value %q is not a valid duration string: %s. Valid examples: '5m', '1h30m', '500ms'", value, err.Error()),
		)
	}
}

// Duration returns a validator which ensures that the configured string value
// is a valid time.Duration.
func Duration() validator.String {
	return durationValidator{}
}
