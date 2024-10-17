package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSystemConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: createSystemConfiguration,
		ReadContext:   readSystemConfiguration,
		UpdateContext: updateSystemConfiguration,
		DeleteContext: deleteSystemConfiguration,
		Schema: map[string]*schema.Schema{
			"audit_log_configuration": {
				Type:       schema.TypeList,
				MaxItems:   1,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete": {
							Type:       schema.TypeList,
							MaxItems:   1,
							Optional:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not FusionAuth should delete the Audit Log based upon this configuration. When true the auditLogConfiguration.delete.numberOfDaysToRetain will be used to identify audit logs that are eligible for deletion. When this value is set to false audit logs will be preserved forever.",
									},
									"number_of_days_to_retain": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     365,
										Description: "The number of days to retain the Audit Log.",
									},
								},
							},
						},
					},
				},
			},
			"cors_configuration": {
				Type:       schema.TypeList,
				MaxItems:   1,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_credentials": {
							Type:        schema.TypeBool,
							Default:     false,
							Optional:    true,
							Description: "The Access-Control-Allow-Credentials response header values as described by MDN Access-Control-Allow-Credentials.",
						},
						"allowed_headers": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "The Access-Control-Allow-Headers response header values as described by MDN Access-Control-Allow-Headers.",
							Computed:    true,
						},
						"allowed_methods": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "The Access-Control-Allow-Methods response header values as described by MDN Access-Control-Allow-Methods.",
							Computed:    true,
						},
						"allowed_origins": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "The Access-Control-Allow-Origin response header values as described by MDN Access-Control-Allow-Origin. If the wildcard * is specified, no additional domains may be specified.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Default:     true,
							Optional:    true,
							Description: "Whether the FusionAuth CORS filter will process requests made to FusionAuth.",
						},
						"exposed_headers": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "The Access-Control-Expose-Headers response header values as described by MDN Access-Control-Expose-Headers.",
							Computed:    true,
						},
						"preflight_max_age_in_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1800,
							Description: "The Access-Control-Max-Age response header values as described by MDN Access-Control-Max-Age.",
						},
					},
				},
			},
			"event_log_configuration": {
				Type:       schema.TypeList,
				MaxItems:   1,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number_to_retain": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     10000,
							Description: "The number of events to retain. Once the the number of event logs exceeds this configured value they will be deleted starting with the oldest event logs.",
						},
					},
				},
			},
			"login_record_configuration": {
				Type:       schema.TypeList,
				MaxItems:   1,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete": {
							Type:       schema.TypeList,
							MaxItems:   1,
							Optional:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not FusionAuth should delete the login records based upon this configuration. When true the loginRecordConfiguration.delete.numberOfDaysToRetain will be used to identify login records that are eligible for deletion. When this value is set to false login records will be preserved forever.",
									},
									"number_of_days_to_retain": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     365,
										Description: "The number of days to retain login records.",
									},
								},
							},
						},
					},
				},
			},
			"report_timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "America/Denver",
				Description: "The time zone used to adjust the stored UTC time when generating reports. Since reports are usually rolled up hourly, this timezone will be used for demarcating the hours.",
			},
			"ui_configuration": {
				Type:       schema.TypeList,
				MaxItems:   1,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A hexadecimal color to override the default menu color in the user interface.",
						},
						"logo_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A URL of a logo to override the default FusionAuth logo in the user interface.",
						},
						"menu_font_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A hexadecimal color to override the default menu font color in the user interface.",
						},
					},
				},
			},
			"webhook_event_log_configuration": {
				Type:       schema.TypeList,
				MaxItems:   1,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete": {
							Type:       schema.TypeList,
							MaxItems:   1,
							Optional:   true,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not FusionAuth should delete the webhook event logs based upon this configuration. When true the webhookEventLogConfiguration.delete.numberOfDaysToRetain will be used to identify webhook event logs that are eligible for deletion. When this value is set to false webhook event logs will be preserved forever.",
									},
									"number_of_days_to_retain": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     365,
										Description: "The number of days to retain webhook event logs.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func createSystemConfiguration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	data.SetId("syscfg")
	return updateSysCfg(buildSystemConfigurationRequest(data), client)
}

func readSystemConfiguration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	resp, err := client.FAClient.RetrieveSystemConfiguration()
	if err != nil {
		return diag.Errorf("RetrieveSystemConfiguration err: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return diag.FromErr(err)
	}

	return buildResourceFromSystemConfiguration(resp.SystemConfiguration, data)
}

func updateSystemConfiguration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	return updateSysCfg(buildSystemConfigurationRequest(data), client)
}

func deleteSystemConfiguration(_ context.Context, _ *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	return updateSysCfg(getDefaultSystemConfigurationRequest(), client)
}

