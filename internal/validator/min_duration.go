package validator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = minDurationValidator{}

// minDurationValidator validates that a string is a valid time.Duration and meets a minimum value.
type minDurationValidator struct {
	minDuration time.Duration
}

func (v minDurationValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Value must be a valid duration string and at least %s.", v.minDuration)
}

func (v minDurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v minDurationValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	duration, err := time.ParseDuration(value)
	if err != nil {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Duration",
			fmt.Sprintf("Value %q is not a valid duration string: %s. Valid examples: '5m', '1h30m', '500ms'", value, err.Error()),
		)
		return
	}

	if duration < v.minDuration {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Duration Too Short",
			fmt.Sprintf("Value %q (%s) is less than the minimum allowed duration of %s", value, duration, v.minDuration),
		)
	}
}

// MinDuration returns a validator which ensures that the configured string value
// is a valid time.Duration and is at least the specified minimum duration.
func MinDuration(minDuration time.Duration) validator.String {
	return minDurationValidator{
		minDuration: minDuration,
	}
}
