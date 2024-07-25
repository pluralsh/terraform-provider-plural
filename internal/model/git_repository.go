package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
)

// TODO: Embed structs like these once it will be supported: https://github.com/hashicorp/terraform-plugin-framework/pull/1021

type GitRepository struct {
	Id  types.String `tfsdk:"id"`
	Url types.String `tfsdk:"url"`
}

func (gr *GitRepository) From(response *gqlclient.GitRepositoryFragment) {
	gr.Id = types.StringValue(response.ID)
	gr.Url = types.StringValue(response.URL)
}

type GitRepositoryExtended struct {
	Id         types.String `tfsdk:"id"`
	Url        types.String `tfsdk:"url"`
	PrivateKey types.String `tfsdk:"private_key"`
	Passphrase types.String `tfsdk:"passphrase"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	UrlFormat  types.String `tfsdk:"url_format"`
	HttpsPath  types.String `tfsdk:"https_path"`
	Decrypt    types.Bool   `tfsdk:"decrypt"`
}

func (gre *GitRepositoryExtended) From(response *gqlclient.GitRepositoryFragment) {
	gre.Id = types.StringValue(response.ID)
	gre.Url = types.StringValue(response.URL)
}

func (gre *GitRepositoryExtended) Attributes() gqlclient.GitAttributes {
	return gqlclient.GitAttributes{
		URL:        gre.Url.ValueString(),
		PrivateKey: gre.PrivateKey.ValueStringPointer(),
		Passphrase: gre.Passphrase.ValueStringPointer(),
		Username:   gre.Username.ValueStringPointer(),
		Password:   gre.Password.ValueStringPointer(),
		HTTPSPath:  gre.HttpsPath.ValueStringPointer(),
		URLFormat:  gre.UrlFormat.ValueStringPointer(),
		Decrypt:    gre.Decrypt.ValueBoolPointer(),
	}
}
