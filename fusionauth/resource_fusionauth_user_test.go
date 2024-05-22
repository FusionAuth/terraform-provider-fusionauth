package fusionauth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gpsinsight/terraform-provider-fusionauth/fusionauth/testdata"
)

func TestAccFusionauthUser_basic(t *testing.T) {
	resourceName := randString10()
	tfResourcePath := fmt.Sprintf("fusionauth_user.test_%s", resourceName)

	startSendSetPasswordEmail, endSendSetPasswordEmail := false, true
	startSkipVerification, endSkipVerification := true, false
	startBirthDate, endBirthDate := "2000-09-01", "2020-09-01"
	startDataHCL := `{
  testBool = true
  testString = "Hello world!"
  testInt    = 1
  testFloat  = 1.23
  testArray  = [ 
    "test:urn:permission1",
  ]
  test = {
    nested = [
      "item",
      "s",
    ]
  }
}`
	endDataHCL := `{
  testBool = false
  testStr = "Bonjour le monde!"
  testInt    = 2
  testFloat  = 3.21
  testArray  = [ 
    "test:urn:permission2",
  ]
  test = {
    nested = "item"
  }
}`
	startDataJSON := `
{
  "testBool":   true,
  "testString": "Hello world!",
  "testInt":    1,
  "testFloat":  1.23,
  "testArray":  [
    "test:urn:permission1"
  ],
  "test": {
    "nested": [
      "item",
      "s"
    ]
  }
}
`
	endDataJSON := `{
  "testBool":  false,
  "testStr":   "Bonjour le monde!",
  "testInt":   2,
  "testFloat": 3.21,
  "testArray": [
    "test:urn:permission2"
  ],
  "test": {
    "nested": "item"
  }
}
`

	startEmail, endEmail := "john.s@example.com", "jon.snow@example.com"
	startEncryptionScheme, endEncryptionScheme := "salted-md5", "bcrypt"
	startExpiry, endExpiry := "7955114522000", "43017447783000"
	startFirstName, endFirstName := "John", "Jon"
	startFullName, endFullName := "test-acc John Smith", "test-acc Jon Snow"
	startImageURL, endImageURL := "https://gravatar.com/avatar/0dc5552bbda9ab3d62a9d7612f471dd9?s=400&d=mp&r=g", "https://gravatar.com/avatar/32efbfa435860a48dd5af9626ea326d3?s=400&d=mp&r=g"
	startLastName, endLastName := "Smith", "Snow"
	startMiddleName, endMiddleName := "A", "'King of the North'"
	startMobilePhone, endMobilePhone := "+642598765432", "+642512345678"
	startParentEmail, endParentEmail := "old.smith@example.com", "eddard.stark@example.com"
	startPassword, endPassword := "P@ssw0rd", "Sup3r|Secr3t"
	startPasswordChangeRequired, endPasswordChangeRequired := false, true
	// preferredLanguages must be a 2 length slice.
	startPreferredLanguages, endPreferredLanguages := []string{"en", "fr"}, []string{"en", "es"}
	startTimezone, endTimezone := "Europe/Paris", "America/Mexico_City"
	startTwoFactorMethodsEmail, endTwoFactorMethodsEmail := "john.smith@example.com", "john.smith@example.com"                 // two factor email address can't be altered once created.
	startTwoFactorMethodsMobilePhone, endTwoFactorMethodsMobilePhone := "+64987654321", "+64987654321"                         // two factor mobile phone number can't be altered once created.
	startTwoFactorMethodsSecret, endTwoFactorMethodsSecret := "Sup3r\nSecr3t\nMFA\nSecr3t", "X7ra\nSup3r\nSecr3t\nMFA\nSecr3t" //nolint:gosec
	startUsername, endUsername := "john.smith", "jon.snow"
	startUsernameStatus, endUsernameStatus := "ACTIVE", "PENDING"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFusionauthUserDestroy,
		Steps: []resource.TestStep{
			{
				// Test resource create
				Config: testAccUserResource(
					resourceName,
					startSendSetPasswordEmail,
					startSkipVerification,
					startBirthDate,
					startDataHCL,
					startEmail,
					startEncryptionScheme,
					startExpiry,
					startFirstName,
					startFullName,
					startImageURL,
					startLastName,
					startMiddleName,
					startMobilePhone,
					startParentEmail,
					startPassword,
					startPasswordChangeRequired,
					startPreferredLanguages,
					startTimezone,
					startTwoFactorMethodsEmail,
					startTwoFactorMethodsMobilePhone,
					startTwoFactorMethodsSecret,
					startUsername,
					startUsernameStatus,
				),
				Check: testUserBasicAccCheckFuncs(
					tfResourcePath,
					startSendSetPasswordEmail,
					startSkipVerification,
					startBirthDate,
					startDataJSON,
					startEmail,
					startEncryptionScheme,
					startExpiry,
					startFirstName,
					startFullName,
					startImageURL,
					startLastName,
					startMiddleName,
					startMobilePhone,
					startParentEmail,
					startPassword,
					startPasswordChangeRequired,
					startPreferredLanguages,
					startTimezone,
					startTwoFactorMethodsEmail,
					startTwoFactorMethodsMobilePhone,
					startTwoFactorMethodsSecret,
					startUsername,
					startUsernameStatus,
				),
			},
			{
				// Test update/mutate resource state
				Config: testAccUserResource(
					resourceName,
					endSendSetPasswordEmail,
					endSkipVerification,
					endBirthDate,
					endDataHCL,
					endEmail,
					endEncryptionScheme,
					endExpiry,
					endFirstName,
					endFullName,
					endImageURL,
					endLastName,
					endMiddleName,
					endMobilePhone,
					endParentEmail,
					endPassword,
					endPasswordChangeRequired,
					endPreferredLanguages,
					endTimezone,
					endTwoFactorMethodsEmail,
					endTwoFactorMethodsMobilePhone,
					endTwoFactorMethodsSecret,
					endUsername,
					endUsernameStatus,
				),
				Check: testUserBasicAccCheckFuncs(
					tfResourcePath,
					endSendSetPasswordEmail,
					endSkipVerification,
					endBirthDate,
					endDataJSON,
					endEmail,
					endEncryptionScheme,
					endExpiry,
					endFirstName,
					endFullName,
					endImageURL,
					endLastName,
					endMiddleName,
					endMobilePhone,
					endParentEmail,
					endPassword,
					endPasswordChangeRequired,
					endPreferredLanguages,
					endTimezone,
					endTwoFactorMethodsEmail,
					endTwoFactorMethodsMobilePhone,
					endTwoFactorMethodsSecret,
					endUsername,
					endUsernameStatus,
				),
			},
			{
				ResourceName:      tfResourcePath,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// The following fields are not returned via RetrieveUser
					// and as such, can't be imported/pushed into the Terraform
					// state.
					"disable_domain_block",
					"encryption_scheme",
					"password",
					"parent_email",
					"send_set_password_email",
					"skip_verification",
					"two_factor_methods.0.secret",
				},
			},
		},
	})
}

