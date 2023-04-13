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
var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

// UserResource defines the resource implementation.
type UserResource struct {
	client *fz.Client
}

// UserResourceModel describes the resource data model.
type UserResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Email              types.String `tfsdk:"email"`
	Role               types.String `tfsdk:"role"`
	LastSignedInAt     types.String `tfsdk:"last_signed_in_at"`
	LastSignedInMethod types.String `tfsdk:"last_signed_in_method"`
	UpdatedAt          types.String `tfsdk:"updated_at"`
	InsertedAt         types.String `tfsdk:"inserted_at"`
	DisabledAt         types.String `tfsdk:"disabled_at"`
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "User resource",

		Attributes: map[string]schema.Attribute{
			"disabled_at": schema.StringAttribute{
				MarkdownDescription: "User disabled at",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "User updated at",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"inserted_at": schema.StringAttribute{
				MarkdownDescription: "User inserted at",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"last_signed_in_method": schema.StringAttribute{
				MarkdownDescription: "User last signed in method",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"last_signed_in_at": schema.StringAttribute{
				MarkdownDescription: "User last signed in at",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "User role",
				Required:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "User email",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "User identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.CreateUser(fz.User{
		Email: data.Email.String(),
		Role:  data.Role.String(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create user, got error: %s", err))
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.

	data.Id = types.StringValue(user.ID)
	data.Email = types.StringValue(user.Email)
	data.Role = types.StringValue(user.Role)
	data.LastSignedInAt = types.StringValue(user.LastSignedInAt)
	data.LastSignedInMethod = types.StringValue(user.LastSignedInMethod)
	data.UpdatedAt = types.StringValue(user.UpdatedAt)
	data.InsertedAt = types.StringValue(user.InsertedAt)
	data.DisabledAt = types.StringValue(user.DisabledAt)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetUser(data.Id.String())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user, got error: %s", err))
		return
	}

	data.Id = types.StringValue(user.ID)
	data.Email = types.StringValue(user.Email)
	data.Role = types.StringValue(user.Role)
	data.LastSignedInAt = types.StringValue(user.LastSignedInAt)
	data.LastSignedInMethod = types.StringValue(user.LastSignedInMethod)
	data.UpdatedAt = types.StringValue(user.UpdatedAt)
	data.InsertedAt = types.StringValue(user.InsertedAt)
	data.DisabledAt = types.StringValue(user.DisabledAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.UpdateUser(data.Id.String(), fz.User{
		Email: data.Email.String(),
		Role:  data.Role.String(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update user, got error: %s", err))
		return
	}

	data.Id = types.StringValue(user.ID)
	data.Email = types.StringValue(user.Email)
	data.Role = types.StringValue(user.Role)
	data.LastSignedInAt = types.StringValue(user.LastSignedInAt)
	data.LastSignedInMethod = types.StringValue(user.LastSignedInMethod)
	data.UpdatedAt = types.StringValue(user.UpdatedAt)
	data.InsertedAt = types.StringValue(user.InsertedAt)
	data.DisabledAt = types.StringValue(user.DisabledAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(data.Id.String())

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete user, got error: %s", err))
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
