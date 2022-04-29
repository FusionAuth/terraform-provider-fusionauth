package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The unique Id of the tenant used to scope this API request.",
				ValidateFunc: validation.IsUUID,
			},
			"user_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"user_id", "username"},
				Description:  "The Id to use for the new User. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"user_id", "username"},
				Description:  "The username of the User.",
			},
			"active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if the User is active. False if the User has been deactivated. Deactivated Users will not be able to login.",
			},
			"birth_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An ISO-8601 formatted date of the User’s birthdate such as YYYY-MM-DD.",
			},
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A JSON serialised string that can hold any information about the User.",
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The User’s email address.",
			},
			"expiry": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The expiration instant of the User’s account. An expired user is not permitted to login.",
			},
			"first_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The first name of the User.",
			},
			"full_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The User’s full name.",
			},
			"image_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL that points to an image file that is the User’s profile image.",
			},
			"last_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The User’s last name.",
			},
			"middle_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The User’s middle name.",
			},
			"mobile_phone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The User’s mobile phone number.",
			},
			"parent_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address of the user’s parent or guardian.",
			},
			"password_change_required": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates that the User’s password needs to be changed during their next login attempt.",
			},
			"preferred_languages": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "An array of locale strings that give, in order, the User’s preferred languages.",
			},
			"timezone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The User’s preferred timezone.",
			},
			"username_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the username. This is used if you are moderating usernames via CleanSpeak.",
			},
		},
	}
}

func dataSourceUserRead(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	oldTenantID := client.FAClient.TenantId
	client.FAClient.TenantId = data.Get("tenant_id").(string)
	defer func() {
		client.FAClient.TenantId = oldTenantID
	}()

	var searchID string
	var resp *fusionauth.UserResponse
	var faErrs *fusionauth.Errors
	var err error

	// Either `user_id` or `username` are guaranteed to be set
	if userID, ok := data.GetOk("user_id"); ok {
		searchID = userID.(string)
		resp, faErrs, err = client.FAClient.RetrieveUser(searchID)
	} else {
		searchID = data.Get("username").(string)
		resp, faErrs, err = client.FAClient.RetrieveUserByUsername(searchID)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return diag.Errorf("couldn't find %s", searchID)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(resp.User.Id)

	if err := data.Set("tenant_id", resp.User.TenantId); err != nil {
		return diag.Errorf("user.tenant_id: %s", err.Error())
	}
	if err := data.Set("user_id", resp.User.Id); err != nil {
		return diag.Errorf("user.user_id: %s", err.Error())
	}
	if err := data.Set("username", resp.User.Username); err != nil {
		return diag.Errorf("user.username: %s", err.Error())
	}

	if err := data.Set("active", resp.User.Active); err != nil {
		return diag.Errorf("user.active: %s", err.Error())
	}
	if err := data.Set("birth_date", resp.User.BirthDate); err != nil {
		return diag.Errorf("user.birth_date: %s", err.Error())
	}
	if userData, diags := mapStringInterfaceToJSONString(resp.User.Data); diags != nil {
		return diags
	} else if err := data.Set("data", userData); err != nil {
		return diag.Errorf("user.data: %s", err.Error())
	}
	if err := data.Set("email", resp.User.Email); err != nil {
		return diag.Errorf("user.email: %s", err.Error())
	}
	if err := data.Set("expiry", resp.User.Expiry); err != nil {
		return diag.Errorf("user.expiry: %s", err.Error())
	}
	if err := data.Set("first_name", resp.User.FirstName); err != nil {
		return diag.Errorf("user.first_name: %s", err.Error())
	}
	if err := data.Set("full_name", resp.User.FullName); err != nil {
		return diag.Errorf("user.full_name: %s", err.Error())
	}
	if err := data.Set("image_url", resp.User.ImageUrl); err != nil {
		return diag.Errorf("user.image_url: %s", err.Error())
	}
	if err := data.Set("last_name", resp.User.LastName); err != nil {
		return diag.Errorf("user.last_name: %s", err.Error())
	}
	if err := data.Set("middle_name", resp.User.MiddleName); err != nil {
		return diag.Errorf("user.middle_name: %s", err.Error())
	}
	if err := data.Set("mobile_phone", resp.User.MobilePhone); err != nil {
		return diag.Errorf("user.mobile_phone: %s", err.Error())
	}
	if err := data.Set("parent_email", resp.User.ParentEmail); err != nil {
		return diag.Errorf("user.parent_email: %s", err.Error())
	}
	if err := data.Set("password_change_required", resp.User.PasswordChangeRequired); err != nil {
		return diag.Errorf("user.password_change_required: %s", err.Error())
	}
	if err := data.Set("preferred_languages", resp.User.PreferredLanguages); err != nil {
		return diag.Errorf("user.preferred_languages: %s", err.Error())
	}
	if err := data.Set("timezone", resp.User.Timezone); err != nil {
		return diag.Errorf("user.timezone: %s", err.Error())
	}
	if err := data.Set("username_status", resp.User.UsernameStatus); err != nil {
		return diag.Errorf("user.username_status: %s", err.Error())
	}

	return nil
}
