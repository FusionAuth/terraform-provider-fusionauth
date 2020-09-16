package fusionauth

import (
	"fmt"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func newTenant() *schema.Resource {
	return &schema.Resource{
		Create: createTenant,
		Read:   readTenant,
		Update: updateTenant,
		Delete: deleteTenant,
		Schema: map[string]*schema.Schema{
			"source_tentant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The optional Id of an existing Tenant to make a copy of. If present, the tenant.id and tenant.name values of the request body will be applied to the new Tenant, all other values will be copied from the source Tenant to the new Tenant.",
				ValidateFunc: validation.IsUUID,
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold any information about the Tenant that should be persisted.",
			},
			"email_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem:     newEmailConfiguration(),
			},
			"event_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"user.action",
								"user.bulk.create",
								"user.create",
								"user.email.verified",
								"user.update",
								"user.deactivate",
								"user.reactivate",
								"user.login.success",
								"user.login.failed",
								"user.password.breach",
								"user.registration.create",
								"user.registration.update",
								"user.registration.delete",
								"user.registration.verified",
								"user.delete",
								"jwt.public-key.update",
								"jwt.refresh-token.revoke",
								"jwt.refresh",
							}, false),
							Description: "the event type",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether or not FusionAuth should send these types of events to any configured Webhooks.",
						},
						"transaction_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"None",
								"Any",
								"SimpleMajority",
								"SuperMajority",
								"AbsoluteMajority",
							}, false),
							Description: "The transaction type that FusionAuth uses when sending these types of events to any configured Webhooks.",
						},
					},
				},
			},
			"external_identifier_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem:     newExternalIdentifierConfiguration(),
			},
			"failed_authentication_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     newFailedAuthenticationConfiguration(),
			},
			"family_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     newFamilyConfiguration(),
			},
			"http_session_max_inactive_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3600,
			},
			"issuer": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The named issuer used to sign tokens, this is generally your public fully qualified domain.",
			},
			"jwt_configuration": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_token_key_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
							Description:  "The unique id of the signing key used to sign the access token.",
						},
						"id_token_key_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
							Description:  "The unique id of the signing key used to sign the Id token.",
						},
						"refresh_token_expiration_policy": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Fixed",
								"SlidingWindow",
							}, false),
							Description: "The refresh token expiration policy.",
						},
						"refresh_token_revocation_policy_on_login_prevented": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When enabled, the refresh token will be revoked when a user action, such as locking an account based on a number of failed login attempts, prevents user login.",
						},
						"refresh_token_revocation_policy_on_password_change": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "When enabled, the refresh token will be revoked when a user changes their password.",
						},
						"refresh_token_time_to_live_in_minutes": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "The length of time in minutes a Refresh Token is valid from the time it was issued. Value must be greater than 0.",
							ValidateFunc: validation.IntAtLeast(1),
						},
						"refresh_token_usage_policy": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Reusable",
								"OneTimeUse",
							}, false),
							Description: "The refresh token usage policy.",
						},
						"time_to_live_in_seconds": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "The length of time in seconds this JWT is valid from the time it was issued. Value must be greater than 0.",
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},
			"logout_url": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The logout redirect URL when sending the user’s browser to the /oauth2/logout URI of the FusionAuth Front End. This value is only used when a logout URL is not defined in your Application.",
			},
			"maximum_password_age": {
				Optional: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"days": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     180,
							Description: "The password maximum age in days. The number of days after which FusionAuth will require a user to change their password. Required when systemConfiguration.maximumPasswordAge.enabled is set to true.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates that the maximum password age is enabled and being enforced.",
						},
					},
				},
			},
			"minimum_password_age": {
				Optional: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     30,
							Description: "e password minimum age in seconds. When enabled FusionAuth will not allow a password to be changed until it reaches this minimum age. Required when systemConfiguration.minimumPasswordAge.enabled is set to true.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates that the minimum password age is enabled and being enforced.",
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Tenant.",
			},
			"password_encryption_configuration": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Description: "The default method for encrypting the User’s password.",
						},
						"encryption_scheme_factor": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     24000,
							Description: "The factor used by the password encryption scheme. If not provided, the PasswordEncryptor provides a default value. Generally this will be used as an iteration count to generate the hash. The actual use of this value is up to the PasswordEncryptor implementation.",
						},
						"modify_encryption_scheme_on_login": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "When enabled a user’s hash configuration will be modified to match these configured settings. This can be useful to increase a password hash strength over time or upgrade imported users to a more secure encryption scheme after an initial import.",
						},
					},
				},
			},
			"password_validation_rules": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     newPasswordValidationRules(),
			},
			"theme_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The unique Id of the theme to be used to style the login page and other end user templates.",
				ValidateFunc: validation.IsUUID,
			},
			"user_delete_policy": {
				Optional: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unverified_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates that users without a verified email address will be permanently deleted after tenant.userDeletePolicy.unverified.numberOfDaysToRetain days.",
						},
						"unverified_number_of_days_to_retain": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The number of days from creation users will be retained before being deleted for not completing email verification. This field is required when tenant.userDeletePolicy.unverified.enabled is set to true. Value must be greater than 0.",
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func newFamilyConfiguration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_child_registrations": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to allow child registrations.",
			},
			"confirm_child_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The unique Id of the email template to use when confirming a child.",
			},
			"delete_orphaned_accounts": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates that child users without parental verification will be permanently deleted after tenant.familyConfiguration.deleteOrphanedAccountsDays days.",
			},
			"delete_orphaned_accounts_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				Description:  "The number of days from creation child users will be retained before being deleted for not completing parental verification. Value must be greater than 0.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether family configuration is enabled.",
			},
			"family_request_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The unique Id of the email template to use when a family request is made.",
			},
			"maximum_child_age": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      12,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The maximum age of a child. Value must be greater than 0.",
			},
			"minimum_owner_age": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Default:      21,
				Description:  "The minimum age to be an owner. Value must be greater than 0.",
			},
			"parent_email_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether a parent email is required.",
			},
			"parent_registration_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The unique Id of the email template to use for parent registration.",
			},
		},
	}
}

