package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GitRepository describes the GitRepository resource and data source model.
type GitRepository struct {
	Id         types.String `tfsdk:"id"`
	Url        types.String `tfsdk:"url"`
	PrivateKey types.String `tfsdk:"private_key"`
	Passphrase types.String `tfsdk:"passphrase"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	UrlFormat  types.String `tfsdk:"url_format"`
	HttpsPath  types.String `tfsdk:"https_path"`
}
