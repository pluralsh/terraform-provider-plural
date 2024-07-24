package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
)

type gitRepository struct {
	Id  types.String `tfsdk:"id"`
	Url types.String `tfsdk:"url"`
}

func (g *gitRepository) From(response *gqlclient.GitRepositoryFragment) {
	g.Id = types.StringValue(response.ID)
	g.Url = types.StringValue(response.URL)
}
