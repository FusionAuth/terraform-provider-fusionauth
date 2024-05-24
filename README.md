# FusionAuth Provider

This provider is used for setting up [FusionAuth](https://fusionauth.io).

For the rendered provider usage documentation, visit the [Terraform Registry](https://registry.terraform.io/providers/gpsinsight/fusionauth/latest/docs).

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
* application
* application/{application_id}/role
* email
* entity
* entity grant
* entity type
* entity type permission
* form
* form field
* group
* generic connector
* key
* imported key
* lambda
* identity provider
    - OpenID Connect
    - Google
    - Apple
    - External JWT
    - Facebook
    - SAML v2
    - Sony PSN
    - Steam
    - Twitch
    - Xbox
* themes
* user
* user action
* webhook
* tenants

## Testing

Please add tests to the relevant files.

To run tests:

```
cd fusionauth
go test
```

## Known issues

If you do not specify permissions when adding an API key, you will get a key created that has no permissions. See the following issues for more details.
- https://github.com/gpsinsight/terraform-provider-fusionauth/issues/126
- https://github.com/FusionAuth/fusionauth-issues/issues/1675
