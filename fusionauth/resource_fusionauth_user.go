package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func newUser() *schema.Resource {
	return &schema.Resource{
		Create: createUser,
		Read:   readUser,
		Update: updateUser,
		Delete: deleteUser,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The unique Id of the tenant used to scope this API request.",
				ValidateFunc: validation.IsUUID,
			},
			"send_set_password_email": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates to FusionAuth to send the User an email asking them to set their password. The Email Template that is used is configured in the System Configuration setting for Set Password Email Template.",
			},
			"skip_verification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates to FusionAuth that it should skip email verification even if it is enabled. This is useful for creating admin or internal User accounts.",
			},
			"birth_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An ISO-8601 formatted date of the User’s birthdate such as YYYY-MM-DD.",
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold any information about a User that should be persisted.",
			},
			"email": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"username", "email"},
				Description:  "The User’s email address. An email address is a unique in FusionAuth and stored in lower case.",
			},
			"encryption_scheme": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"salted-md5",
					"salted-sha256",
					"salted-hmac-sha256",
					"salted-pbkdf2-hmac-sha256",
					"bcrypt",
				}, false),
				Description: "The method for encrypting the User’s password.",
			},
			"expiry": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The expiration instant of the User’s account. An expired user is not permitted to login.",
			},
			// "factor": {
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// },
			"first_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The first name of the User.",
			},
			"full_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The User’s full name as a separate field that is not calculated from firstName and lastName.",
			},
			"image_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL that points to an image file that is the User’s profile image.",
			},
			"last_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The User’s last name.",
			},
			"middle_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The User’s middle name.",
			},
			"mobile_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The User’s mobile phone number. This is useful is you will be sending push notifications or SMS messages to the User.",
			},
			"parent_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email address of the user’s parent or guardian. This field is used to allow a child user to identify their parent so FusionAuth can make a request to the parent to confirm the parent relationship.",
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(8, 256),
				Description:  "The User’s plain texts password. This password will be hashed and the provided value will never be stored and cannot be retrieved.",
			},
			"password_change_required": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Indicates that the User’s password needs to be changed during their next login attempt.",
			},
			"preferred_languages": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An array of locale strings that give, in order, the User’s preferred languages. These are important for email templates and other localizable text.",
			},
			"timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The User’s preferred timezone. The string must be in an IANA time zone format. For example:",
			},
			"two_factor_delivery": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"None",
					"TextMessage",
				}, false),
				Default:     "None",
				Description: "The User’s preferred delivery for verification codes during a two factor login request.",
			},
			"two_factor_enabled": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Determines if the User has two factor authentication enabled for their account or not.",
			},
			"two_factor_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The Base64 encoded secret used to generate Two Factor verification codes.",
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"username", "email"},
				Description:  "The username of the User. The username is stored and returned as a case sensitive value, however a username is considered unique regardless of the case. bob is considered equal to BoB so either version of this username can be used whenever providing it as input to an API.",
			},
			"username_status": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ACTIVE",
					"PENDING",
					"REJECTED",
				}, false),
				Default:     "ACTIVE",
				Description: "The current status of the username. This is used if you are moderating usernames via CleanSpeak.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func buildUser(data *schema.ResourceData) fusionauth.UserRequest {
	u := fusionauth.UserRequest{
		User: fusionauth.User{
			TenantId:           data.Get("tenant_id").(string),
			BirthDate:          data.Get("birth_date").(string),
			Data:               data.Get("data").(map[string]interface{}),
			Email:              data.Get("email").(string),
			Expiry:             int64(data.Get("expiry").(int)),
			FirstName:          data.Get("first_name").(string),
			FullName:           data.Get("full_name").(string),
			ImageUrl:           data.Get("image_url").(string),
			LastName:           data.Get("last_name").(string),
			MiddleName:         data.Get("middle_name").(string),
			MobilePhone:        data.Get("mobile_phone").(string),
			ParentEmail:        data.Get("parent_email").(string),
			PreferredLanguages: handleStringSlice("preferred_languages", data),
			Timezone:           data.Get("timezone").(string),
			SecureIdentity: fusionauth.SecureIdentity{
				EncryptionScheme:       data.Get("encryption_scheme").(string),
				Password:               data.Get("password").(string),
				PasswordChangeRequired: data.Get("password_change_required").(bool),
				TwoFactorDelivery:      fusionauth.TwoFactorDelivery(data.Get("two_factor_delivery").(string)),
				TwoFactorEnabled:       data.Get("two_factor_enabled").(bool),
				TwoFactorSecret:        data.Get("two_factor_secret").(string),
				Username:               data.Get("username").(string),
				UsernameStatus:         fusionauth.ContentStatus(data.Get("username_status").(string)),
			},
		},
		SendSetPasswordEmail: data.Get("send_set_password_email").(bool),
		SkipVerification:     data.Get("skip_verification").(bool),
	}
	return u
}

