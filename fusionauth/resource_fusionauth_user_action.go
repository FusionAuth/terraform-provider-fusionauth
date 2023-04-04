package fusionauth

import (
	"context"
	"net/http"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceUserAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: createUserAction,
		ReadContext:   readUserAction,
		UpdateContext: updateUserAction,
		DeleteContext: deleteUserAction,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of this User Action.",
			},
			"user_action_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The id of this User Action.",
				ValidateFunc: validation.IsUUID,
			},
			"cancel_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template that is used when User Actions are canceled.",
				ValidateFunc: validation.IsUUID,
			},
			"end_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template that is used when User Actions expired automatically (end).",
				ValidateFunc: validation.IsUUID,
			},
			"include_email_in_event_json": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to include the email information in the JSON that is sent to the Webhook when a user action is taken.",
			},
			"localized_names": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A mapping of localized names for this User Action. The key is the Locale and the value is the name of the User Action for that language.",
			},
			"modify_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template that is used when User Actions are modified.",
				ValidateFunc: validation.IsUUID,
			},
			"options": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The list of User Action Options.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of this User Action Option.",
						},
						"localized_names": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "A mapping of localized names for this User Action Option. The key is the Locale and the value is the name of the User Action Option for that language.",
						},
					},
				},
			},
			"prevent_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not this User Action will prevent user login. When this value is set to true the action must also be marked as a time based action. See `temporal`.",
			},
			"send_end_event": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not FusionAuth will send events to any registered Webhooks when this User Action expires.",
			},
			"start_email_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The Id of the Email Template that is used when User Actions are started (created).",
				ValidateFunc: validation.IsUUID,
			},
			"temporal": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not this User Action is time-based (temporal).",
			},
			"user_emailing_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not email is enabled for this User Action.",
			},
			"user_notifications_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not user notifications are enabled for this User Action.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildUserAction(data *schema.ResourceData) fusionauth.UserAction {
	ua := fusionauth.UserAction{
		Name: data.Get("name").(string),
	}

	if d, ok := data.GetOk("user_action_id"); ok {
		ua.Id = d.(string)
	}
	if d, ok := data.GetOk("cancel_email_template_id"); ok {
		ua.CancelEmailTemplateId = d.(string)
	}
	if d, ok := data.GetOk("end_email_template_id"); ok {
		ua.EndEmailTemplateId = d.(string)
	}
	if d, ok := data.GetOk("include_email_in_event_json"); ok {
		ua.IncludeEmailInEventJSON = d.(bool)
	}
	if i, ok := data.GetOk("localized_names"); ok {
		ua.LocalizedNames = intMapToStringMap(i.(map[string]interface{}))
	}
	if d, ok := data.GetOk("modify_email_template_id"); ok {
		ua.ModifyEmailTemplateId = d.(string)
	}
	if d, ok := data.GetOk("options"); ok {
		ua.Options = buildUserActionOptions(d)
	}
	if d, ok := data.GetOk("prevent_login"); ok {
		ua.PreventLogin = d.(bool)
	}
	if d, ok := data.GetOk("send_end_event"); ok {
		ua.SendEndEvent = d.(bool)
	}
	if d, ok := data.GetOk("start_email_template_id"); ok {
		ua.StartEmailTemplateId = d.(string)
	}
	if d, ok := data.GetOk("temporal"); ok {
		ua.Temporal = d.(bool)
	}
	if d, ok := data.GetOk("user_emailing_enabled"); ok {
		ua.UserEmailingEnabled = d.(bool)
	}
	if d, ok := data.GetOk("user_notifications_enabled"); ok {
		ua.UserNotificationsEnabled = d.(bool)
	}

	return ua
}

