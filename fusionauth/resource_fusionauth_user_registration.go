package fusionauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func newRegistration() *schema.Resource {
	return &schema.Resource{
		CreateContext: createRegistration,
		ReadContext:   readRegistration,
		UpdateContext: updateRegistration,
		DeleteContext: deleteRegistration,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Id of the User that is registering for the Application.",
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
			},
			"authentication_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The authentication token that may be used in place of the User’s password when authenticating against this application represented by this registration. This parameter is ignored if generateAuthenticationToken is set to true and instead the value will be generated.",
				Computed:    true,
				Sensitive:   true,
			},
			"generate_authentication_token": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if FusionAuth should generate an Authentication Token for this registration.",
			},
			"application_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Id of the Application that this registration is for.",
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "An object that can hold any information about the User for this registration that should be persisted.",
			},
			"preferred_languages": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An array of locale strings that give, in order, the User’s preferred languages for this registration. These are important for email templates and other localizable text.",
			},
			"roles": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The list of roles that the User has for this registration.",
			},
			"timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The User’s preferred timezone for this registration. The string will be in an IANA time zone format.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username of the User for this registration only.",
			},
			"skip_registration_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates to FusionAuth that it should skip registration verification even if it is enabled for the Application.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func buildRegistration(data *schema.ResourceData) fusionauth.RegistrationRequest {
	return fusionauth.RegistrationRequest{
		Registration: fusionauth.UserRegistration{
			ApplicationId:      data.Get("application_id").(string),
			Data:               data.Get("data").(map[string]interface{}),
			PreferredLanguages: handleStringSlice("preferred_languages", data),
			Roles:              handleStringSlice("roles", data),
			Timezone:           data.Get("timezone").(string),
			Username:           data.Get("username").(string),
		},
		GenerateAuthenticationToken:  data.Get("generate_authentication_token").(bool),
		SkipRegistrationVerification: data.Get("skip_registration_validation").(bool),
	}
}

func createRegistration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	reg := struct {
		Registration                 fusionauth.UserRegistration `json:"registration,omitempty"`
		SkipRegistrationVerification bool                        `json:"skipRegistrationVerification"`
	}{
		Registration: fusionauth.UserRegistration{
			ApplicationId:       data.Get("application_id").(string),
			AuthenticationToken: data.Get("authentication_token").(string),
			Data:                data.Get("data").(map[string]interface{}),
			PreferredLanguages:  handleStringSlice("preferred_languages", data),
			Roles:               handleStringSlice("roles", data),
			Timezone:            data.Get("timezone").(string),
			Username:            data.Get("username").(string),
		},
		SkipRegistrationVerification: data.Get("skip_registration_validation").(bool),
	}

	client := i.(Client)
	b, _ := json.Marshal(reg)
	b, err := sendCreateRegistration(b, data.Get("user_id").(string), data.Get("application_id").(string), client)
	if err != nil {
		return diag.Errorf("register err: %v", err)
	}

	_ = json.Unmarshal(b, &reg)

	data.SetId(reg.Registration.Id)
	return buildResourceDataFromRegistration(reg.Registration, data)
}

func sendCreateRegistration(b []byte, uid string, aid string, client Client) ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%s/%s/%s", strings.TrimRight(client.Host, "/"), "api/user/registration", uid, aid),
		bytes.NewBuffer(b),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", client.APIKey)
	req.Header.Add("Content-Type", "application/json")

	hc := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, _ = io.ReadAll(resp.Body)

	if err := checkResponse(resp.StatusCode, nil); err != nil {
		return nil, err
	}
	return b, nil
}

func readRegistration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	resp, faErrs, err := client.FAClient.RetrieveRegistration(data.Get("user_id").(string), data.Get("application_id").(string))
	if err != nil {
		return diag.Errorf("RetrieveRegistration err: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		data.SetId("")
		return nil
	}
	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return buildResourceDataFromRegistration(resp.Registration, data)
}

func buildResourceDataFromRegistration(r fusionauth.UserRegistration, data *schema.ResourceData) diag.Diagnostics {
	if err := data.Set("authentication_token", r.AuthenticationToken); err != nil {
		return diag.Errorf("registration.authentication_token: %s", err.Error())
	}
	if err := data.Set("application_id", r.ApplicationId); err != nil {
		return diag.Errorf("registration.application_id: %s", err.Error())
	}
	if err := data.Set("data", r.Data); err != nil {
		return diag.Errorf("registration.data: %s", err.Error())
	}
	if err := data.Set("preferred_languages", r.PreferredLanguages); err != nil {
		return diag.Errorf("registration.preferred_languages: %s", err.Error())
	}
	if err := data.Set("roles", r.Roles); err != nil {
		return diag.Errorf("registration.roles: %s", err.Error())
	}
	if err := data.Set("timezone", r.Timezone); err != nil {
		return diag.Errorf("registration.timezone: %s", err.Error())
	}
	if err := data.Set("username", r.Username); err != nil {
		return diag.Errorf("registration.username: %s", err.Error())
	}

	return nil
}

func updateRegistration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	ur := buildRegistration(data)

	resp, faErrs, err := client.FAClient.UpdateRegistration(data.Get("user_id").(string), ur)
	if err != nil {
		return diag.Errorf("UpdateRegistration err: %v", err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteRegistration(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)

	resp, faErrs, err := client.FAClient.DeleteRegistration(data.Get("user_id").(string), data.Get("application_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponse(resp.StatusCode, faErrs); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
