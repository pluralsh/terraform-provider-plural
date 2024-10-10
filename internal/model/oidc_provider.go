package model

import (
	"context"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/pluralsh/polly/algorithms"
	"github.com/samber/lo"
)

type OIDCProvider struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	Description  types.String `tfsdk:"description"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	AuthMethod   types.String `tfsdk:"auth_method"`
	RedirectURIs types.Set    `tfsdk:"redirect_uris"`
}

func (p *OIDCProvider) Attributes(ctx context.Context, d diag.Diagnostics) gqlclient.OidcProviderAttributes {
	return gqlclient.OidcProviderAttributes{
		Name:         p.Name.ValueString(),
		Description:  p.Description.ValueStringPointer(),
		AuthMethod:   p.authMethodAttribute(),
		RedirectUris: p.redirectURIsAttribute(ctx, d),
	}
}

func (p *OIDCProvider) authMethodAttribute() *gqlclient.OidcAuthMethod {
	if p.AuthMethod.IsNull() {
		return nil
	}

	return lo.ToPtr(gqlclient.OidcAuthMethod(p.AuthMethod.ValueString()))
}

func (p *OIDCProvider) redirectURIsAttribute(ctx context.Context, d diag.Diagnostics) []*string {
	redirectURIs := make([]types.String, len(p.RedirectURIs.Elements()))
	d.Append(p.RedirectURIs.ElementsAs(ctx, &redirectURIs, false)...)
	return algorithms.Map(redirectURIs, func(v types.String) *string { return v.ValueStringPointer() })
}

func (p *OIDCProvider) TypeAttribute() gqlclient.OidcProviderType {
	return gqlclient.OidcProviderType(p.Type.ValueString())
}

func (p *OIDCProvider) From(response *gqlclient.OIDCProviderFragment, ctx context.Context, d diag.Diagnostics) {
	p.ID = types.StringValue(response.ID)
	p.Name = types.StringValue(response.Name)
	p.Description = types.StringPointerValue(response.Description)
	p.ClientID = types.StringValue(response.ClientID)
	p.ClientSecret = types.StringValue(response.ClientSecret)
	p.AuthMethod = p.authMethodFrom(response.AuthMethod)
	p.RedirectURIs = common.SetFrom(response.RedirectUris, ctx, d)

}

func (p *OIDCProvider) authMethodFrom(authMethod *gqlclient.OidcAuthMethod) types.String {
	if authMethod == nil {
		return types.StringNull()
	}

	return types.StringValue(string(*authMethod))
}
