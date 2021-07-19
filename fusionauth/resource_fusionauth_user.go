package fusionauth

import (
	"fmt"
	"net/http"

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
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The factor used by the password encryption scheme. If not provided, the PasswordEncryptor provides a default value. Generally this will be used as an iteration count to generate the hash. The actual use of this value is up to the PasswordEncryptor implementation.",
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
			"two_factor_methods": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authenticator_algorithm": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HmacSHA1",
							}, false),
							Default:     "HmacSHA1",
							Description: "The algorithm used by the TOTP authenticator. With the current implementation, this will always be HmacSHA1.",
						},
						"authenticator_code_length": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     6,
							Description: "The length of code generated by the TOTP. With the current implementation, this will always be 6.",
						},
						"authenticator_time_step": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     30,
							Description: "The time-step size in seconds. With the current implementation, this will always be 30.",
						},
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value of the email address for this method. Only present if user.twoFactor.methods[x].method is email.",
						},
						"method": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"authenticator",
								"email",
								"sms",
							}, false),
							Description: "The type of this method. There will also be an object with the same value containing additional information about this method.",
						},
						"mobile_phone": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value of the mobile phone for this method. Only present if user.twoFactor.methods[x].method is sms.",
						},
						"secret": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "A base64 encoded secret",
						},
					},
				},
			},
			"two_factor_recovery_codes": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Sensitive:   true,
				Computed:    true,
				Description: "A list of recovery codes. These may be used in place of a code provided by an MFA factor. They are single use. If a recovery code is used in a disable request, all MFA methods are removed. If, after that, a Multi-Factor method is added, a new set of recovery codes will be generated.",
			},
			"two_factor_delivery": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"None",
					"TextMessage",
				}, false),
				Default:     "None",
				Removed:     "Removed in version 1.26.0",
				Description: "The User’s preferred delivery for verification codes during a two factor login request.",
			},
			"two_factor_enabled": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Removed:     "Removed in version 1.26.0",
				Description: "Determines if the User has two factor authentication enabled for their account or not.",
			},
			"two_factor_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Removed:     "Removed in version 1.26.0",
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
				Username:               data.Get("username").(string),
				UsernameStatus:         fusionauth.ContentStatus(data.Get("username_status").(string)),
			},
			TwoFactor: fusionauth.UserTwoFactorConfiguration{
				Methods:       buildTwoFactorMethod("two_factor_methods", data),
				RecoveryCodes: handleStringSlice("two_factor_recovery_codes", data),
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

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}
	data.SetId(resp.User.Id)
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

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	if err := data.Set("tenant_id", resp.User.TenantId); err != nil {
		return fmt.Errorf("user.tenant_id: %s", err.Error())
	}
	if err := data.Set("birth_date", resp.User.BirthDate); err != nil {
		return fmt.Errorf("user.birth_date: %s", err.Error())
	}
	if err := data.Set("data", resp.User.Data); err != nil {
		return fmt.Errorf("user.data: %s", err.Error())
	}
	if err := data.Set("email", resp.User.Email); err != nil {
		return fmt.Errorf("user.email: %s", err.Error())
	}
	if err := data.Set("expiry", resp.User.Expiry); err != nil {
		return fmt.Errorf("user.expiry: %s", err.Error())
	}
	if err := data.Set("first_name", resp.User.FirstName); err != nil {
		return fmt.Errorf("user.first_name: %s", err.Error())
	}
	if err := data.Set("full_name", resp.User.FullName); err != nil {
		return fmt.Errorf("user.full_name: %s", err.Error())
	}
	if err := data.Set("image_url", resp.User.ImageUrl); err != nil {
		return fmt.Errorf("user.image_url: %s", err.Error())
	}
	if err := data.Set("last_name", resp.User.LastName); err != nil {
		return fmt.Errorf("user.last_name: %s", err.Error())
	}
	if err := data.Set("middle_name", resp.User.MiddleName); err != nil {
		return fmt.Errorf("user.middle_name: %s", err.Error())
	}
	if err := data.Set("mobile_phone", resp.User.MobilePhone); err != nil {
		return fmt.Errorf("user.mobile_phone: %s", err.Error())
	}
	if err := data.Set("parent_email", resp.User.ParentEmail); err != nil {
		return fmt.Errorf("user.parent_email: %s", err.Error())
	}
	if err := data.Set("preferred_languages", resp.User.PreferredLanguages); err != nil {
		return fmt.Errorf("user.preferred_languages: %s", err.Error())
	}
	if err := data.Set("timezone", resp.User.Timezone); err != nil {
		return fmt.Errorf("user.timezone: %s", err.Error())
	}
	if err := data.Set("username", resp.User.Username); err != nil {
		return fmt.Errorf("user.username: %s", err.Error())
	}
	if err := data.Set("username_status", resp.User.UsernameStatus); err != nil {
		return fmt.Errorf("user.username_status: %s", err.Error())
	}
	if err := data.Set("password_change_required", resp.User.PasswordChangeRequired); err != nil {
		return fmt.Errorf("user.password_change_required: %s", err.Error())
	}

	if err := data.Set("two_factor_recovery_codes", resp.User.TwoFactor.RecoveryCodes); err != nil {
		return fmt.Errorf("user.two_factor_recovery_codes: %s", err.Error())
	}

	tfms := make([]map[string]interface{}, 0, len(resp.User.TwoFactor.Methods))
	for idx := range resp.User.TwoFactor.Methods {
		m := resp.User.TwoFactor.Methods[idx]
		tfms = append(tfms, map[string]interface{}{
			"authenticator_algorithm":   m.Authenticator.Algorithm,
			"authenticator_code_length": m.Authenticator.CodeLength,
			"authenticator_time_step":   m.Authenticator.TimeStep,
			"email":                     m.Email,
			"method":                    m.Method,
			"mobile_phone":              m.MobilePhone,
			"secret":                    m.Secret,
		})
	}
	if err = data.Set("two_factor_methods", tfms); err != nil {
		return fmt.Errorf("user.two_factor_methods: %s", err.Error())
	}
	return nil
}

func updateUser(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	u := buildUser(data)

	resp, faErrs, err := client.FAClient.UpdateUser(data.Id(), u)
	if err != nil {
		return err
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return nil
}

func deleteUser(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteUser(id)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return err
	}

	return nil
}

func buildTwoFactorMethod(key string, data *schema.ResourceData) []fusionauth.TwoFactorMethod {
	s := data.Get(key)
	set, ok := s.(*schema.Set)
	if !ok {
		return []fusionauth.TwoFactorMethod{}
	}
	l := set.List()

	tfms := make([]fusionauth.TwoFactorMethod, 0, len(l))
	for _, x := range l {
		r := x.(map[string]interface{})
		tfms = append(tfms, fusionauth.TwoFactorMethod{
			Authenticator: fusionauth.AuthenticatorConfiguration{
				Algorithm:  fusionauth.TOTPAlgorithm(r["authenticator_algorithm"].(string)),
				CodeLength: r["authenticator_code_length"].(int),
				TimeStep:   r["authenticator_time_step"].(int),
			},
			Email:       r["email"].(string),
			Method:      r["method"].(string),
			MobilePhone: r["mobile_phone"].(string),
			Secret:      r["secret"].(string),
		})
	}

	return tfms
}
