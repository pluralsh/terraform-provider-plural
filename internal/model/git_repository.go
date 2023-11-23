package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console-client-go"
)

type GitRepositoryBase struct {
	Id         types.String `tfsdk:"id"`
	Url        types.String `tfsdk:"url"`
}

func (this *GitRepositoryBase) From(response *gqlclient.GitRepositoryFragment) {
	this.Id = types.StringValue(response.ID)
	this.Url = types.StringValue(response.URL)
}

// GitRepository describes the Git repository resource and data source model.
type GitRepository struct {
	GitRepositoryBase

	PrivateKey types.String `tfsdk:"private_key"`
	Passphrase types.String `tfsdk:"passphrase"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	UrlFormat  types.String `tfsdk:"url_format"`
	HttpsPath  types.String `tfsdk:"https_path"`
}

func (this *GitRepository) Attributes() gqlclient.GitAttributes {
	return gqlclient.GitAttributes{
		URL:        this.Url.ValueString(),
		PrivateKey: this.PrivateKey.ValueStringPointer(),
		Passphrase: this.Passphrase.ValueStringPointer(),
		Username:   this.Username.ValueStringPointer(),
		Password:   this.Password.ValueStringPointer(),
		HTTPSPath:  this.HttpsPath.ValueStringPointer(),
		URLFormat:  this.UrlFormat.ValueStringPointer(),
	}
}