// testTenantAccTestCheckFuncs abstracts the test case setup required between
// create and update testing.
func testUserBasicAccCheckFuncs(
	tfResourcePath string,
	sendSetPasswordEmail bool,
	skipVerification bool,
	birthDate string,
	data string,
	email string,
	encryptionScheme string,
	expiry string,
	firstName string,
	fullName string,
	imageURL string,
	lastName string,
	middleName string,
	mobilePhone string,
	parentEmail string,
	password string,
	passwordChangeRequired bool,
	// preferredLanguages must be a 2 length slice.
	preferredLanguages []string,
	timezone string,
	twoFactorMethodsEmail string,
	twoFactorMethodsMobilePhone string,
	twoFactorMethodsSecret string,
	username string,
	usernameStatus string,
) resource.TestCheckFunc {
	testChecks := []resource.TestCheckFunc{
		testAccCheckFusionauthUserExists(tfResourcePath),
		resource.TestCheckResourceAttrSet(tfResourcePath, "user_id"),
		resource.TestCheckResourceAttrSet(tfResourcePath, "tenant_id"),
		resource.TestCheckResourceAttr(tfResourcePath, "send_set_password_email", fmt.Sprintf("%t", sendSetPasswordEmail)),
		resource.TestCheckResourceAttr(tfResourcePath, "skip_verification", fmt.Sprintf("%t", skipVerification)),
		// User Object Properties
		resource.TestCheckResourceAttr(tfResourcePath, "birth_date", birthDate),
		testCheckResourceAttrJSON(tfResourcePath, "data", data),
		resource.TestCheckResourceAttr(tfResourcePath, "email", email),
		resource.TestCheckResourceAttr(tfResourcePath, "encryption_scheme", encryptionScheme),
		resource.TestCheckResourceAttr(tfResourcePath, "expiry", expiry),
		resource.TestCheckResourceAttr(tfResourcePath, "first_name", firstName),
		resource.TestCheckResourceAttr(tfResourcePath, "full_name", fullName),
		resource.TestCheckResourceAttr(tfResourcePath, "image_url", imageURL),
		resource.TestCheckResourceAttr(tfResourcePath, "last_name", lastName),
		resource.TestCheckResourceAttr(tfResourcePath, "middle_name", middleName),
		resource.TestCheckResourceAttr(tfResourcePath, "mobile_phone", mobilePhone),
		resource.TestCheckResourceAttr(tfResourcePath, "parent_email", parentEmail),
		resource.TestCheckResourceAttr(tfResourcePath, "password", password),
		resource.TestCheckResourceAttr(tfResourcePath, "password_change_required", fmt.Sprintf("%t", passwordChangeRequired)),
		resource.TestCheckResourceAttr(tfResourcePath, "preferred_languages.0", preferredLanguages[0]),
		resource.TestCheckResourceAttr(tfResourcePath, "preferred_languages.1", preferredLanguages[1]),
		resource.TestCheckResourceAttr(tfResourcePath, "timezone", timezone),
		resource.TestCheckResourceAttr(tfResourcePath, "username", username),
		resource.TestCheckResourceAttr(tfResourcePath, "username_status", usernameStatus),
	}

	mfaTestCount := 0
	if twoFactorMethodsSecret != "" {
		keyPrefix := fmt.Sprintf("two_factor_methods.%d", mfaTestCount)
		testChecks = append(
			testChecks,
			resource.TestCheckResourceAttr(tfResourcePath, fmt.Sprintf("%s.method", keyPrefix), "authenticator"),
			resource.TestCheckResourceAttr(tfResourcePath, fmt.Sprintf("%s.authenticator_algorithm", keyPrefix), "HmacSHA1"),
			resource.TestCheckResourceAttr(tfResourcePath, fmt.Sprintf("%s.authenticator_code_length", keyPrefix), "6"),
			resource.TestCheckResourceAttr(tfResourcePath, fmt.Sprintf("%s.authenticator_time_step", keyPrefix), "30"),
			resource.TestCheckResourceAttr(tfResourcePath, fmt.Sprintf("%s.secret", keyPrefix), base64.StdEncoding.EncodeToString([]byte(twoFactorMethodsSecret))),
		)
		mfaTestCount++
	}

	if twoFactorMethodsEmail != "" {
		keyPrefix := fmt.Sprintf("two_factor_methods.%d", mfaTestCount)
		testChecks = append(
			testChecks,
			resource.TestCheckResourceAttr(tfResourcePath, fmt.Sprintf("%s.method", keyPrefix), "email"),
			resource.TestCheckResourceAttr(tfResourcePath, fmt.Sprintf("%s.email", keyPrefix), twoFactorMethodsEmail),
		)
		mfaTestCount++
	}

	if twoFactorMethodsMobilePhone != "" {
		keyPrefix := fmt.Sprintf("two_factor_methods.%d", mfaTestCount)
		testChecks = append(
			testChecks,
			resource.TestCheckResourceAttr(tfResourcePath, fmt.Sprintf("%s.method", keyPrefix), "sms"),
			resource.TestCheckResourceAttr(tfResourcePath, fmt.Sprintf("%s.mobile_phone", keyPrefix), twoFactorMethodsMobilePhone),
		)
	}

	return resource.ComposeTestCheckFunc(testChecks...)
}

func testAccCheckFusionauthUserExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource id is set")
		}

		faClient := fusionauthClient()
		user, faErrs, err := faClient.RetrieveUser(rs.Primary.ID)
		if err != nil {
			// low-level error performing api request
			return fmt.Errorf("retrieveuser error: %#+v", err)
		}

		if faErrs != nil && faErrs.Present() {
			// Fusionauth has errors to report
			return fmt.Errorf("retrieveuser fusionauth errors: %#+v", faErrs)
		}

		if user == nil || user.StatusCode != http.StatusOK {
			// This is a weird edge case...
			return fmt.Errorf("retrieveuser failed to get user: %#+v", user)
		}

		return nil
	}
}

func testAccCheckFusionauthUserDestroy(s *terraform.State) error {
	faClient := fusionauthClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fusionauth_user" {
			continue
		}

		// Retry for eventual consistency
		err := retry.RetryContext(context.Background(), 10*time.Second, func() *retry.RetryError {
			user, faErrs, err := faClient.RetrieveUser(rs.Primary.ID)
			if err != nil {
				// low-level error performing api request
				return retry.NonRetryableError(err)
			}

			if faErrs != nil && faErrs.Present() {
				// Fusionauth has errors to report
				return retry.NonRetryableError(faErrs)
			}

			if user != nil && user.StatusCode == http.StatusNotFound {
				// resource destroyed!
				return nil
			}

			return retry.RetryableError(fmt.Errorf("fusionauth user still exists: %s", rs.Primary.ID))
		})

		if err != nil {
			// We failed destroying the user...
			return err
		}
	}

	return nil
}

