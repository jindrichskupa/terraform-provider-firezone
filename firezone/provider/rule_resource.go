package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	fz "github.com/jindrichskupa/firezone-client-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RuleResource{}
var _ resource.ResourceWithImportState = &RuleResource{}

func NewRuleResource() resource.Resource {
	return &RuleResource{}
}

// RuleResource defines the resource implementation.
type RuleResource struct {
	client *fz.Client
}

// RuleResourceModel describes the resource data model.
type RuleResourceModel struct {
	Id          types.String `tfsdk:"id"`
	UserId      types.String `tfsdk:"user_id"`
	Action      types.String `tfsdk:"action"`
	Destination types.String `tfsdk:"destination"`
	PortRange   types.String `tfsdk:"port_range"`
	PortType    types.String `tfsdk:"port_type"`
}

func (r *RuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rule"
}

func (r *RuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Rule resource",

		Attributes: map[string]schema.Attribute{
			"port_type": schema.StringAttribute{
				MarkdownDescription: "Rule port type",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					// These are example validators from terraform-plugin-framework-validators
					stringvalidator.LengthBetween(1, 256),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(tcp|udp)$`),
						"must be either 'tcp' or 'udp'",
					),
				},
			},
			"port_range": schema.StringAttribute{
				MarkdownDescription: "Rule port range",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					// These are example validators from terraform-plugin-framework-validators
					stringvalidator.LengthBetween(1, 256),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^([1-9]?[0-9]*|[1-9]?[0-9]* - [1-9]?[0-9]*)$`),
						"must contain single port or port range in format 'port' or 'port - port' (with spaces)",
					),
				},
			},
			"destination": schema.StringAttribute{
				MarkdownDescription: "Rule destination",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Rule action",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					// These are example validators from terraform-plugin-framework-validators
					stringvalidator.LengthBetween(1, 256),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(accept|drop)$`),
						"must be either 'accept' or 'drop'",
					),
				},
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "Rule user id",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Rule identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
		},
	}
}

func (r *RuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *RuleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.CreateRule(fz.Rule{
		UserId:      data.UserId.ValueString(),
		Action:      data.Action.ValueString(),
		Destination: data.Destination.ValueString(),
		PortRange:   data.PortRange.ValueString(),
		PortType:    data.PortType.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create rule, got error: %s", err))
		return
	}

	data.Id = types.StringValue(rule.ID)
	data.UserId = types.StringValue(rule.UserId)
	data.Action = types.StringValue(rule.Action)
	data.Destination = types.StringValue(rule.Destination)
	data.PortRange = types.StringValue(rule.PortRange)
	data.PortType = types.StringValue(rule.PortType)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *RuleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetRule(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read rule, got error: %s", err))
		return
	}

	data.Id = types.StringValue(rule.ID)
	data.UserId = types.StringValue(rule.UserId)
	data.Action = types.StringValue(rule.Action)
	data.Destination = types.StringValue(rule.Destination)
	data.PortRange = types.StringValue(rule.PortRange)
	data.PortType = types.StringValue(rule.PortType)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *RuleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.UpdateRule(data.Id.ValueString(), fz.Rule{
		UserId:      data.UserId.ValueString(),
		Action:      data.Action.ValueString(),
		Destination: data.Destination.ValueString(),
		PortRange:   data.PortRange.ValueString(),
		PortType:    data.PortType.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update rule, got error: %s", err))
		return
	}

	data.Id = types.StringValue(rule.ID)
	data.UserId = types.StringValue(rule.UserId)
	data.Action = types.StringValue(rule.Action)
	data.Destination = types.StringValue(rule.Destination)
	data.PortRange = types.StringValue(rule.PortRange)
	data.PortType = types.StringValue(rule.PortType)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *RuleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRule(data.Id.ValueString())

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete rule, got error: %s", err))
		return
	}
}

func (r *RuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
