package validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pluralsh/polly/algorithms"
)

var _ validator.String = alsoRequiresIfValidator{}

// RequiresIf returns true when destination path should be required, false otherwise.
type RequiresIf func(source path.Path, sourceValue attr.Value, destination path.Path) bool

var (
	RequiresIfSourceValueOneOf = func(arr []string) RequiresIf {
		return func(_ path.Path, sourceValue attr.Value, _ path.Path) bool {
			return algorithms.Index(arr, func(value string) bool {
				return sourceValue.Equal(types.StringValue(value))
			}) > -1
		}
	}
)

// alsoRequiresIfValidator validates that the value matches one of expected values.
type alsoRequiresIfValidator struct {
	PathExpressions path.Expressions
	f               RequiresIf
}

type alsoRequiresValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type alsoRequiresValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

func (a alsoRequiresIfValidator) Description(ctx context.Context) string {
	return a.MarkdownDescription(ctx)
}

func (a alsoRequiresIfValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensure that if an attribute is set, these might require to be set: %q", a.PathExpressions)
}

func (a alsoRequiresIfValidator) Validate(ctx context.Context, req alsoRequiresValidatorRequest, res *alsoRequiresValidatorResponse) {
	// If attribute configuration is null, there is nothing else to validate
	if req.ConfigValue.IsNull() {
		return
	}

	expressions := req.PathExpression.MergeExpressions(a.PathExpressions...)

	for _, expression := range expressions {
		matchedPaths, diags := req.Config.PathMatches(ctx, expression)

		res.Diagnostics.Append(diags...)

		// Collect all errors
		if diags.HasError() {
			continue
		}

		for _, mp := range matchedPaths {
			// If the user specifies the same attribute this validator is applied to,
			// also as part of the input, skip it
			if mp.Equal(req.Path) {
				continue
			}

			var mpVal attr.Value
			diags := req.Config.GetAttribute(ctx, mp, &mpVal)
			res.Diagnostics.Append(diags...)

			// Collect all errors
			if diags.HasError() {
				continue
			}

			// Delay validation until all involved attribute have a known value
			if mpVal.IsUnknown() {
				return
			}

			if mpVal.IsNull() && a.f(req.Path, req.ConfigValue, mp) {
				res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
					req.Path,
					fmt.Sprintf("Attribute %q must be specified when %q with value %s is specified", mp, req.Path, req.ConfigValue),
				))
			}
		}
	}
}

func (a alsoRequiresIfValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := alsoRequiresValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &alsoRequiresValidatorResponse{}

	a.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

// AlsoRequiresIf todo
func AlsoRequiresIf(f RequiresIf, expressions ...path.Expression) validator.String {
	return &alsoRequiresIfValidator{
		PathExpressions: expressions,
		f:               f,
	}
}
