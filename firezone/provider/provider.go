package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	fz "github.com/jindrichskupa/firezone-client-go/client"
)

// Ensure FirezoneProvider satisfies various provider interfaces.
var _ provider.Provider = &FirezoneProvider{}

// FirezoneProvider defines the provider implementation.
type FirezoneProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// FirezoneProviderModel describes the provider data model.
type FirezoneProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiKey   types.String `tfsdk:"api_key"`
}

func (p *FirezoneProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "firezone"
	resp.Version = p.version
}

func (p *FirezoneProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Firezone API endpoint",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Firezone API key",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *FirezoneProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data FirezoneProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("FIREZONE_ENDPOINT")
	api_key := os.Getenv("FIREZONE_API_KEY")

	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}

	if !data.ApiKey.IsNull() {
		api_key = data.ApiKey.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing firezone endpoint",
			"Alternative: FIREZONE_ENDPOINT environment variable",
		)
	}

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing firezone api key",
			"Alternative: FIREZONE_API_KEY environment variable",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := fz.NewClient(endpoint, api_key)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create firezone API Client",
			"An unexpected error occurred "+err.Error(),
		)
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *FirezoneProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
		NewRuleResource,
	}
}

func (p *FirezoneProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUserDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FirezoneProvider{
			version: version,
		}
	}
}