func newFailedAuthenticationConfiguration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action_duration": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The duration of the User Action. This value along with the actionDurationUnit will be used to set the duration of the User Action. Value must be greater than 0.",
			},
			"action_duration_unit": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MINUTES",
					"HOURS",
					"DAYS",
					"WEEKS",
					"MONTHS",
					"YEARS",
				}, false),
				Default:     "MINUTES",
				Description: "The unit of time associated with a duration.",
			},
			"reset_count_in_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				Description:  "The length of time in seconds before the failed authentication count will be reset. Value must be greater than 0.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"too_many_attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				Description:  "The number of failed attempts considered to be too many. Once this threshold is reached the specified User Action will be applied to the user for the duration specified. Value must be greater than 0.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"user_action_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the User Action that is applied when the threshold is reached for too many failed authentication attempts.",
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func newExternalIdentifierConfiguration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"authorization_grant_id_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The time in seconds until a OAuth authorization code in no longer valid to be exchanged for an access token. This is essentially the time allowed between the start of an Authorization request during the Authorization code grant and when you request an access token using this authorization code on the Token endpoint.",
				ValidateFunc: validation.IntBetween(1, 600),
			},
			"change_password_id_generator": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The length of the secure generator used for generating the change password Id.",
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"randomAlpha",
								"randomAlphaNumeric",
								"randomBytes",
								"randomDigits.",
							}, false),
							Description: "The type of the secure generator used for generating the change password Id.",
						},
					},
				},
			},
			"change_password_id_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The time in seconds until a change password Id is no longer valid and cannot be used by the Change Password API. Value must be greater than 0.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"device_code_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The time in seconds until a device code Id is no longer valid and cannot be used by the Token API. Value must be greater than 0.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"device_user_code_id_generator": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The length of the secure generator used for generating the change password Id.",
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"randomAlpha",
								"randomAlphaNumeric",
								"randomBytes",
								"randomDigits.",
							}, false),
							Description: "The type of the secure generator used for generating the change password Id.",
						},
					},
				},
			},
			"email_verification_id_generator": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The length of the secure generator used for generating the change password Id.",
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"randomAlpha",
								"randomAlphaNumeric",
								"randomBytes",
								"randomDigits.",
							}, false),
							Description: "The type of the secure generator used for generating the change password Id.",
						},
					},
				},
			},
			"email_verification_id_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The time in seconds until a email verification Id is no longer valid and cannot be used by the Verify Email API. Value must be greater than 0.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"external_authentication_id_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The time in seconds until an external authentication Id is no longer valid and cannot be used by the Token API. Value must be greater than 0.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"one_time_password_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The time in seconds until a One Time Password is no longer valid and cannot be used by the Login API. Value must be greater than 0.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"passwordless_login_generator": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The length of the secure generator used for generating the change password Id.",
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"randomAlpha",
								"randomAlphaNumeric",
								"randomBytes",
								"randomDigits.",
							}, false),
							Description: "The type of the secure generator used for generating the change password Id.",
						},
					},
				},
			},
			"passwordless_login_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The time in seconds until a passwordless code is no longer valid and cannot be used by the Passwordless API. Value must be greater than 0.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"registration_verification_id_generator": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The length of the secure generator used for generating the change password Id.",
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"randomAlpha",
								"randomAlphaNumeric",
								"randomBytes",
								"randomDigits.",
							}, false),
							Description: "The type of the secure generator used for generating the change password Id.",
						},
					},
				},
			},
			"registration_verification_id_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The time in seconds until a registration verification Id is no longer valid and cannot be used by the Verify Registration API. Value must be greater than 0.",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"setup_password_id_generator": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The length of the secure generator used for generating the change password Id.",
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"randomAlpha",
								"randomAlphaNumeric",
								"randomBytes",
								"randomDigits.",
							}, false),
							Description: "The type of the secure generator used for generating the change password Id.",
						},
					},
				},
			},
			"setup_password_id_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The time in seconds until a setup password Id is no longer valid and cannot be used by the Change Password API. Value must be greater than 0.",
			},
			"two_factor_id_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The time in seconds until a two factor Id is no longer valid and cannot be used by the Two Factor Login API. Value must be greater than 0.",
			},
			"two_factor_trust_id_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The time in seconds until an issued Two Factor trust Id is no longer valid and the User will be required to complete Two Factor authentication during the next authentication attempt. Value must be greater than 0.",
			},
		},
	}
}

