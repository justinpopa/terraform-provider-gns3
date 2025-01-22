package provider

import (
	"context"

	"client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &gns3Provider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &gns3Provider{
			version: version,
		}
	}
}

// hashicupsProvider is the provider implementation.
type gns3Provider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type gns3ProviderModel struct {
	Host types.String         `tfsdk:"host"`
	Port basetypes.Int32Value `tfsdk:"port"`
}

// Metadata returns the provider type name.
func (p *gns3Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "gns3"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *gns3Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    false,
				Description: "The host of the GNS3 server.",
			},
			"port": schema.Int32Attribute{
				Optional:    false,
				Description: "The port for the GNS3 server.",
			},
		},
	}
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *gns3Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config gns3ProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown GNS3 Host",
			"The provider cannot create the GNS3 client as there is an unknown configuration value for the GNS3 server host. "+
				"Set the value statically in the provider configuration.",
		)
	}

	if config.Port.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("port"),
			"Unknown GNS3 Port",
			"The provider cannot create the GNS3 client as there is an unknown configuration value for the GNS3 server port. "+
				"Set the value statically in the provider configuration.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := client.NewClient(config.Host.Value, int(config.Port.Value))

	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *gns3Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *gns3Provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
