package fusionauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func newUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: createUser,
		ReadContext:   readUser,
		UpdateContext: updateUser,
		DeleteContext: deleteUser,
		Schema:        userSchemaV1().Schema,
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{{
			Type:    userSchemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: upgradeUserSchemaV0ToV1,
			Version: 0,
		}},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func createUser(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	req, diags := dataToUserRequest(data)
	if diags != nil {
		return diags
	}

	oldTenantID := client.FAClient.TenantId
	client.FAClient.TenantId = req.User.TenantId
	defer func() {
		client.FAClient.TenantId = oldTenantID
	}()

	resp, faErrs, err := client.FAClient.CreateUser(req.User.Id, req)
	if err != nil {
		return diag.Errorf("CreateUser err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return userResponseToData(data, resp)
}

func readUser(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.RetrieveUser(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return userResponseToData(data, resp)
}

func updateUser(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	req, diags := dataToUserRequest(data)
	if diags != nil {
		return diags
	}

	resp, faErrs, err := client.FAClient.UpdateUser(data.Id(), req)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return userResponseToData(data, resp)
}

func deleteUser(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	id := data.Id()

	resp, faErrs, err := client.FAClient.DeleteUser(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		// User successfully deleted
		data.SetId("")
		return nil
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func dataToUserRequest(data *schema.ResourceData) (req fusionauth.UserRequest, diags diag.Diagnostics) {
	var userID string
	if datum, ok := data.GetOk("user_id"); ok {
		userID = datum.(string)
	}

	twoFactorMethods, subDiags := dataToTwoFactorMethods(data)
	if subDiags != nil {
		diags = append(diags, subDiags...)
	}

	resourceData, subDiags := jsonStringToMapStringInterface(data.Get("data").(string))
	if subDiags != nil {
		diags = append(diags, subDiags...)
	}

	req = fusionauth.UserRequest{
		ApplicationId:      data.Get("application_id").(string),
		DisableDomainBlock: data.Get("disable_domain_block").(bool),
		User: fusionauth.User{
			TenantId:           data.Get("tenant_id").(string),
			BirthDate:          data.Get("birth_date").(string),
			Data:               resourceData,
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
				Id:                     userID,
				EncryptionScheme:       data.Get("encryption_scheme").(string),
				Password:               data.Get("password").(string),
				PasswordChangeRequired: data.Get("password_change_required").(bool),
				Username:               data.Get("username").(string),
				UsernameStatus:         fusionauth.ContentStatus(data.Get("username_status").(string)),
			},
			TwoFactor: fusionauth.UserTwoFactorConfiguration{
				Methods:       twoFactorMethods,
				RecoveryCodes: handleStringSlice("two_factor_recovery_codes", data),
			},
		},
		SendSetPasswordEmail: data.Get("send_set_password_email").(bool),
		SkipVerification:     data.Get("skip_verification").(bool),
	}

	return req, diags
}

func userResponseToData(data *schema.ResourceData, resp *fusionauth.UserResponse) diag.Diagnostics {
	data.SetId(resp.User.Id)

	if err := data.Set("user_id", resp.User.Id); err != nil {
		return diag.Errorf("user.user_id: %s", err.Error())
	}
	if err := data.Set("tenant_id", resp.User.TenantId); err != nil {
		return diag.Errorf("user.tenant_id: %s", err.Error())
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
	// Do not set parent_email in TF state as the server never returns the data.
	if err := data.Set("preferred_languages", resp.User.PreferredLanguages); err != nil {
		return diag.Errorf("user.preferred_languages: %s", err.Error())
	}
	if err := data.Set("timezone", resp.User.Timezone); err != nil {
		return diag.Errorf("user.timezone: %s", err.Error())
	}
	if err := data.Set("username", resp.User.Username); err != nil {
		return diag.Errorf("user.username: %s", err.Error())
	}
	if err := data.Set("username_status", resp.User.UsernameStatus); err != nil {
		return diag.Errorf("user.username_status: %s", err.Error())
	}
	if err := data.Set("password_change_required", resp.User.PasswordChangeRequired); err != nil {
		return diag.Errorf("user.password_change_required: %s", err.Error())
	}

	if err := data.Set("two_factor_recovery_codes", resp.User.TwoFactor.RecoveryCodes); err != nil {
		return diag.Errorf("user.two_factor_recovery_codes: %s", err.Error())
	}

	currentTwoFactorMethods := dataToTwoFactorMethodMap(data)
	twoFactorMethodsData := make([]map[string]interface{}, len(resp.User.TwoFactor.Methods))
	for i, twoFactorMethod := range resp.User.TwoFactor.Methods {
		var secret string
		if strings.ToLower(twoFactorMethod.Method) == "authenticator" {
			if currentMethod, ok := currentTwoFactorMethods[twoFactorMethod.Method]; ok {
				// FusionAuth doesn't return the secret via the API, so we need
				// to do some manual i/o to ensure terraform state lines up
				// and changes can be detected properly.
				secret = currentMethod.Secret
			}
		}

		twoFactorMethodsData[i] = map[string]interface{}{
			"two_factor_method_id":      twoFactorMethod.Id,
			"authenticator_algorithm":   twoFactorMethod.Authenticator.Algorithm,
			"authenticator_code_length": twoFactorMethod.Authenticator.CodeLength,
			"authenticator_time_step":   twoFactorMethod.Authenticator.TimeStep,
			"email":                     twoFactorMethod.Email,
			"method":                    twoFactorMethod.Method,
			"mobile_phone":              twoFactorMethod.MobilePhone,
			"secret":                    secret,
		}
	}
	if err := data.Set("two_factor_methods", twoFactorMethodsData); err != nil {
		return diag.Errorf("user.two_factor_methods: %s", err.Error())
	}

	return nil
}

func dataToTwoFactorMethods(data *schema.ResourceData) (twoFactorMethods []fusionauth.TwoFactorMethod, diags diag.Diagnostics) {
	twoFactorMethodsData, ok := data.Get("two_factor_methods").([]interface{})
	if twoFactorMethodsData == nil || !ok {
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Unable to convert two factor methods",
				Detail:   "two_factor_methods unable to be typecast to []interface{}",
			})
		}

		// Nothing to do here!
		return twoFactorMethods, diags
	}

	twoFactorMethods = make([]fusionauth.TwoFactorMethod, len(twoFactorMethodsData))
	for i, twoFactorMethodsDatum := range twoFactorMethodsData {
		if twoFactorMethod, ok := twoFactorMethodsDatum.(map[string]interface{}); !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Unable to convert a two factor method",
				Detail:   fmt.Sprintf("two_factor_methods.%d: %#+v unable to be typecast to []interface{}", i, twoFactorMethodsDatum),
			})
		} else {
			twoFactorMethods[i] = fusionauth.TwoFactorMethod{
				Id: twoFactorMethod["two_factor_method_id"].(string),
				Authenticator: fusionauth.AuthenticatorConfiguration{
					Algorithm:  fusionauth.TOTPAlgorithm(twoFactorMethod["authenticator_algorithm"].(string)),
					CodeLength: twoFactorMethod["authenticator_code_length"].(int),
					TimeStep:   twoFactorMethod["authenticator_time_step"].(int),
				},
				Email:       twoFactorMethod["email"].(string),
				Method:      twoFactorMethod["method"].(string),
				MobilePhone: twoFactorMethod["mobile_phone"].(string),
				Secret:      twoFactorMethod["secret"].(string),
			}
		}
	}

	return twoFactorMethods, diags
}

func dataToTwoFactorMethodMap(data *schema.ResourceData) (methodMap map[string]fusionauth.TwoFactorMethod) {
	methodMap = map[string]fusionauth.TwoFactorMethod{}

	methodList, errs := dataToTwoFactorMethods(data)
	if errs != nil {
		return
	}

	for _, method := range methodList {
		methodMap[method.Method] = method
	}

	return
}
