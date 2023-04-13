package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	fz "github.com/jindrichskupa/firezone-client-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DeviceResource{}
var _ resource.ResourceWithImportState = &DeviceResource{}

func NewDeviceResource() resource.Resource {
	return &DeviceResource{}
}

// DeviceResource defines the resource implementation.
type DeviceResource struct {
	client *fz.Client
}

// DeviceResourceModel describes the resource data model.
type DeviceResourceModel struct {
	Id                            types.String `tfsdk:"id"`
	AllowedIPs                    []string     `tfsdk:"allowed_ips"`
	Description                   types.String `tfsdk:"description"`
	DNS                           []string     `tfsdk:"dns"`
	Endpoint                      types.String `tfsdk:"endpoint"`
	InsertedAt                    types.String `tfsdk:"inserted_at"`
	IPv4                          types.String `tfsdk:"ipv4"`
	IPv6                          types.String `tfsdk:"ipv6"`
	LatestHandshake               types.String `tfsdk:"latest_handshake"`
	MTU                           types.Int64  `tfsdk:"mtu"`
	Name                          types.String `tfsdk:"name"`
	PersistentKeepalive           types.Int64  `tfsdk:"persistent_keepalive"`
	PresharedKey                  types.String `tfsdk:"preshared_key"`
	PublicKey                     types.String `tfsdk:"public_key"`
	RemoteIP                      types.String `tfsdk:"remote_ip"`
	RXBytes                       interface{}  `tfsdk:"rx_bytes"`
	ServerPublicKey               types.String `tfsdk:"server_public_key"`
	TXBytes                       interface{}  `tfsdk:"tx_bytes"`
	UpdatedAt                     types.String `tfsdk:"updated_at,omitempty"`
	UseDefaultAllowedIPs          types.Bool   `tfsdk:"use_default_allowed_ips"`
	UseDefaultDNS                 types.Bool   `tfsdk:"use_default_dns"`
	UseDefaultEndpoint            types.Bool   `tfsdk:"use_default_endpoint"`
	UseDefaultMTU                 types.Bool   `tfsdk:"use_default_mtu"`
	UseDefaultPersistentKeepalive types.Bool   `tfsdk:"use_default_persistent_keepalive"`
	UserId                        types.String `tfsdk:"user_id"`
}

func (r *DeviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (r *DeviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Device resource",

		Attributes: map[string]schema.Attribute{
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Device updated at",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"inserted_at": schema.StringAttribute{
				MarkdownDescription: "Device inserted at",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},

			"user_id": schema.StringAttribute{
				MarkdownDescription: "Device user id",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Device identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *DeviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*fz.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DeviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DeviceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	device, err := r.client.CreateDevice(fz.Device{})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create device, got error: %s", err))
		return
	}

	data.Id = types.StringValue(device.ID)
	data.UserId = types.StringValue(device.UserId)

	data.UpdatedAt = types.StringValue(device.UpdatedAt)
	data.InsertedAt = types.StringValue(device.InsertedAt)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DeviceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	device, err := r.client.GetDevice(data.Id.String())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read device, got error: %s", err))
		return
	}

	data.Id = types.StringValue(device.ID)
	data.UserId = types.StringValue(device.UserId)

	data.UpdatedAt = types.StringValue(device.UpdatedAt)
	data.InsertedAt = types.StringValue(device.InsertedAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DeviceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	device, err := r.client.UpdateDevice(data.Id.String(), fz.Device{})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update device, got error: %s", err))
		return
	}

	data.Id = types.StringValue(device.ID)
	data.UserId = types.StringValue(device.UserId)

	data.UpdatedAt = types.StringValue(device.UpdatedAt)
	data.InsertedAt = types.StringValue(device.InsertedAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DeviceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDevice(data.Id.String())

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete device, got error: %s", err))
		return
	}
}

func (r *DeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