func newEmailConfiguration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"default_from_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default From Name used in sending emails when a from name is not provided on an individual email template. This is the display name part of the email address ( i.e. Jared Dunn <jared@piedpiper.com>).",
			},
			"default_from_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default email address that emails will be sent from when a from address is not provided on an individual email template. This is the address part email address (i.e. Jared Dunn <jared@piedpiper.com>).",
			},
			"forgot_password_email_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Id of the Email Template that is used when a user is sent a forgot password email.",
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The host name of the SMTP server that FusionAuth will use.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "An optional password FusionAuth will use to authenticate with the SMTP server.",
			},
			"passwordless_email_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Id of the Passwordless Email Template.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The port of the SMTP server that FusionAuth will use.",
			},
			"properties": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Additional Email Configuration in a properties file formatted String.",
			},
			"security": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NONE",
					"SSL",
					"TLS",
				}, false),
				Default:     "NONE",
				Description: "The type of security protocol FusionAuth will use when connecting to the SMTP server.",
			},
			"set_password_email_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Id of the Email Template that is used when a user had their account created for them and they must set their password manually and they are sent an email to set their password.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional username FusionAuth will to authenticate with the SMTP server.",
			},
			"verification_email_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The If of the Email Template that is used to send the verification emails to users. These emails are used to verify that a user’s email address is valid. If either the verifyEmail or verifyEmailWhenChanged fields are true this field is required.",
			},
			"verify_email": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the user’s email addresses are verified when the registers with your application.",
			},
			"verify_email_when_changed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the user’s email addresses are verified when the user changes them.",
			},
		},
	}
}

func newPasswordValidationRules() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"breach_detection": {
				Optional: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable Reactor breach detection. Requires an activated license.",
						},
						"match_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"High",
								"Medium",
								"Low",
							}, false),
							Description: "The level of severity where Reactor will consider a breach.",
						},
						"notify_user_email_template_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
							Description:  "The Id of the email template to use when notifying user of breached password. Required if tenant.passwordValidationRules.breachDetection.onLogin is set to NotifyUser.",
						},
						"on_login": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Off",
								"RecordOnly",
								"NotifyUser",
								"RequireChange",
							}, false),
							Description: "The behavior when detecting breaches at time of user login",
						},
					},
				},
			},
			"max_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  256,
			},
			"min_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  8,
			},
			"remember_previous_passwords": {
				Optional: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The number of previous passwords to remember. Value must be greater than 0.",
							ValidateFunc: validation.IntAtLeast(1),
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to prevent a user from using any of their previous passwords.",
						},
					},
				},
			},
			"required_mixed_case": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to force the user to use at least one uppercase and one lowercase character.",
			},
			"require_non_alpha": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to force the user to use at least one non-alphanumeric character.",
			},
			"require_number": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to force the user to use at least one number.",
			},
			"validate_on_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When enabled the user’s password will be validated during login. If the password does not meet the currently configured validation rules the user will be required to change their password.",
			},
		},
	}
}

