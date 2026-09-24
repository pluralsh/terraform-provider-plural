package resource

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

// enumValues converts the values of a generated Console API enum, e.g. gqlclient.AllMonitorType,
// to strings, so validators stay in sync with the API.
func enumValues[T ~string](values []T) []string {
	return lo.Map(values, func(v T, _ int) string { return string(v) })
}

// withAllowedValues appends the allowed values to an attribute description.
func withAllowedValues(description string, values []string, markdown bool) string {
	format := `"%s"`
	if markdown {
		format = "`%s`"
	}

	quoted := lo.Map(values, func(v string, _ int) string { return fmt.Sprintf(format, v) })
	if len(quoted) > 1 {
		quoted = append(quoted[:len(quoted)-2], quoted[len(quoted)-2]+" and "+quoted[len(quoted)-1])
	}

	return fmt.Sprintf("%s Allowed values include %s.", description, strings.Join(quoted, ", "))
}