func updateSysCfg(req fusionauth.SystemConfigurationRequest, client Client) diag.Diagnostics {
	resp, faErrs, err := client.FAClient.UpdateSystemConfiguration(req)
	if err != nil {
		return diag.Errorf("UpdateSystemConfiguration err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func buildSystemConfigurationRequest(data *schema.ResourceData) fusionauth.SystemConfigurationRequest {
	sc := getDefaultSystemConfigurationRequest()
	if v, ok := data.GetOk("audit_log_configuration.0.delete.0.enabled"); ok {
		sc.SystemConfiguration.AuditLogConfiguration.Delete.Enabled = v.(bool)
	}
	if v, ok := data.GetOk("audit_log_configuration.0.delete.0.number_of_days_to_retain"); ok {
		sc.SystemConfiguration.AuditLogConfiguration.Delete.NumberOfDaysToRetain = v.(int)
	}

	if v, ok := data.GetOk("cors_configuration.0.allow_credentials"); ok {
		sc.SystemConfiguration.CorsConfiguration.AllowCredentials = v.(bool)
	}

	if _, ok := data.GetOk("cors_configuration.0.allowed_headers"); ok {
		sc.SystemConfiguration.CorsConfiguration.AllowedHeaders = handleStringSlice("cors_configuration.0.allowed_headers", data)
	}

	if v, ok := data.GetOk("cors_configuration.0.allowed_methods"); ok {
		set := v.(*schema.Set)
		l := set.List()
		s := make([]fusionauth.HTTPMethod, 0, len(l))
		for _, x := range l {
			s = append(s, fusionauth.HTTPMethod(x.(string)))
		}
		sc.SystemConfiguration.CorsConfiguration.AllowedMethods = s
	}

	if _, ok := data.GetOk("cors_configuration.0.allowed_origins"); ok {
		sc.SystemConfiguration.CorsConfiguration.AllowedOrigins = handleStringSlice("cors_configuration.0.allowed_origins", data)
	}

	if v, ok := data.GetOk("cors_configuration.0.enabled"); ok {
		sc.SystemConfiguration.CorsConfiguration.Enabled = v.(bool)
	}

	if _, ok := data.GetOk("cors_configuration.0.exposed_headers"); ok {
		sc.SystemConfiguration.CorsConfiguration.ExposedHeaders = handleStringSlice("cors_configuration.0.exposed_headers", data)
	}

	if v, ok := data.GetOk("cors_configuration.0.preflight_max_age_in_seconds"); ok {
		sc.SystemConfiguration.CorsConfiguration.PreflightMaxAgeInSeconds = v.(int)
	}

	if v, ok := data.GetOk("event_log_configuration.0.number_to_retain"); ok {
		sc.SystemConfiguration.EventLogConfiguration.NumberToRetain = v.(int)
	}

	if v, ok := data.GetOk("login_record_configuration.0.delete.0.enabled"); ok {
		sc.SystemConfiguration.LoginRecordConfiguration.Delete.Enabled = v.(bool)
	}
	if v, ok := data.GetOk("login_record_configuration.0.delete.0.number_of_days_to_retain"); ok {
		sc.SystemConfiguration.LoginRecordConfiguration.Delete.NumberOfDaysToRetain = v.(int)
	}

	sc.SystemConfiguration.ReportTimezone = data.Get("report_timezone").(string)

	if v, ok := data.GetOk("ui_configuration.0.header_color"); ok {
		sc.SystemConfiguration.UiConfiguration.HeaderColor = v.(string)
	}

	if v, ok := data.GetOk("ui_configuration.0.logo_url"); ok {
		sc.SystemConfiguration.UiConfiguration.LogoURL = v.(string)
	}

	if v, ok := data.GetOk("ui_configuration.0.menu_font_color"); ok {
		sc.SystemConfiguration.UiConfiguration.MenuFontColor = v.(string)
	}

	if v, ok := data.GetOk("webhook_event_log_configuration.0.delete.0.enabled"); ok {
		sc.SystemConfiguration.WebhookEventLogConfiguration.Delete.Enabled = v.(bool)
	}
	if v, ok := data.GetOk("webhook_event_log_configuration.0.delete.0.number_of_days_to_retain"); ok {
		sc.SystemConfiguration.WebhookEventLogConfiguration.Delete.NumberOfDaysToRetain = v.(int)
	}

	return sc
}

func buildResourceFromSystemConfiguration(sc fusionauth.SystemConfiguration, data *schema.ResourceData) diag.Diagnostics {
	err := data.Set("audit_log_configuration", []map[string]interface{}{
		{
			"delete": []map[string]interface{}{
				{
					"enabled":                  sc.AuditLogConfiguration.Delete.Enabled,
					"number_of_days_to_retain": sc.AuditLogConfiguration.Delete.NumberOfDaysToRetain,
				},
			},
		},
	})
	if err != nil {
		return diag.Errorf("system_configuration.audit_log_configuration: %s", err.Error())
	}

	err = data.Set("cors_configuration", []map[string]interface{}{
		{
			"allow_credentials":            sc.CorsConfiguration.AllowCredentials,
			"allowed_headers":              sc.CorsConfiguration.AllowedHeaders,
			"allowed_methods":              sc.CorsConfiguration.AllowedMethods,
			"allowed_origins":              sc.CorsConfiguration.AllowedOrigins,
			"enabled":                      sc.CorsConfiguration.Enabled,
			"exposed_headers":              sc.CorsConfiguration.ExposedHeaders,
			"preflight_max_age_in_seconds": sc.CorsConfiguration.PreflightMaxAgeInSeconds,
		},
	})
	if err != nil {
		return diag.Errorf("system_configuration.cors_configuration: %s", err.Error())
	}

	err = data.Set("event_log_configuration", []map[string]interface{}{
		{
			"number_to_retain": sc.EventLogConfiguration.NumberToRetain,
		},
	})
	if err != nil {
		return diag.Errorf("system_configuration.event_log_configuration: %s", err.Error())
	}

	err = data.Set("login_record_configuration", []map[string]interface{}{
		{
			"delete": []map[string]interface{}{
				{
					"enabled":                  sc.LoginRecordConfiguration.Delete.Enabled,
					"number_of_days_to_retain": sc.LoginRecordConfiguration.Delete.NumberOfDaysToRetain,
				},
			},
		},
	})
	if err != nil {
		return diag.Errorf("system_configuration.login_record_configuration: %s", err.Error())
	}

	if err := data.Set("report_timezone", sc.ReportTimezone); err != nil {
		return diag.Errorf("system_configuration.report_timezone: %s", err.Error())
	}

	err = data.Set("ui_configuration", []map[string]interface{}{
		{
			"header_color":    sc.UiConfiguration.HeaderColor,
			"logo_url":        sc.UiConfiguration.LogoURL,
			"menu_font_color": sc.UiConfiguration.MenuFontColor,
		},
	})
	if err != nil {
		return diag.Errorf("system_configuration.ui_configuration: %s", err.Error())
	}

	err = data.Set("webhook_event_log_configuration", []map[string]interface{}{
		{
			"delete": []map[string]interface{}{
				{
					"enabled":                  sc.WebhookEventLogConfiguration.Delete.Enabled,
					"number_of_days_to_retain": sc.WebhookEventLogConfiguration.Delete.NumberOfDaysToRetain,
				},
			},
		},
	})
	if err != nil {
		return diag.Errorf("system_configuration.webhook_event_log_configuration: %s", err.Error())
	}

	return nil
}

func getDefaultSystemConfigurationRequest() fusionauth.SystemConfigurationRequest {
	return fusionauth.SystemConfigurationRequest{
		SystemConfiguration: fusionauth.SystemConfiguration{
			AuditLogConfiguration: fusionauth.AuditLogConfiguration{
				Delete: fusionauth.DeleteConfiguration{
					Enableable: fusionauth.Enableable{
						Enabled: false,
					},
				},
			},
			CorsConfiguration: fusionauth.CORSConfiguration{
				AllowCredentials: false,
				AllowedHeaders: []string{
					"Accept",
					"Access-Control-Request-Headers",
					"Access-Control-Request-Method",
					"Authorization",
					"Content-Type",
					"Last-Modified",
					"Origin",
					"X-FusionAuth-TenantId",
					"X-Requested-With",
				},
				AllowedMethods: []fusionauth.HTTPMethod{
					fusionauth.HTTPMethod_GET,
					fusionauth.HTTPMethod_OPTIONS,
				},
				Enableable: fusionauth.Enableable{
					Enabled: true,
				},
				ExposedHeaders: []string{
					"Access-Control-Allow-Origin",
					"Access-Control-Allow-Credentials",
				},
				PreflightMaxAgeInSeconds: 1800,
			},
			EventLogConfiguration: fusionauth.EventLogConfiguration{
				NumberToRetain: 10000,
			},
			LoginRecordConfiguration: fusionauth.LoginRecordConfiguration{
				Delete: fusionauth.DeleteConfiguration{
					Enableable: fusionauth.Enableable{
						Enabled: false,
					},
				},
			},
			ReportTimezone: "America/Denver",
			WebhookEventLogConfiguration: fusionauth.WebhookEventLogConfiguration{
				Delete: fusionauth.DeleteConfiguration{
					Enableable: fusionauth.Enableable{
						Enabled: false,
					},
				},
			},
		},
	}
}