func buildTentant(data *schema.ResourceData) fusionauth.Tenant {
	return fusionauth.Tenant{
		Data: data.Get("data").(map[string]interface{}),
		EmailConfiguration: fusionauth.EmailConfiguration{
			ForgotPasswordEmailTemplateId: data.Get("email_configuration.0.forgot_password_email_template_id").(string),
			Host:                          data.Get("email_configuration.0.host").(string),
			Password:                      data.Get("email_configuration.0.password").(string),
			PasswordlessEmailTemplateId:   data.Get("email_configuration.0.passwordless_email_template_id").(string),
			Port:                          data.Get("email_configuration.0.port").(int),
			Properties:                    data.Get("email_configuration.0.properties").(string),
			Security:                      fusionauth.EmailSecurityType(data.Get("email_configuration.0.security").(string)),
			SetPasswordEmailTemplateId:    data.Get("email_configuration.0.set_password_email_template_id").(string),
			Username:                      data.Get("email_configuration.0.username").(string),
			VerificationEmailTemplateId:   data.Get("email_configuration.0.verification_email_template_id").(string),
			VerifyEmail:                   data.Get("email_configuration.0.verify_email").(bool),
			VerifyEmailWhenChanged:        data.Get("email_configuration.0.verify_email_when_changed").(bool),
		},
		EventConfiguration: buildEventConfiguration("event_configuration", data),
		ExternalIdentifierConfiguration: fusionauth.ExternalIdentifierConfiguration{
			AuthorizationGrantIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.authorization_grant_id_time_to_live_in_seconds",
			).(int),
			ChangePasswordIdGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.change_password_id_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.change_password_id_generator.0.type").(string),
				),
			},
			ChangePasswordIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.change_password_id_time_to_live_in_seconds",
			).(int),
			DeviceCodeTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.device_code_time_to_live_in_seconds",
			).(int),
			DeviceUserCodeIdGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.device_user_code_id_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.device_user_code_id_generator.0.type").(string),
				),
			},
			EmailVerificationIdGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.email_verification_id_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.email_verification_id_generator.0.type").(string),
				),
			},
			EmailVerificationIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.email_verification_id_time_to_live_in_seconds",
			).(int),
			ExternalAuthenticationIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.external_authentication_id_time_to_live_in_seconds",
			).(int),
			OneTimePasswordTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.one_time_password_time_to_live_in_seconds",
			).(int),
			PasswordlessLoginGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.passwordless_login_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.passwordless_login_generator.0.type").(string),
				),
			},
			PasswordlessLoginTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.passwordless_login_time_to_live_in_seconds").(int),
			RegistrationVerificationIdGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.registration_verification_id_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.registration_verification_id_generator.0.type").(string),
				),
			},
			RegistrationVerificationIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.registration_verification_id_time_to_live_in_seconds",
			).(int),
			SetupPasswordIdGenerator: fusionauth.SecureGeneratorConfiguration{
				Length: data.Get("external_identifier_configuration.0.setup_password_id_generator.0.length").(int),
				Type: fusionauth.SecureGeneratorType(
					data.Get("external_identifier_configuration.0.setup_password_id_generator.0.type").(string),
				),
			},
			SetupPasswordIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.setup_password_id_time_to_live_in_seconds",
			).(int),
			TwoFactorIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.two_factor_id_time_to_live_in_seconds",
			).(int),
			TwoFactorTrustIdTimeToLiveInSeconds: data.Get(
				"external_identifier_configuration.0.two_factor_trust_id_time_to_live_in_seconds",
			).(int),
		},
		FailedAuthenticationConfiguration: fusionauth.FailedAuthenticationConfiguration{
			ActionDuration: int64(data.Get("failed_authentication_configuration.0.action_duration").(int)),
			ActionDurationUnit: fusionauth.ExpiryUnit(
				data.Get("failed_authentication_configuration.0.action_duration_unit").(string),
			),
			ResetCountInSeconds: data.Get("failed_authentication_configuration.0.reset_count_in_seconds").(int),
			TooManyAttempts:     data.Get("failed_authentication_configuration.0.too_many_attempts").(int),
			UserActionId:        data.Get("failed_authentication_configuration.0.user_action_id").(string),
		},
		FamilyConfiguration: fusionauth.FamilyConfiguration{
			AllowChildRegistrations:           data.Get("family_configuration.0.allow_child_registrations").(bool),
			ConfirmChildEmailTemplateId:       data.Get("family_configuration.0.confirm_child_email_template_id").(string),
			DeleteOrphanedAccounts:            data.Get("family_configuration.0.delete_orphaned_accounts").(bool),
			DeleteOrphanedAccountsDays:        data.Get("family_configuration.0.delete_orphaned_accounts_days").(int),
			Enableable:                        buildEnableable("family_configuration.0.enabled", data),
			FamilyRequestEmailTemplateId:      data.Get("family_configuration.0.family_request_email_template_id").(string),
			MaximumChildAge:                   data.Get("family_configuration.0.maximum_child_age").(int),
			MinimumOwnerAge:                   data.Get("family_configuration.0.minimum_owner_age").(int),
			ParentEmailRequired:               data.Get("family_configuration.0.parent_email_required").(bool),
			ParentRegistrationEmailTemplateId: data.Get("family_configuration.0.parent_registration_email_template_id").(string),
		},
		HttpSessionMaxInactiveInterval: data.Get("http_session_max_inactive_interval").(int),
		Issuer:                         data.Get("issuer").(string),
		JwtConfiguration: fusionauth.JWTConfiguration{
			AccessTokenKeyId:             data.Get("jwt_configuration.0.access_token_key_id").(string),
			IdTokenKeyId:                 data.Get("jwt_configuration.0.id_token_key_id").(string),
			RefreshTokenExpirationPolicy: fusionauth.RefreshTokenExpirationPolicy(data.Get("jwt_configuration.0.refresh_token_expiration_policy").(string)),
			RefreshTokenRevocationPolicy: fusionauth.RefreshTokenRevocationPolicy{
				OnLoginPrevented:  data.Get("jwt_configuration.0.refresh_token_revocation_policy_on_login_prevented").(bool),
				OnPasswordChanged: data.Get("jwt_configuration.0.refresh_token_revocation_policy_on_password_change").(bool),
			},
			RefreshTokenTimeToLiveInMinutes: data.Get("jwt_configuration.0.refresh_token_time_to_live_in_minutes").(int),
			RefreshTokenUsagePolicy:         fusionauth.RefreshTokenUsagePolicy(data.Get("jwt_configuration.0.refresh_token_usage_policy").(string)),
			TimeToLiveInSeconds:             data.Get("jwt_configuration.0.time_to_live_in_seconds").(int),
		},
		LogoutURL: data.Get("logout_url").(string),
		MaximumPasswordAge: fusionauth.MaximumPasswordAge{
			Enableable: buildEnableable("maximum_password_age.0.enabled", data),
			Days:       data.Get("maximum_password_age.0.days").(int),
		},
		MinimumPasswordAge: fusionauth.MinimumPasswordAge{
			Enableable: buildEnableable("minimum_password_age.0.enabled", data),
			Seconds:    data.Get("minimum_password_age.0.seconds").(int),
		},
		Name: data.Get("name").(string),
		PasswordEncryptionConfiguration: fusionauth.PasswordEncryptionConfiguration{
			EncryptionScheme:       data.Get("password_encryption_configuration.0.encryption_scheme").(string),
			EncryptionSchemeFactor: data.Get("password_encryption_configuration.0.encryption_scheme_factor").(int),
			ModifyEncryptionSchemeOnLogin: data.Get(
				"password_encryption_configuration.0.modify_encryption_scheme_on_login",
			).(bool),
		},
		PasswordValidationRules: fusionauth.PasswordValidationRules{
			BreachDetection: fusionauth.PasswordBreachDetection{
				Enableable: buildEnableable("password_validation_rules.0.breach_detection.0.enabled", data),
				MatchMode: fusionauth.BreachMatchMode(
					data.Get("password_validation_rules.0.breach_detection.0.match_mode").(string),
				),
				NotifyUserEmailTemplateId: data.Get(
					"password_validation_rules.0.breach_detection.0.notify_user_email_template_id",
				).(string),
				OnLogin: fusionauth.BreachAction(
					data.Get("password_validation_rules.0.breach_detection.0.on_login").(string),
				),
			},
			MaxLength: data.Get("password_validation_rules.0.max_length").(int),
			MinLength: data.Get("password_validation_rules.0.min_length").(int),
			RememberPreviousPasswords: fusionauth.RememberPreviousPasswords{
				Enableable: buildEnableable("password_validation_rules.0.remember_previous_passwords.0.enabled", data),
				Count:      data.Get("password_validation_rules.0.remember_previous_passwords.0.count").(int),
			},
			RequireMixedCase: data.Get("password_validation_rules.0.required_mixed_case").(bool),
			RequireNonAlpha:  data.Get("password_validation_rules.0.require_non_alpha").(bool),
			RequireNumber:    data.Get("password_validation_rules.0.require_number").(bool),
			ValidateOnLogin:  data.Get("password_validation_rules.0.validate_on_login").(bool),
		},
		ThemeId: data.Get("theme_id").(string),
		UserDeletePolicy: fusionauth.TenantUserDeletePolicy{
			Unverified: fusionauth.TimeBasedDeletePolicy{
				Enableable:           buildEnableable("user_delete_policy.0.unverified_enabled", data),
				NumberOfDaysToRetain: data.Get("user_delete_policy.0.unverified_number_of_days_to_retain").(int),
			},
		},
	}
}

