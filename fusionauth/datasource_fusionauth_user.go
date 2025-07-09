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
				ExactlyOneOf: []string{"email", "user_id", "username"},
				Description:  "The Id to use for the new User. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"email", "phone_number", "user_id", "username"},
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"email", "user_id", "username"},
				Description:  "The User’s email address.",
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
			"identities": {
				Type:        schema.TypeList,
				Computed:    true,
				Default:     nil,
				Description: "The list of identities that exist for a User.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display value for the identity. Only used for username type identities. If the unique username feature is not enabled, this value will be the same as user.identities[x].value. Otherwise, it will be the username the User has chosen. For primary username identities, this will be the same value as user.username .",
						},
						"insert_instant": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The instant when the identity was created.",
						},
						"last_login_instant": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The instant when the identity was last used to log in. If a User has multiple identity types (username, email, and phoneNumber), then this value will represent the specific identity they last used to log in. This contrasts with user.lastLoginInstant, which represents the last time any of the User’s identities was used to log in.",
						},
						"last_update_instant": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The instant when the identity was last updated.",
						},
						"moderation_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The current status of the username. This is used if you are moderating usernames via CleanSpeak.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The identity type.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value represented by the identity.",
						},
						"verified": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether verification was actually performed on the identity by FusionAuth.",
						},
						"verified_instant": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The instant when verification was performed on the identity.",
						},
						"verified_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reason the User’s identity was verified or not verified.",
						},
					},
				},
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
			"phone_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The User’s phone number.",
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
			"verification_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Default:     nil,
				Description: "The list of all verifications that exist for a user. This includes the email and phone identities that a user may have. The values from emailVerificationId and emailVerificationOneTimeCode are legacy fields and will also be present in this list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"verification_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A verification Id.",
						},
						"one_time_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A one time code that will be paired with the verificationIds[x].id .",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The identity type that the verification Id is for. This identity type, along with verificationIds[x].value , matches exactly one identity via user.identities[x].type .",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The identity value that the verification Id is for. This identity value, along with verificationIds[x].type , matches exactly one identity via user.identities[x].value .",
						},
					},
				},
			},
		},
	}
}

//nolint:gocyclo,gocognit
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
	} else if username, ok := data.GetOk("username"); ok {
		searchID = username.(string)
		resp, faErrs, err = client.FAClient.RetrieveUserByUsername(searchID)
	} else if email, ok := data.GetOk("email"); ok {
		searchID = email.(string)
		resp, faErrs, err = client.FAClient.RetrieveUserByEmail(searchID)
	} else if phoneNumber, ok := data.GetOk("phone_number"); ok {
		searchID = phoneNumber.(string)
		resp, faErrs, err = client.FAClient.RetrieveUserByLoginId(searchID)
	} else {
		return diag.Errorf("user_id, username, email or phone_number must be set")
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
	if resp.User.Identities != nil {
		identities := make([]map[string]interface{}, len(resp.User.Identities))
		for i, identity := range resp.User.Identities {
			identities[i] = map[string]interface{}{
				"display_value":       identity.DisplayValue,
				"insert_instant":      identity.InsertInstant,
				"last_login_instant":  identity.LastLoginInstant,
				"last_update_instant": identity.LastUpdateInstant,
				"moderation_status":   identity.ModerationStatus,
				"type":                identity.Type,
				"value":               identity.Value,
				"verified":            identity.Verified,
				"verified_instant":    identity.VerifiedInstant,
				"verified_reason":     identity.VerifiedReason,
			}
		}
		if err := data.Set("identities", identities); err != nil {
			return diag.Errorf("user.identities: %s", err.Error())
		}
	} else {
		// If Identities is nil in the response, set it as an empty list in the state.
		if err := data.Set("identities", []map[string]interface{}{}); err != nil {
			return diag.Errorf("user.identities: %s", err.Error())
		}
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
	if err := data.Set("phone_number", resp.User.PhoneNumber); err != nil {
		return diag.Errorf("user.phone_number: %s", err.Error())
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
	if resp.VerificationIds != nil {
		verificationIDs := make([]map[string]interface{}, len(resp.VerificationIds))
		for i, verificationID := range resp.VerificationIds {
			verificationIDs[i] = map[string]interface{}{
				"verification_id": verificationID.Id,
				"one_time_code":   verificationID.OneTimeCode,
				"type":            verificationID.Type,
				"value":           verificationID.Value,
			}
		}
		if err := data.Set("verification_ids", verificationIDs); err != nil {
			return diag.Errorf("user.verification_ids: %s", err.Error())
		}
	} else {
		// If VerificationIds is nil in the response, set it as an empty list in the state.
		if err := data.Set("verification_ids", []map[string]interface{}{}); err != nil {
			return diag.Errorf("user.verification_ids: %s", err.Error())
		}
	}

	return nil
}