func createUser(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	u := buildUser(data)

	oldTenantID := client.FAClient.TenantId
	client.FAClient.TenantId = u.User.TenantId
	defer func() {
		client.FAClient.TenantId = oldTenantID
	}()

	resp, faErrs, err := client.FAClient.CreateUser("", u)
	if err != nil {
		return fmt.Errorf("CreateUser err: %v", err)
	}

	if faErrs != nil {
		return fmt.Errorf("CreateUser errors: %v", faErrs)
	}
	data.SetId(resp.User.Id)
	_ = data.Set("send_set_password_email", nil)
	_ = data.Set("skip_verification", nil)
	if u.User.TenantId == "" {
		_ = data.Set("tenant_id", resp.User.TenantId)
	}
	return nil
}

func readUser(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveUser(id)
	if err != nil {
		return err
	}

	if faErrs != nil {
		return fmt.Errorf("RetrieveUser errors: %v", faErrs)
	}

	_ = data.Set("tenant_id", resp.User.TenantId)
	_ = data.Set("birth_date", resp.User.BirthDate)
	_ = data.Set("data", resp.User.Data)
	_ = data.Set("email", resp.User.Email)
	_ = data.Set("expiry", resp.User.Expiry)
	_ = data.Set("first_name", resp.User.FirstName)
	_ = data.Set("full_name", resp.User.FullName)
	_ = data.Set("image_url", resp.User.ImageUrl)
	_ = data.Set("last_name", resp.User.LastName)
	_ = data.Set("middle_name", resp.User.MiddleName)
	_ = data.Set("mobile_phone", resp.User.MobilePhone)
	_ = data.Set("parent_email", resp.User.ParentEmail)
	_ = data.Set("preferred_languages", resp.User.PreferredLanguages)
	_ = data.Set("timezone", resp.User.Timezone)
	_ = data.Set("two_factor_delivery", resp.User.TwoFactorDelivery)
	_ = data.Set("two_factor_enabled", resp.User.TwoFactorEnabled)
	_ = data.Set("two_factor_secret", resp.User.TwoFactorSecret)
	_ = data.Set("username", resp.User.Username)
	_ = data.Set("username_status", resp.User.UsernameStatus)
	_ = data.Set("encryption_scheme", resp.User.EncryptionScheme)
	_ = data.Set("password_change_required", resp.User.PasswordChangeRequired)

	return nil
}

func updateUser(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	u := buildUser(data)

	_, faErrs, err := client.FAClient.UpdateUser(data.Id(), u)
	if err != nil {
		return err
	}

	if faErrs != nil {
		return fmt.Errorf("UpdateUser errors: %v", faErrs)
	}

	return nil
}

func deleteUser(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	_, faErrs, err := client.FAClient.DeleteUser(id)
	if err != nil {
		return err
	}

	if faErrs != nil {
		return fmt.Errorf("DeleteUser errors: %v", faErrs)
	}

	return nil
}