func buildEventConfiguration(key string, data *schema.ResourceData) fusionauth.EventConfiguration {
	s := data.Get(key)
	set, ok := s.(*schema.Set)
	if !ok {
		return fusionauth.EventConfiguration{}
	}
	l := set.List()

	ev := make(map[fusionauth.EventType]fusionauth.EventConfigurationData)

	for _, x := range l {
		r := x.(map[string]interface{})
		ev[fusionauth.EventType(r["event"].(string))] = fusionauth.EventConfigurationData{
			TransactionType: fusionauth.TransactionType(r["transaction_type"].(string)),
			Enableable: fusionauth.Enableable{
				Enabled: r["enabled"].(bool),
			},
		}
	}

	return fusionauth.EventConfiguration{Events: ev}
}

func buildResourceDataFromTenant(t fusionauth.Tenant, data *schema.ResourceData) {
	_ = data.Set("data", t.Data)
	_ = data.Set(
		"email_configuration.0.forgot_password_email_template_id",
		t.EmailConfiguration.ForgotPasswordEmailTemplateId,
	)
	_ = data.Set("email_configuration.0.host", t.EmailConfiguration.Host)
	_ = data.Set("email_configuration.0.password", t.EmailConfiguration.Password)
	_ = data.Set("email_configuration.0.passwordless_email_template_id", t.EmailConfiguration.PasswordlessEmailTemplateId)
	_ = data.Set("email_configuration.0.port", t.EmailConfiguration.Port)
	_ = data.Set("email_configuration.0.properties", t.EmailConfiguration.Properties)
	_ = data.Set("email_configuration.0.security", t.EmailConfiguration.Security)
	_ = data.Set("email_configuration.0.set_password_email_template_id", t.EmailConfiguration.SetPasswordEmailTemplateId)
	_ = data.Set("email_configuration.0.username", t.EmailConfiguration.Username)
	_ = data.Set("email_configuration.0.verification_email_template_id", t.EmailConfiguration.VerificationEmailTemplateId)
	_ = data.Set("email_configuration.0.verify_email", t.EmailConfiguration.VerifyEmail)
	_ = data.Set("email_configuration.0.verify_email_when_changed", t.EmailConfiguration.VerifyEmailWhenChanged)
	_ = data.Set(
		"external_identifier_configuration.0.authorization_grant_id_time_to_live_in_seconds",
		t.ExternalIdentifierConfiguration.AuthorizationGrantIdTimeToLiveInSeconds,
	)
	_ = data.Set(
		"external_identifier_configuration.0.change_password_id_generator.0.length",
		t.ExternalIdentifierConfiguration.ChangePasswordIdGenerator.Length,
	)
	_ = data.Set(
		"external_identifier_configuration.0.change_password_id_generator.0.type",
		t.ExternalIdentifierConfiguration.ChangePasswordIdGenerator.Type,
	)
	_ = data.Set(
		"external_identifier_configuration.0.change_password_id_time_to_live_in_seconds",
		t.ExternalIdentifierConfiguration.ChangePasswordIdTimeToLiveInSeconds,
	)
	_ = data.Set(
		"external_identifier_configuration.0.device_code_time_to_live_in_seconds",
		t.ExternalIdentifierConfiguration.DeviceCodeTimeToLiveInSeconds,
	)
	_ = data.Set(
		"external_identifier_configuration.0.device_user_code_id_generator.0.length",
		t.ExternalIdentifierConfiguration.DeviceUserCodeIdGenerator.Length,
	)
	_ = data.Set(
		"external_identifier_configuration.0.device_user_code_id_generator.0.type",
		t.ExternalIdentifierConfiguration.DeviceUserCodeIdGenerator.Type,
	)
	_ = data.Set(
		"external_identifier_configuration.0.email_verification_id_generator.0.length",
		t.ExternalIdentifierConfiguration.EmailVerificationIdGenerator.Length,
	)
	_ = data.Set(
		"external_identifier_configuration.0.email_verification_id_generator.0.type",
		t.ExternalIdentifierConfiguration.EmailVerificationIdGenerator.Type,
	)
	_ = data.Set(
		"external_identifier_configuration.0.email_verification_id_time_to_live_in_seconds",
		t.ExternalIdentifierConfiguration.EmailVerificationIdTimeToLiveInSeconds,
	)
	_ = data.Set(
		"external_identifier_configuration.0.external_authentication_id_time_to_live_in_seconds",
		t.ExternalIdentifierConfiguration.ExternalAuthenticationIdTimeToLiveInSeconds,
	)
	_ = data.Set(
		"external_identifier_configuration.0.one_time_password_time_to_live_in_seconds",
		t.ExternalIdentifierConfiguration.OneTimePasswordTimeToLiveInSeconds,
	)
	_ = data.Set(
		"external_identifier_configuration.0.passwordless_login_generator.0.length",
		t.ExternalIdentifierConfiguration.PasswordlessLoginGenerator.Length,
	)
	_ = data.Set(
		"external_identifier_configuration.0.passwordless_login_generator.0.type",
		t.ExternalIdentifierConfiguration.PasswordlessLoginGenerator.Type,
	)
	_ = data.Set(
		"external_identifier_configuration.0.passwordless_login_time_to_live_in_seconds",
		t.ExternalIdentifierConfiguration.PasswordlessLoginTimeToLiveInSeconds,
	)
	_ = data.Set(
		"external_identifier_configuration.0.registration_verification_id_generator.0.length",
		t.ExternalIdentifierConfiguration.RegistrationVerificationIdTimeToLiveInSeconds,
	)
	_ = data.Set(
		"external_identifier_configuration.0.registration_verification_id_generator.0.type",
		t.ExternalIdentifierConfiguration.RegistrationVerificationIdGenerator.Type,
	)
	_ = data.Set(
		"external_identifier_configuration.0.registration_verification_id_time_to_live_in_seconds",
		t.ExternalIdentifierConfiguration.RegistrationVerificationIdTimeToLiveInSeconds,
	)
	_ = data.Set(
		"external_identifier_configuration.0.setup_password_id_generator.0.length",
		t.ExternalIdentifierConfiguration.SetupPasswordIdGenerator.Length,
	)
	_ = data.Set(
		"external_identifier_configuration.0.setup_password_id_generator.0.type",
		t.ExternalIdentifierConfiguration.SetupPasswordIdGenerator.Type,
	)
	_ = data.Set(
		"external_identifier_configuration.0.setup_password_id_time_to_live_in_seconds",
		t.ExternalIdentifierConfiguration.SetupPasswordIdTimeToLiveInSeconds,
	)
	_ = data.Set(
		"external_identifier_configuration.0.two_factor_id_time_to_live_in_seconds",
		t.ExternalIdentifierConfiguration.TwoFactorIdTimeToLiveInSeconds,
	)
	_ = data.Set(
		"external_identifier_configuration.0.two_factor_trust_id_time_to_live_in_seconds",
		t.ExternalIdentifierConfiguration.TwoFactorTrustIdTimeToLiveInSeconds,
	)
	_ = data.Set(
		"failed_authentication_configuration.0.action_duration", t.FailedAuthenticationConfiguration.ActionDuration,
	)
	_ = data.Set(
		"failed_authentication_configuration.0.action_duration_unit", t.FailedAuthenticationConfiguration.ActionDurationUnit,
	)
	_ = data.Set(
		"failed_authentication_configuration.0.reset_count_in_seconds",
		t.FailedAuthenticationConfiguration.ResetCountInSeconds,
	)
	_ = data.Set(
		"failed_authentication_configuration.0.too_many_attempts", t.FailedAuthenticationConfiguration.TooManyAttempts,
	)
	_ = data.Set("failed_authentication_configuration.0.user_action_id", t.FailedAuthenticationConfiguration.UserActionId)
	_ = data.Set("family_configuration.0.allow_child_registrations", t.FamilyConfiguration.AllowChildRegistrations)
	_ = data.Set(
		"family_configuration.0.confirm_child_email_template_id",
		t.FamilyConfiguration.ConfirmChildEmailTemplateId,
	)
	_ = data.Set("family_configuration.0.delete_orphaned_accounts", t.FamilyConfiguration.DeleteOrphanedAccounts)
	_ = data.Set("family_configuration.0.delete_orphaned_accounts_days", t.FamilyConfiguration.DeleteOrphanedAccountsDays)
	_ = data.Set("family_configuration.0.enabled", t.FamilyConfiguration.Enabled)
	_ = data.Set(
		"family_configuration.0.family_request_email_template_id",
		t.FamilyConfiguration.FamilyRequestEmailTemplateId,
	)
	_ = data.Set("family_configuration.0.maximum_child_age", t.FamilyConfiguration.MaximumChildAge)
	_ = data.Set("family_configuration.0.minimum_owner_age", t.FamilyConfiguration.MinimumOwnerAge)
	_ = data.Set("family_configuration.0.parent_email_required", t.FamilyConfiguration.ParentEmailRequired)
	_ = data.Set(
		"family_configuration.0.parent_registration_email_template_id",
		t.FamilyConfiguration.ParentRegistrationEmailTemplateId,
	)
	_ = data.Set("http_session_max_inactive_interval", t.HttpSessionMaxInactiveInterval)
	_ = data.Set("issuer", t.Issuer)
	_ = data.Set("jwt_configuration.0.access_token_key_id", t.JwtConfiguration.AccessTokenKeyId)
	_ = data.Set("jwt_configuration.0.id_token_key_id", t.JwtConfiguration.IdTokenKeyId)
	_ = data.Set("jwt_configuration.0.refresh_token_expiration_policy", t.JwtConfiguration.RefreshTokenExpirationPolicy)
	_ = data.Set("jwt_configuration.0.refresh_token_revocation_policy_on_login_prevented", t.JwtConfiguration.RefreshTokenRevocationPolicy.OnLoginPrevented)
	_ = data.Set("jwt_configuration.0.refresh_token_revocation_policy_on_password_change", t.JwtConfiguration.RefreshTokenRevocationPolicy.OnPasswordChanged)
	_ = data.Set(
		"jwt_configuration.0.refresh_token_usage_policy",
		t.JwtConfiguration.RefreshTokenUsagePolicy,
	)
	_ = data.Set("jwt_configuration.0.id_token_key_id", t.JwtConfiguration.IdTokenKeyId)
	_ = data.Set("jwt_configuration.0.time_to_live_in_seconds", t.JwtConfiguration.TimeToLiveInSeconds)
	_ = data.Set("logout_url", t.LogoutURL)
	_ = data.Set("maximum_password_age.0.enabled", t.MaximumPasswordAge.Enabled)
	_ = data.Set("maximum_password_age.0.days", t.MaximumPasswordAge.Days)
	_ = data.Set("minimum_password_age.0.enabled", t.MinimumPasswordAge.Enabled)
	_ = data.Set("minimum_password_age.0.seconds", t.MinimumPasswordAge.Seconds)
	_ = data.Set("name", t.Name)
	_ = data.Set(
		"password_encryption_configuration.0.encryption_scheme",
		t.PasswordEncryptionConfiguration.EncryptionScheme,
	)
	_ = data.Set(
		"password_encryption_configuration.0.encryption_scheme_factor",
		t.PasswordEncryptionConfiguration.EncryptionSchemeFactor,
	)
	_ = data.Set(
		"password_encryption_configuration.0.modify_encryption_scheme_on_login",
		t.PasswordEncryptionConfiguration.ModifyEncryptionSchemeOnLogin,
	)
	_ = data.Set(
		"password_validation_rules.0.breach_detection.0.enabled",
		t.PasswordValidationRules.BreachDetection.Enabled)
	_ = data.Set(
		"password_validation_rules.0.breach_detection.0.match_mode",
		t.PasswordValidationRules.BreachDetection.MatchMode,
	)
	_ = data.Set(
		"password_validation_rules.0.breach_detection.0.notify_user_email_template_id",
		t.PasswordValidationRules.BreachDetection.NotifyUserEmailTemplateId,
	)
	_ = data.Set(
		"password_validation_rules.0.breach_detection.0.on_login", t.PasswordValidationRules.BreachDetection.OnLogin,
	)
	_ = data.Set("password_validation_rules.0.max_length", t.PasswordValidationRules.MaxLength)
	_ = data.Set("password_validation_rules.0.min_length", t.PasswordValidationRules.MinLength)
	_ = data.Set(
		"password_validation_rules.0.remember_previous_passwords.0.enabled",
		t.PasswordValidationRules.RememberPreviousPasswords.Enabled,
	)
	_ = data.Set(
		"password_validation_rules.0.remember_previous_passwords.0.count",
		t.PasswordValidationRules.RememberPreviousPasswords.Count,
	)
	_ = data.Set("password_validation_rules.0.required_mixed_case", t.PasswordValidationRules.RequireMixedCase)
	_ = data.Set("password_validation_rules.0.require_non_alpha", t.PasswordValidationRules.RequireNonAlpha)
	_ = data.Set("password_validation_rules.0.require_number", t.PasswordValidationRules.RequireNumber)
	_ = data.Set("password_validation_rules.0.validate_on_login", t.PasswordValidationRules.ValidateOnLogin)
	_ = data.Set("theme_id", t.ThemeId)
	_ = data.Set("user_delete_policy.0.unverified_enabled", t.UserDeletePolicy.Unverified.Enabled)
	_ = data.Set(
		"user_delete_policy.0.unverified_number_of_days_to_retain",
		t.UserDeletePolicy.Unverified.NumberOfDaysToRetain,
	)

	e := make([]map[string]interface{}, 0, len(t.EventConfiguration.Events))
	for k, v := range t.EventConfiguration.Events {
		e = append(e, map[string]interface{}{
			"event":            k,
			"transaction_type": v.TransactionType,
			"enabled":          v.Enabled,
		})
	}
	_ = data.Set("event_configuration", e)
}