func createUserAction(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	userAction := buildUserAction(data)

	resp, faErrs, err := client.FAClient.CreateUserAction(userAction.Id, fusionauth.UserActionRequest{
		UserAction: userAction,
	})

	if err != nil {
		return diag.Errorf("CreateUser err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(resp.UserAction.Id)

	return readUserAction(ctx, data, i)
}

func readUserAction(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, err := client.FAClient.RetrieveUserAction(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}

	if err := data.Set("name", resp.UserAction.Name); err != nil {
		return diag.Errorf("user_action.name: %s", err.Error())
	}
	if err := data.Set("user_action_id", resp.UserAction.Id); err != nil {
		return diag.Errorf("user_action.id: %s", err.Error())
	}
	if err := data.Set("cancel_email_template_id", resp.UserAction.CancelEmailTemplateId); err != nil {
		return diag.Errorf("user_action.cancel_email_template_id: %s", err.Error())
	}
	if err := data.Set("end_email_template_id", resp.UserAction.EndEmailTemplateId); err != nil {
		return diag.Errorf("user_action.end_email_template_id: %s", err.Error())
	}
	if err := data.Set("include_email_in_event_json", resp.UserAction.IncludeEmailInEventJSON); err != nil {
		return diag.Errorf("user_action.include_email_in_event_json: %s", err.Error())
	}
	if err := data.Set("localized_names", resp.UserAction.LocalizedNames); err != nil {
		return diag.Errorf("user_action.localized_names: %s", err.Error())
	}
	if err := data.Set("modify_email_template_id", resp.UserAction.ModifyEmailTemplateId); err != nil {
		return diag.Errorf("user_action.modify_email_template_id: %s", err.Error())
	}

	options := make([]map[string]interface{}, 0, len(resp.UserAction.Options))
	for _, opt := range resp.UserAction.Options {
		options = append(options, map[string]interface{}{
			"name":            opt.Name,
			"localized_names": opt.LocalizedNames,
		})
	}
	if err := data.Set("options", options); err != nil {
		return diag.Errorf("user_action.options: %s", err.Error())
	}
	if err := data.Set("prevent_login", resp.UserAction.PreventLogin); err != nil {
		return diag.Errorf("user_action.prevent_login: %s", err.Error())
	}
	if err := data.Set("send_end_event", resp.UserAction.SendEndEvent); err != nil {
		return diag.Errorf("user_action.send_end_event: %s", err.Error())
	}
	if err := data.Set("start_email_template_id", resp.UserAction.StartEmailTemplateId); err != nil {
		return diag.Errorf("user_action.start_email_template_id: %s", err.Error())
	}
	if err := data.Set("temporal", resp.UserAction.Temporal); err != nil {
		return diag.Errorf("user_action.temporal: %s", err.Error())
	}
	if err := data.Set("user_emailing_enabled", resp.UserAction.UserEmailingEnabled); err != nil {
		return diag.Errorf("user_action.user_emailing_enabled: %s", err.Error())
	}
	if err := data.Set("user_notifications_enabled", resp.UserAction.UserNotificationsEnabled); err != nil {
		return diag.Errorf("user_action.user_notifications_enabled: %s", err.Error())
	}

	return nil
}

func updateUserAction(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	resp, faErrs, err := client.FAClient.UpdateUserAction(data.Id(), fusionauth.UserActionRequest{
		UserAction: buildUserAction(data),
	})
	if err != nil {
		return diag.Errorf("UpdateUserAction err: %v", err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return readUserAction(ctx, data, i)
}

func deleteUserAction(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	resp, faErrs, err := client.FAClient.DeleteUserAction(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildUserActionOptions(d interface{}) []fusionauth.UserActionOption {
	opts := make([]fusionauth.UserActionOption, 0, len(d.([]*schema.ResourceData)))
	for _, v := range d.([]*schema.ResourceData) {
		uao := fusionauth.UserActionOption{
			Name: v.Get("name").(string),
		}

		if i, ok := v.GetOk("localized_names"); ok {
			uao.LocalizedNames = intMapToStringMap(i.(map[string]interface{}))
		}

		opts = append(opts, uao)
	}

	return opts
}
