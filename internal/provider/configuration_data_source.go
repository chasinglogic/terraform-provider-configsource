package provider

import (
	"context"
	"fmt"

	"github.com/config-source/terraform-provider-cdb/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ConfigurationDataSource{}

func NewConfigurationDataSource() datasource.DataSource {
	return &ConfigurationDataSource{}
}

// ConfigurationDataSource defines the data source implementation.
type ConfigurationDataSource struct {
	client *client.Client
}

// ConfigurationDataSourceModel describes the data source data model.
type ConfigurationDataSourceModel struct {
	Environment types.String `tfsdk:"environment"`
	Key         types.String `tfsdk:"key"`

	StrValue   types.String  `tfsdk:"str_value"`
	IntValue   types.Int64   `tfsdk:"int_value"`
	FloatValue types.Float64 `tfsdk:"float_value"`
	BoolValue  types.Bool    `tfsdk:"bool_value"`
	Id         types.String  `tfsdk:"id"`
}

func (d *ConfigurationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (d *ConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Configuration data source, used for retrieving a key for a given environment",

		Attributes: map[string]schema.Attribute{
			"key": schema.StringAttribute{
				MarkdownDescription: "Configuration key to retrieve for environment",
				Optional:            true,
			},
			"environment": schema.StringAttribute{
				MarkdownDescription: "Name of the environment to get the configuration for",
				Optional:            true,
			},
			"str_value": schema.StringAttribute{
				MarkdownDescription: "String value stored at key if key is an string type",
				Optional:            true,
			},
			"int_value": schema.Int64Attribute{
				MarkdownDescription: "Integer value stored at key if key is an integer type",
				Optional:            true,
			},
			"float_value": schema.Float64Attribute{
				MarkdownDescription: "Float value stored at key if key is a float type",
				Optional:            true,
			},
			"bool_value": schema.BoolAttribute{
				MarkdownDescription: "Boolean value stored at key if key is a boolean type",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this Configuration Value",
				Computed:            true,
			},
		},
	}
}

func (d *ConfigurationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConfigurationDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	cv, err := d.client.GetConfigValue(ctx, data.Environment.ValueString(), data.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	data.Id = types.StringValue(fmt.Sprintf("config_value_%s_%s", data.Environment.ValueString(), data.Key.ValueString()))
	data.Key = types.StringValue(cv.Key)
	switch cv.ValueType {
	case "string":
		data.StrValue = types.StringValue(cv.StrValue)
	case "integer":
		data.IntValue = types.Int64Value(cv.IntValue)
	case "float":
		data.FloatValue = types.Float64Value(cv.FloatValue)
	case "boolean":
		data.BoolValue = types.BoolValue(cv.BoolValue)
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
