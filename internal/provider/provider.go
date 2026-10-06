// Copyright Thayne McCombs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure ExtLibProvider satisfies various provider interfaces.
var _ provider.Provider = &ExtLibProvider{}
var _ provider.ProviderWithFunctions = &ExtLibProvider{}

// ExtLibProvider defines the provider implementation.
type ExtLibProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ExtLibProvider{
			version: version,
		}
	}
}

func (p *ExtLibProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "extlib"
	resp.Version = p.version
}

func (p *ExtLibProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *ExtLibProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *ExtLibProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *ExtLibProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *ExtLibProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewBoolnumFunction,
		NewBoolsetFunction,
	}
}
