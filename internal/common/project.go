package common

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console-client-go"
)

func ProjectFrom(project *gqlclient.TinyProjectFragment) types.String {
	if project != nil {
		return types.StringValue(project.ID)
	}

	return types.StringNull()
}
