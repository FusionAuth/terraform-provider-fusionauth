package fusionauth

import (
	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newTenant() *schema.Resource {
	return &schema.Resource{
		CreateContext: createTenant,
		ReadContext:   readTenant,
		UpdateContext: updateTenant,
		DeleteContext: deleteTenant,
		Schema: map[string]*schema.Schema{
			"source_tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The optional Id of an existing Tenant to make a copy of. If present, the tenant.id and tenant.name values of the request body will be applied to the new Tenant, all other values will be copied from the source Tenant to the new Tenant.",
				ValidateFunc: validation.IsUUID,
			},
			"webhook_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An array of Webhook Ids. For Webhooks that are not already configured for All Tenants, specifying an Id on this request will indicate the associated Webhook should handle events for this tenant.",
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The Id to use for the new Tenant. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
			},
			"access_control_configuration": {
				Type:       schema.TypeList,
				MaxItems:   1,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ui_ip_access_control_list_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
							Description:  "The Id of the IP Access Control List limiting access to all applications in this tenant.",
						},
					},
				},
			},
			"captcha_configuration": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether captcha configuration is enabled.",
						},
						"captcha_method": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"GoogleRecaptchaV2",
								"GoogleRecaptchaV3",
								"HCaptcha",
								"HCaptchaEnterprise",
							}, false),
							Description: "The type of captcha method to use. This field is required when tenant.captchaConfiguration.enabled is set to true.",
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The secret key for this captcha method. This field is required when tenant.captchaConfiguration.enabled is set to true.",
						},
						"site_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The site key for this captcha method. This field is required when tenant.captchaConfiguration.enabled is set to true.",
						},
						"threshold": {
							Type:         schema.TypeFloat,
							Optional:     true,
							Default:      0.5,
							Description:  "The numeric threshold which separates a passing score from a failing one. This value only applies if using either the Google v3 or HCaptcha Enterprise method, otherwise this value is ignored.",
							ValidateFunc: validation.FloatBetween(0.0, 1.0),
						},
					},
				},
			},
			"connector_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "A list of Connector policies. Users will be authenticated against Connectors in order. Each Connector can be included in this list at most once and must exist. If not specified a policy for the default connector will be created",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connector_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The identifier of the Connector to which this policy refers.",
						},
						"domains": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "A list of email domains to which this connector should apply. A value of [\"*\"] indicates this connector applies to all users.",
						},
						"migrate": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If true, the user’s data will be migrated to FusionAuth at first successful authentication",
						},
					},
				},
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
								"audit-log.create",
								"event-log.create",
								"jwt.public-key.update",
								"jwt.refresh",
								"jwt.refresh-token.revoke",
								"kickstart.success",
								"group.create",
								"group.create.complete",
								"group.delete",
								"group.delete.complete",
								"group.member.add",
								"group.member.add.complete",
								"group.member.remove",
								"group.member.remove.complete",
								"group.member.update",
								"group.member.update.complete",
								"group.update",
								"group.update.complete",
								"user.action",
								"user.bulk.create",
								"user.create",
								"user.create.complete",
								"user.deactivate",
								"user.delete",
								"user.delete.complete",
								"user.email.update",
								"user.email.verified",
								"user.identity-provider.link",
								"user.identity-provider.unlink",
								"user.loginId.duplicate.create",
								"user.loginId.duplicate.update",
								"user.login.failed",
								"user.login.new-device",
								"user.login.success",
								"user.login.suspicious",
								"user.password.breach",
								"user.password.reset.send",
								"user.password.reset.start",
								"user.password.reset.success",
								"user.password.update",
								"user.reactivate",
								"user.registration.create",
								"user.registration.create.complete",
								"user.registration.delete",
								"user.registration.delete.complete",
								"user.registration.update",
								"user.registration.update.complete",
								"user.registration.verified",
								"user.two-factor.method.add",
								"user.two-factor.method.remove",
								"user.update",
								"user.update.complete",
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
				Computed: true,
				Elem:     newFailedAuthenticationConfiguration(),
			},
			"family_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem:     newFamilyConfiguration(),
			},
			"form_configuration": {
				Type:       schema.TypeList,
				MaxItems:   1,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_user_form_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IsUUID,
							Description:  "The unique Id of the form to use for the Add and Edit User form when used in the FusionAuth admin UI.",
						},
					},
				},
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
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IsUUID,
							Description:  "The unique id of the signing key used to sign the access token.",
						},
						"id_token_key_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IsUUID,
							Description:  "The unique id of the signing key used to sign the Id token.",
						},
						"refresh_token_expiration_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Fixed",
							ValidateFunc: validation.StringInSlice([]string{
								"Fixed",
								"SlidingWindow",
								"SlidingWindowWithMaximumLifetime",
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
						"refresh_token_sliding_window_maximum_time_to_live_in_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The maximum lifetime of a refresh token when using a refresh token expiration policy of SlidingWindowWithMaximumLifetime. Value must be greater than 0.",
							ValidateFunc: validation.IntAtLeast(1),
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
							Default:  "Reusable",
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
			"login_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"require_authentication": {
							Type:        schema.TypeBool,
							Computed:    true,
							Optional:    true,
							Description: "Indicates whether to require an API key for the Login API when an applicationId is not provided. When an applicationId is provided to the Login API call, the application configuration will take precedence. In almost all cases, you will want to this to be true.",
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
				Computed: true,
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
				Computed: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     30,
							Description: "The password minimum age in seconds. When enabled FusionAuth will not allow a password to be changed until it reaches this minimum age. Required when systemConfiguration.minimumPasswordAge.enabled is set to true.",
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
			"multi_factor_configuration": {
				Type:       schema.TypeList,
				MaxItems:   1,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"login_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "Enabled",
							ValidateFunc: validation.StringInSlice([]string{fusionauth.MultiFactorLoginPolicy_Enabled.String(), fusionauth.MultiFactorLoginPolicy_Disabled.String(), fusionauth.MultiFactorLoginPolicy_Required.String()}, false),
							Description:  "When set to Enabled and a user has one or more two-factor methods configured, the user will be required to complete a two-factor challenge during login. When set to Disabled, even when a user has configured one or more two-factor methods, the user will not be required to complete a two-factor challenge during login. When the login policy is to Required, a two-factor challenge will be required during login. If a user does not have configured two-factor methods, they will not be able to log in.",
						},
						"authenticator": {
							Type:       schema.TypeList,
							MaxItems:   1,
							Optional:   true,
							Computed:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "When enabled, users may utilize an authenticator application to complete a multi-factor authentication request. This method uses TOTP (Time-Based One-Time Password) as defined in RFC 6238 and often uses an native mobile app such as Google Authenticator.",
									},
								},
							},
						},
						"email": {
							Type:       schema.TypeList,
							MaxItems:   1,
							Optional:   true,
							Computed:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "When enabled, users may utilize an email address to complete a multi-factor authentication request.",
									},
									"template_id": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsUUID,
										Description:  "The Id of the email template that is used when notifying a user to complete a multi-factor authentication request.",
									},
								},
							},
						},
						"sms": {
							Type:       schema.TypeList,
							MaxItems:   1,
							Optional:   true,
							Computed:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "When enabled, users may utilize a mobile phone number to complete a multi-factor authentication request.",
									},
									"messenger_id": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsUUID,
										Description:  "The messenger that is used to deliver a SMS multi-factor authentication request.",
									},
									"template_id": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsUUID,
										Description:  "The Id of the SMS template that is used when notifying a user to complete a multi-factor authentication request.",
									},
								},
							},
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Tenant.",
			},
			"oauth_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_credentials_access_token_populate_lambda_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
							Description:  "The Id of a lambda that will be called to populate the JWT during a client credentials grant.",
						},
					},
				},
			},
			"password_encryption_configuration": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				Computed: true,
				Elem:     newPasswordValidationRules(),
			},
			"rate_limit_configuration": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"failed_login": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether rate limiting is enabled for failed login.",
									},
									"limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      5,
										Description:  "The number of times a user can fail to login within the configured timePeriodInSeconds duration. If a Failed authentication action has been configured then it will take precedence. Value must be greater than 0.",
										ValidateFunc: validation.IntAtLeast(1),
									},
									"time_period_in_seconds": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      60,
										Description:  "The duration for the number of times a user can fail login before being rate limited. Value must be greater than 0.",
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
						"forgot_password": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether rate limiting is enabled for forgot password.",
									},
									"limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      5,
										Description:  "The number of times a user can request a forgot password email within the configured `time_period_in_seconds` duration. Value must be greater than 0.",
										ValidateFunc: validation.IntAtLeast(1),
									},
									"time_period_in_seconds": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      60,
										Description:  "The duration for the number of times a user can request a forgot password email before being rate limited. Value must be greater than 0.",
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
						"send_email_verification": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether rate limiting is enabled for send email verification.",
									},
									"limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      5,
										Description:  "The number of times a user can request a verification email within the configured `time_period_in_seconds` duration. Value must be greater than 0.",
										ValidateFunc: validation.IntAtLeast(1),
									},
									"time_period_in_seconds": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      60,
										Description:  "The duration for the number of times a user can request a verification email before being rate limited. Value must be greater than 0.",
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
						"send_passwordless": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether rate limiting is enabled for send passwordless.",
									},
									"limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      5,
										Description:  "The number of times a user can request a passwordless login email within the configured `time_period_in_seconds` duration. Value must be greater than 0.",
										ValidateFunc: validation.IntAtLeast(1),
									},
									"time_period_in_seconds": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      60,
										Description:  "The duration for the number of times a user can request a passwordless login email before being rate limited. Value must be greater than 0.",
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
						"send_registration_verification": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether rate limiting is enabled for send registration verification.",
									},
									"limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      5,
										Description:  "The number of times a user can request a registration verification email within the configured `time_period_in_seconds` duration. Value must be greater than 0.",
										ValidateFunc: validation.IntAtLeast(1),
									},
									"time_period_in_seconds": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      60,
										Description:  "The duration for the number of times a user can request a registration verification email before being rate limited. Value must be greater than 0.",
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
						"send_two_factor": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether rate limiting is enabled for send two factor.",
									},
									"limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      5,
										Description:  "The number of times a user can request a two-factor code by email or SMS within the configured `time_period_in_seconds` duration. Value must be greater than 0.",
										ValidateFunc: validation.IntAtLeast(1),
									},
									"time_period_in_seconds": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      60,
										Description:  "The duration for the number of times a user can request a two-factor code by email or SMS before being rate limited. Value must be greater than 0.",
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
					},
				},
			},
			"registration_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem:     newTenantRegistrationConfiguration(),
			},
			"theme_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The unique Id of the theme to be used to style the login page and other end user templates.",
				ValidateFunc: validation.IsUUID,
			},
			"user_delete_policy": {
				Optional: true,
				Computed: true,
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
							Default:      120,
							Description:  "The number of days from creation users will be retained before being deleted for not completing email verification. This field is required when tenant.userDeletePolicy.unverified.enabled is set to true. Value must be greater than 0.",
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},
			"username_configuration": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique": {
							Optional: true,
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "When true, FusionAuth will handle username collisions by generating a random suffix.",
									},
									"number_of_digits": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      5,
										Description:  "The maximum number of digits to use when building a unique suffix for a username. A number will be randomly selected and will be 1 or more digits up to this configured value in length. For example, if this value is 5, the suffix will be a number between 00001 and 99999, inclusive.",
										ValidateFunc: validation.IntInSlice([]int{3, 4, 5, 6, 7, 8, 9, 10}),
									},
									"separator": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "#",
										Description: "A single character to use as a separator from the requested username and a unique suffix that is added when a duplicate username is detected. This value can be a single non-alphanumeric ASCII character.",
									},
									"strategy": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "OnCollision",
										ValidateFunc: validation.StringInSlice([]string{
											"Always",
											"OnCollision",
										}, false),
										Description: "This strategy instructions FusionAuth when to append a unique suffix to the username. ",
									},
								},
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func newTenantRegistrationConfiguration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"blocked_domains": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A list of unique domains that are not allowed to register when self service is enabled.",
			},
		},
	}
}

