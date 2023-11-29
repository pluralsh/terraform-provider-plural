package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console-client-go"
)

type GitRepository struct {
	Id  types.String `tfsdk:"id"`
	Url types.String `tfsdk:"url"`
}

func (g *GitRepository) From(response *gqlclient.GitRepositoryFragment) {
	g.Id = types.StringValue(response.ID)
	g.Url = types.StringValue(response.URL)
}
