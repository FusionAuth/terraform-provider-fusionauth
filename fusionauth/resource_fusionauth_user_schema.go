package fusionauth

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func userSchemaV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The Id to use for the new User. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "The unique Id of the tenant used to scope this API request.",
				ValidateFunc: validation.IsUUID,
			},
			"application_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "An optional Application Id. When this value is provided, it will be used to resolve an application specific email template if you have configured transactional emails such as setup password, email verification and others.",
				ValidateFunc: validation.IsUUID,
			},
			"disable_domain_block": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "A tenant has the option to configure one or more email domains to be blocked in order to restrict email domains during user create or update.",
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
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "An object that can hold any information about a User that should be persisted. Must be a JSON encoded string.",
				DiffSuppressFunc: diffSuppressJSON,
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
				Default:  "salted-pbkdf2-hmac-sha256",
				ValidateDiagFunc: WarnStringInSlice([]string{
					"salted-md5",
					"salted-sha256",
					"salted-hmac-sha256",
					"salted-pbkdf2-hmac-sha256",
					"salted-pbkdf2-hmac-sha256-512",
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
			// 	Computed:    true,
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"two_factor_method_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique Id of the method.",
						},
						"authenticator_algorithm": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HmacSHA1",
							}, false),
							Description: "The algorithm used by the TOTP authenticator. With the current implementation, this will always be HmacSHA1.",
						},
						"authenticator_code_length": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The length of code generated by the TOTP. With the current implementation, this will always be 6.",
						},
						"authenticator_time_step": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
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
	}
}

func upgradeUserSchemaV0ToV1(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if len(rawState) == 0 {
		// if raw state is nil or empty, return.
		return map[string]interface{}{}, nil
	}

	// Remove deprecated fields from state.
	delete(rawState, "two_factor_delivery")
	delete(rawState, "two_factor_enabled")
	delete(rawState, "two_factor_secret")

	// Migrate data types.
	if currentValue, ok := rawState["data"]; ok && currentValue != nil {
		if v, ok := currentValue.(map[string]interface{}); ok && len(v) > 0 {
			bytes, err := json.Marshal(currentValue)
			if err != nil {
				return rawState, err
			}

			rawState["data"] = string(bytes)
		} else {
			rawState["data"] = ""
		}
	}

	return rawState, nil
}

func userSchemaV0() *schema.Resource {
	v0Diff := &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Datatype changes
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold any information about a User that should be persisted.",
			},

			// Property Deprecations
			"two_factor_delivery": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"None",
					"TextMessage",
				}, false),
				Default:     "None",
				Deprecated:  "Removed in Fusionauth version 1.26.0",
				Description: "The User’s preferred delivery for verification codes during a two factor login request.",
			},
			"two_factor_enabled": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Deprecated:  "Removed in Fusionauth version 1.26.0",
				Description: "Determines if the User has two factor authentication enabled for their account or not.",
			},
			"two_factor_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Deprecated:  "Removed in Fusionauth version 1.26.0",
				Description: "The Base64 encoded secret used to generate Two Factor verification codes.",
			},
		},
	}

	return injectSchemaChanges(userSchemaV1(), v0Diff)
}
