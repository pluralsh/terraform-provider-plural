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

var _ validator.String = conflictsWithIfValidator{}

// ConflictsIf returns true when destination path should not be provided, false otherwise.
type ConflictsIf func(source path.Path, sourceValue attr.Value, target path.Path, targetValue attr.Value) bool

var (
	ConflictsIfTargetValueOneOf = func(arr []string) ConflictsIf {
		return func(_ path.Path, _ attr.Value, _ path.Path, targetValue attr.Value) bool {
			return algorithms.Index(arr, func(value string) bool {
				return targetValue.Equal(types.StringValue(value))
			}) > -1
		}
	}
)

// conflictsWithIfValidator validates that the value matches one of expected values.
type conflictsWithIfValidator struct {
	PathExpressions path.Expressions
	f               ConflictsIf
}

type conflictsWithValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type conflictsWithValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

func (a conflictsWithIfValidator) Description(ctx context.Context) string {
	return a.MarkdownDescription(ctx)
}

func (a conflictsWithIfValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensure that if an attribute is set, these might not be allowed to be set: %q", a.PathExpressions)
}

func (a conflictsWithIfValidator) Validate(ctx context.Context, req conflictsWithValidatorRequest, res *conflictsWithValidatorResponse) {
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

			if !mpVal.IsNull() && a.f(req.Path, req.ConfigValue, mp, mpVal) {
				res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
					req.Path,
					fmt.Sprintf("Attribute %q with value %s cannot be specified when %q is specified", mp, mpVal, req.Path),
				))
			}
		}
	}
}

func (a conflictsWithIfValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := conflictsWithValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &conflictsWithValidatorResponse{}

	a.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

// ConflictsWithIf checks that a set of path.Expression,
// including the attribute the validator is applied to,
// do not have a value simultaneously and specified condition is met.
//
// Relative path.Expression will be resolved using the attribute being
// validated.
func ConflictsWithIf(f ConflictsIf, expressions ...path.Expression) validator.String {
	return &conflictsWithIfValidator{
		PathExpressions: expressions,
		f:               f,
	}
}
