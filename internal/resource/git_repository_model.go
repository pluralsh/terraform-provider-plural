package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console-client-go"
)

type gitRepository struct {
	Id         types.String `tfsdk:"id"`
	Url        types.String `tfsdk:"url"`
	PrivateKey types.String `tfsdk:"private_key"`
	Passphrase types.String `tfsdk:"passphrase"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	UrlFormat  types.String `tfsdk:"url_format"`
	HttpsPath  types.String `tfsdk:"https_path"`
}

func (g *gitRepository) From(response *gqlclient.GitRepositoryFragment) {
	g.Id = types.StringValue(response.ID)
	g.Url = types.StringValue(response.URL)
}

func (g *gitRepository) Attributes() gqlclient.GitAttributes {
	return gqlclient.GitAttributes{
		URL:        g.Url.ValueString(),
		PrivateKey: g.PrivateKey.ValueStringPointer(),
		Passphrase: g.Passphrase.ValueStringPointer(),
		Username:   g.Username.ValueStringPointer(),
		Password:   g.Password.ValueStringPointer(),
		HTTPSPath:  g.HttpsPath.ValueStringPointer(),
		URLFormat:  g.UrlFormat.ValueStringPointer(),
	}
}
