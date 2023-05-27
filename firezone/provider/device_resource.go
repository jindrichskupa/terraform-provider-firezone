package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	Id types.String `tfsdk:"id"`
	// AllowedIPs                    types.List   `tfsdk:"allowed_ips"`
	Description types.String `tfsdk:"description"`
	// DNS                           types.List   `tfsdk:"dns"`
	Endpoint                      types.String `tfsdk:"endpoint"`
	IPv4                          types.String `tfsdk:"ipv4"`
	IPv6                          types.String `tfsdk:"ipv6"`
	MTU                           types.Int64  `tfsdk:"mtu"`
	Name                          types.String `tfsdk:"name"`
	PersistentKeepalive           types.Int64  `tfsdk:"persistent_keepalive"`
	PresharedKey                  types.String `tfsdk:"preshared_key"`
	PublicKey                     types.String `tfsdk:"public_key"`
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
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Device endpoint",
				Optional:            true,
				Computed:            true,
			},
			"preshared_key": schema.StringAttribute{
				MarkdownDescription: "Device preshared key",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"mtu": schema.Int64Attribute{
				MarkdownDescription: "Device MTU",
				Optional:            true,
				Computed:            true,
			},
			"use_default_dns": schema.BoolAttribute{
				MarkdownDescription: "Device use default DNS",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"use_default_endpoint": schema.BoolAttribute{
				MarkdownDescription: "Device use default endpoint",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"use_default_mtu": schema.BoolAttribute{
				MarkdownDescription: "Device use default MTU",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"use_default_allowed_ips": schema.BoolAttribute{
				MarkdownDescription: "Device use default allowed ips",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"use_default_persistent_keepalive": schema.BoolAttribute{
				MarkdownDescription: "Device use default persistent keepalive",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"persistent_keepalive": schema.Int64Attribute{
				MarkdownDescription: "Device persistent keepalive",
				Optional:            true,
				Computed:            true,
			},
			"ipv6": schema.StringAttribute{
				MarkdownDescription: "Device IPv6",
				Optional:            true,
				Computed:            true,
			},
			"ipv4": schema.StringAttribute{
				MarkdownDescription: "Device IPv4",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Device description",
				Optional:            true,
				Computed:            true,
			},
			// "allowed_ips": schema.ListAttribute{
			// 	MarkdownDescription: "Device allowed ips",
			// 	Optional:            true,
			// },
			"public_key": schema.StringAttribute{
				MarkdownDescription: "Device public key",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Device name",
				Required:            true,
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

	device, err := r.client.CreateDevice(fz.Device{
		UserId:      data.UserId.ValueString(),
		Name:        data.Name.ValueString(),
		PublicKey:   data.PublicKey.ValueString(),
		Description: data.Description.ValueString(),
		IPv4:        data.IPv4.ValueString(),
		IPv6:        data.IPv6.ValueString(),
		// AllowedIPs:  data.AllowedIPs.ValueList(),
		Endpoint:     data.Endpoint.ValueString(),
		PresharedKey: data.PresharedKey.ValueString(),
		MTU:          int(data.MTU.ValueInt64()),
		// DNS:                           data.DNS.ValueString(),
		PersistentKeepalive:           int(data.PersistentKeepalive.ValueInt64()),
		UseDefaultDNS:                 data.UseDefaultDNS.ValueBool(),
		UseDefaultEndpoint:            data.UseDefaultEndpoint.ValueBool(),
		UseDefaultMTU:                 data.UseDefaultMTU.ValueBool(),
		UseDefaultAllowedIPs:          data.UseDefaultAllowedIPs.ValueBool(),
		UseDefaultPersistentKeepalive: data.UseDefaultPersistentKeepalive.ValueBool(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create device, got error: %s", err))
		return
	}

	data.Id = types.StringValue(device.ID)
	data.UserId = types.StringValue(device.UserId)
	data.Name = types.StringValue(device.Name)
	data.PublicKey = types.StringValue(device.PublicKey)
	data.Description = types.StringValue(device.Description)
	data.IPv4 = types.StringValue(device.IPv4)
	data.IPv6 = types.StringValue(device.IPv6)
	// data.AllowedIPs, = types.ListValue(device.AllowedIPs)
	data.Endpoint = types.StringValue(device.Endpoint)
	data.PresharedKey = types.StringValue(device.PresharedKey)
	data.MTU = types.Int64Value(int64(int64(device.MTU)))
	// data.DNS = types.StringValue(device.DNS)
	data.PersistentKeepalive = types.Int64Value(int64(device.PersistentKeepalive))
	data.UseDefaultAllowedIPs = types.BoolValue(device.UseDefaultAllowedIPs)
	data.UseDefaultDNS = types.BoolValue(device.UseDefaultDNS)
	data.UseDefaultEndpoint = types.BoolValue(device.UseDefaultEndpoint)
	data.UseDefaultMTU = types.BoolValue(device.UseDefaultMTU)
	data.UseDefaultPersistentKeepalive = types.BoolValue(device.UseDefaultPersistentKeepalive)

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

	device, err := r.client.GetDevice(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read device, got error: %s", err))
		return
	}

	data.Id = types.StringValue(device.ID)
	data.UserId = types.StringValue(device.UserId)
	data.Name = types.StringValue(device.Name)
	data.PublicKey = types.StringValue(device.PublicKey)
	data.Description = types.StringValue(device.Description)
	data.IPv4 = types.StringValue(device.IPv4)
	data.IPv6 = types.StringValue(device.IPv6)
	// data.AllowedIPs, = types.ListValue(device.AllowedIPs)
	data.Endpoint = types.StringValue(device.Endpoint)
	data.PresharedKey = types.StringValue(device.PresharedKey)
	data.MTU = types.Int64Value(int64(device.MTU))
	// data.DNS = types.StringValue(device.DNS)
	data.PersistentKeepalive = types.Int64Value(int64(device.PersistentKeepalive))
	data.UseDefaultAllowedIPs = types.BoolValue(device.UseDefaultAllowedIPs)
	data.UseDefaultDNS = types.BoolValue(device.UseDefaultDNS)
	data.UseDefaultEndpoint = types.BoolValue(device.UseDefaultEndpoint)
	data.UseDefaultMTU = types.BoolValue(device.UseDefaultMTU)
	data.UseDefaultPersistentKeepalive = types.BoolValue(device.UseDefaultPersistentKeepalive)

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

	device, err := r.client.UpdateDevice(data.Id.ValueString(), fz.Device{
		UserId:      data.UserId.ValueString(),
		Name:        data.Name.ValueString(),
		PublicKey:   data.PublicKey.ValueString(),
		Description: data.Description.ValueString(),
		IPv4:        data.IPv4.ValueString(),
		IPv6:        data.IPv6.ValueString(),
		// AllowedIPs:  data.AllowedIPs.ValueList(),
		Endpoint:     data.Endpoint.ValueString(),
		PresharedKey: data.PresharedKey.ValueString(),
		MTU:          int(data.MTU.ValueInt64()),
		// DNS:                           data.DNS.ValueString(),
		PersistentKeepalive:           int(data.PersistentKeepalive.ValueInt64()),
		UseDefaultDNS:                 data.UseDefaultDNS.ValueBool(),
		UseDefaultEndpoint:            data.UseDefaultEndpoint.ValueBool(),
		UseDefaultMTU:                 data.UseDefaultMTU.ValueBool(),
		UseDefaultAllowedIPs:          data.UseDefaultAllowedIPs.ValueBool(),
		UseDefaultPersistentKeepalive: data.UseDefaultPersistentKeepalive.ValueBool(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update device, got error: %s", err))
		return
	}

	data.Id = types.StringValue(device.ID)
	data.UserId = types.StringValue(device.UserId)
	data.Name = types.StringValue(device.Name)
	data.PublicKey = types.StringValue(device.PublicKey)
	data.Description = types.StringValue(device.Description)
	data.IPv4 = types.StringValue(device.IPv4)
	data.IPv6 = types.StringValue(device.IPv6)
	// data.AllowedIPs, = types.ListValue(device.AllowedIPs)
	data.Endpoint = types.StringValue(device.Endpoint)
	data.PresharedKey = types.StringValue(device.PresharedKey)
	data.MTU = types.Int64Value(int64(device.MTU))
	// data.DNS = types.StringValue(device.DNS)
	data.PersistentKeepalive = types.Int64Value(int64(device.PersistentKeepalive))
	data.UseDefaultAllowedIPs = types.BoolValue(device.UseDefaultAllowedIPs)
	data.UseDefaultDNS = types.BoolValue(device.UseDefaultDNS)
	data.UseDefaultEndpoint = types.BoolValue(device.UseDefaultEndpoint)
	data.UseDefaultMTU = types.BoolValue(device.UseDefaultMTU)
	data.UseDefaultPersistentKeepalive = types.BoolValue(device.UseDefaultPersistentKeepalive)

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

	err := r.client.DeleteDevice(data.Id.ValueString())

	if err != nil {
		// resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete device, got error: %s", err))
		return
	}
}

func (r *DeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