func newFailedAuthenticationConfiguration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
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
			"action_cancel_policy_on_password_reset": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether you want the user to be able to self-service unlock their account prior to the action duration by completing a password reset workflow.",
			},
			"email_user": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates you would like to email the user when the user’s account is locked due to this action being taken. This requires the User Action specified by the tenant.failedAuthenticationConfiguration.userActionId to also be configured for email. If the User Action is not configured to be able to email the user, this configuration will be ignored.",
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(secureGeneratorTypes(), false),
							Description:  "The type of the secure generator used for generating the change password Id.",
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(secureGeneratorTypes(), false),
							Description:  "The type of the secure generator used for generating the change password Id.",
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(secureGeneratorTypes(), false),
							Description:  "The type of the secure generator used for generating the change password Id.",
						},
					},
				},
			},
			"email_verification_one_time_code_generator": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The length of the secure generator used for generating the email verification one time code.",
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      string(fusionauth.SecureGeneratorType_RandomAlphaNumeric),
							ValidateFunc: validation.StringInSlice(secureGeneratorTypes(), false),
							Description:  "The type of the secure generator used for generating the email verification one time code.",
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
			"login_intent_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The number of seconds before the Login Timeout identifier is no longer valid to complete post-authentication steps in the OAuth workflow. Must be greater than 0.",
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(secureGeneratorTypes(), false),
							Description:  "The type of the secure generator used for generating the change password Id.",
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(secureGeneratorTypes(), false),
							Description:  "The type of the secure generator used for generating the change password Id.",
						},
					},
				},
			},
			"registration_verification_one_time_code_generator": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The length of the secure generator used for generating the registration verification one time code.",
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      string(fusionauth.SecureGeneratorType_RandomAlphaNumeric),
							ValidateFunc: validation.StringInSlice(secureGeneratorTypes(), false),
							Description:  "The type of the secure generator used for generating the registration verification one time code.",
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
			"saml_v2_authn_request_id_ttl_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "The time in seconds that a SAML AuthN request will be eligible for use to authenticate with FusionAuth.",
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(secureGeneratorTypes(), false),
							Description:  "The type of the secure generator used for generating the change password Id.",
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
			"trust_token_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      180,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The number of seconds before the Trust Token is no longer valid to complete a request that requires trust. Value must be greater than 0.",
			},
			"pending_account_link_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3600,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The number of seconds before the pending account link identifier is no longer valid to complete an account link request. Value must be greater than 0.",
			},
			"two_factor_one_time_code_id_generator": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The length of the secure generator used for generating the the two factor code Id.",
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      string(fusionauth.SecureGeneratorType_RandomAlphaNumeric),
							ValidateFunc: validation.StringInSlice(secureGeneratorTypes(), false),
							Description:  "The type of the secure generator used for generating the two factor one time code Id.",
						},
					},
				},
			},
			"two_factor_id_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The time in seconds until a two factor Id is no longer valid and cannot be used by the Two Factor Login API. Value must be greater than 0.",
			},
			"two_factor_one_time_code_id_time_to_live_in_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The number of seconds before the Two-Factor One Time Code used to enable or disable a two-factor method is no longer valid. Must be greater than 0.",
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
			"additional_headers": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The additional SMTP headers to be added to each outgoing email. Each SMTP header consists of a name and a value.",
			},
			"default_from_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default From Name used in sending emails when a from name is not provided on an individual email template. This is the display name part of the email address ( i.e. Jared Dunn <jared@piedpiper.com>).",
			},
			"default_from_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true, // Fusionauth defaults to `change-me@example.com` if not configured.
				Description: "The default email address that emails will be sent from when a from address is not provided on an individual email template. This is the address part email address (i.e. Jared Dunn <jared@piedpiper.com>).",
			},
			// "email_verification_email_template_id": {
			// 	Type:         schema.TypeString,
			// 	Optional:     true,
			// 	Description:  "The Id of the Email Template used to send emails to users to verify that their email address is valid.",
			// 	ValidateFunc: validation.IsUUID,
			// },
			"email_update_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template that is used when a user is sent a forgot password email.",
				ValidateFunc: validation.IsUUID,
			},
			"email_verified_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template used to verify user emails.",
				ValidateFunc: validation.IsUUID,
			},
			"forgot_password_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template that is used when a user is sent a forgot password email.",
				ValidateFunc: validation.IsUUID,
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The host name of the SMTP server that FusionAuth will use.",
			},
			"implicit_email_verification_allowed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When set to true, this allows email to be verified as a result of completing a similar email based workflow such as change password. When set to false, the user must explicitly complete the email verification workflow even if the user has already completed a similar email workflow such as change password.",
			},
			"login_id_in_use_on_create_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template used to send emails to users when another user attempts to create an account with their login Id.",
				ValidateFunc: validation.IsUUID,
			},
			"login_id_in_use_on_update_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template used to send emails to users when another user attempts to create an account with their login Id.",
				ValidateFunc: validation.IsUUID,
			},
			"login_new_device_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template used to send emails to users when they log in on a new device.",
				ValidateFunc: validation.IsUUID,
			},
			"login_suspicious_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template used to send emails to users when a suspicious login occurs.",
				ValidateFunc: validation.IsUUID,
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "An optional password FusionAuth will use to authenticate with the SMTP server.",
			},
			"passwordless_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Passwordless Email Template.",
				ValidateFunc: validation.IsUUID,
			},
			"password_reset_success_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template used to send emails to users when they have completed a 'forgot password' workflow and their password has been reset.",
				ValidateFunc: validation.IsUUID,
			},
			"password_update_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template used to send emails to users when they have completed a 'forgot password' workflow and their password has been reset.",
				ValidateFunc: validation.IsUUID,
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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template that is used when a user had their account created for them and they must set their password manually and they are sent an email to set their password.",
				ValidateFunc: validation.IsUUID,
			},
			"two_factor_method_add_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template used to send emails to users when a MFA method has been added to their account.",
				ValidateFunc: validation.IsUUID,
			},
			"two_factor_method_remove_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template used to send emails to users when a MFA method has been removed from their account.",
				ValidateFunc: validation.IsUUID,
			},
			"unverified": {
				Optional:   true,
				Type:       schema.TypeList,
				MaxItems:   1,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_email_change_when_gated": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "When this value is set to true, the user is allowed to change their email address when they are gated because they haven’t verified their email address.",
						},
						"behavior": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Allow",
								"Gated",
							}, false),
							Default:     "Allow",
							Description: "The behavior when detecting breaches at time of user login",
						},
					},
				},
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional username FusionAuth will to authenticate with the SMTP server.",
			},
			"verification_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template that is used to send the verification emails to users. These emails are used to verify that a user’s email address is valid. If either the verifyEmail or verifyEmailWhenChanged fields are true this field is required.",
				ValidateFunc: validation.IsUUID,
			},
			"verification_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ClickableLink",
				Description:  "The process by which the user will verify their email address.",
				ValidateFunc: validation.StringInSlice([]string{"ClickableLink", "FormField"}, false),
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
				Optional:   true,
				Type:       schema.TypeList,
				MaxItems:   1,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
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
							Default:      1,
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

// secureGeneratorTypes returns a list of the valid secure generator types.
func secureGeneratorTypes() []string {
	return []string{
		string(fusionauth.SecureGeneratorType_RandomAlpha),
		string(fusionauth.SecureGeneratorType_RandomAlphaNumeric),
		string(fusionauth.SecureGeneratorType_RandomBytes),
		string(fusionauth.SecureGeneratorType_RandomDigits),
	}
}