func createTenant(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	t := fusionauth.TenantRequest{
		Tenant:         buildTentant(data),
		SourceTenantId: data.Get("source_tentant_id").(string),
	}

	resp, faErrs, err := client.FAClient.CreateTenant("", t)

	if err != nil {
		return fmt.Errorf("CreateTenant err: %v", err)
	}

	if faErrs != nil {
		return fmt.Errorf("CreateTenant errors: %v", faErrs)
	}

	data.SetId(resp.Tenant.Id)
	return nil
}

func readTenant(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveTenant(id)
	if err != nil {
		return err
	}
	if faErrs != nil {
		return fmt.Errorf("RetrieveTenant errors: %v", faErrs)
	}
	buildResourceDataFromTenant(resp.Tenant, data)

	return nil
}

func updateTenant(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)

	t := fusionauth.TenantRequest{
		Tenant:         buildTentant(data),
		SourceTenantId: data.Get("source_tentant_id").(string),
	}

	_, faErrs, err := client.FAClient.UpdateTenant(data.Id(), t)

	if err != nil {
		return fmt.Errorf("UpdateTenant err: %v", err)
	}

	if faErrs != nil {
		return fmt.Errorf("UpdateTenant errors: %v", faErrs)
	}

	return nil
}

func deleteTenant(data *schema.ResourceData, i interface{}) error {
	client := i.(Client)
	_, faErrs, err := client.FAClient.DeleteTenant(data.Id())
	if err != nil {
		return err
	}

	if faErrs != nil {
		return fmt.Errorf("DeleteTenant errors: %v", faErrs)
	}

	return nil
}
