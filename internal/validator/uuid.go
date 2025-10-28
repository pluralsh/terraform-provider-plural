package validator

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = uuidValidator{}

// uuidValidator validates that a string is a valid UUID (v4 format).
type uuidValidator struct{}

// UUID regex pattern that matches standard UUID format (8-4-4-4-12 hex digits)
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func (v uuidValidator) Description(_ context.Context) string {
	return "Value must be a valid UUID."
}

func (v uuidValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v uuidValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if !uuidRegex.MatchString(value) {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid UUID",
			fmt.Sprintf("Value %q is not a valid UUID. Expected format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", value),
		)
	}
}

// UUID returns a validator which ensures that the configured string value
// is a valid UUID.
func UUID() validator.String {
	return uuidValidator{}
}