// testAccUserResource returns the required terraform configuration to test terraform user operations.
func testAccUserResource(
	resourceName string,
	sendSetPasswordEmail bool,
	skipVerification bool,
	birthDate string,
	data string,
	email string,
	encryptionScheme string,
	expiry string,
	firstName string,
	fullName string,
	imageURL string,
	lastName string,
	middleName string,
	mobilePhone string,
	parentEmail string,
	password string,
	passwordChangeRequired bool,
	preferredLanguages []string,
	timezone string,
	twoFactorMethodsEmail string,
	twoFactorMethodsMobilePhone string,
	twoFactorMethodsSecret string,
	username string,
	usernameStatus string,
) string {
	if data == "" {
		data = "{}"
	}

	var mfaAuthenticator string
	if twoFactorMethodsSecret != "" {
		mfaAuthenticator = fmt.Sprintf(`
  two_factor_methods {
    method                    = "authenticator"
    authenticator_algorithm   = "HmacSHA1"  # With the current implementation, this will always be HmacSHA1. 
    authenticator_code_length = 6           # With the current implementation, this will always be 6.
    authenticator_time_step   = 30          # With the current implementation, this will always be 30.
    secret                    = "%s"
  }
`, base64.StdEncoding.EncodeToString([]byte(twoFactorMethodsSecret)))
	}

	var mfaEmail string
	if twoFactorMethodsEmail != "" {
		mfaEmail = fmt.Sprintf(`
  two_factor_methods {
    method                    = "email"
    email                     = "%s"
  }
`, twoFactorMethodsEmail)
	}

	var mfaSms string
	if twoFactorMethodsMobilePhone != "" {
		mfaSms = fmt.Sprintf(`
  two_factor_methods {
    method                    = "sms"
    mobile_phone              = "%s"
  }
`, twoFactorMethodsMobilePhone)
	}

	return testAccUserResourceConfigBase(resourceName) +
		fmt.Sprintf(`
resource "fusionauth_user" "test_%[1]s" {
  tenant_id               = fusionauth_tenant.test_%[1]s.id
  send_set_password_email = %[2]t
  skip_verification       = %[3]t

  birth_date               = "%[4]s"
  data                     = jsonencode(%[5]s)
  email                    = "%[6]s"
  encryption_scheme        = "%[7]s"
  expiry                   = %[8]s
  first_name               = "%[9]s" 
  full_name                = "%[10]s" 
  image_url                = "%[11]s" 
  last_name                = "%[12]s" 
  middle_name              = "%[13]s" 
  mobile_phone             = "%[14]s" 
  parent_email             = "%[15]s" 
  password                 = "%[16]s" 
  password_change_required = %[17]t 
  preferred_languages      = %[18]s 
  timezone                 = "%[19]s"

  # Two Factor Methods%[20]s%[21]s%[22]s

  username        = "%[23]s"
  username_status = "%[24]s"
}
`,
			resourceName,
			sendSetPasswordEmail,
			skipVerification,
			birthDate,
			data,
			email,
			encryptionScheme,
			expiry,
			firstName,
			fullName,
			imageURL,
			lastName,
			middleName,
			mobilePhone,
			parentEmail,
			password,
			passwordChangeRequired,
			stringSliceToHCL(preferredLanguages),
			timezone,
			mfaAuthenticator,
			mfaEmail,
			mfaSms,
			username,
			usernameStatus,
		)
}

// testAccUserResourceConfigBase generates the prerequisite resources required
// for creating a user.
func testAccUserResourceConfigBase(resourceName string) string {
	return testAccAccessTokenKeyResourceConfig(resourceName) +
		testAccIDTokenKeyResourceConfig(resourceName) +
		testAccThemeResourceConfig(
			resourceName,
			testdata.MessageProperties(""),
			"/* stylez */",
			generateFusionAuthTemplate(),
		) +
		testAccTenantResourceConfig(
			resourceName,
			resourceName,
			"no-reply@example.com",
			30,
			false,
			false,
		)
}
