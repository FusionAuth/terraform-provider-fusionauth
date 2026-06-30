# FusionAuth Provider

This provider is used for setting up [FusionAuth](https://fusionauth.io).

For the rendered provider usage documentation, visit the [Terraform Registry](https://registry.terraform.io/providers/FusionAuth/fusionauth/latest/docs).

## Please Read

November 16th, 2023
This Terraform Provider has moved to the [FusionAuth](https://github.com/FusionAuth) organization.

FusionAuth would like to thank [GPS Insight](https://github.com/gpsinsight) for all of their efforts to build and maintain this provider for the past three years!

The purpose of this change in ownership is to allow FusionAuth to be in a better position to manage pull requests, and work towards full parity with the FusionAuth API.

Please continue to use and provide feedback on this provider as you have in the past, we are happy to accept pull requests.

## Argument Reference

* `api_key` - (Required) The API Key for the FusionAuth instance
* `host` - (Required) Host for FusionAuth instance

## Resources Available

* API Key
* Application
* Application OAuth Scope
* Application Role
* Consent
* Email
* Entity
* Entity Grant
* Entity Type
* Entity Type Permission
* Form
* Form Field
* Generic Connector
* Generic Messenger
* Group
* Identity Provider
  * Apple
  * External JWT
  * Facebook
  * Google
  * LinkedIn
  * OpenID Connect
  * SAML v2
  * SAML v2 IdP Initiated
  * Sony PSN
  * Steam
  * Twitch
  * Xbox
* Imported Key
* Key
* Lambda
* LDAP Connector
* Reactor
* Registration
* SMS Message Template
* System Configuration
* Tenant
* Tenant Manager Configuration
* Theme
* Twilio Messenger
* User
* User Action
* User Group Membership
* Webhook

## Data Sources Available

* Application
* Application OAuth Scope
* Application Role
* Consent
* Email
* Form
* Form Field
* Generic Connector
* Generic Messenger
* Identity Provider
* Lambda
* LDAP Connector
* SMS Message Template
* Tenant
* Theme
* Twilio Messenger
* User
* User Group Membership

## Testing

Please add tests to the relevant files.

To run tests:

```
cd fusionauth
go test
```
