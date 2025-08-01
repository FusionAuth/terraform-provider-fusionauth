# FusionAuth Provider

This provider is used for setting up [FusionAuth](https://fusionauth.io).

Learn more about [using FusionAuth and Terraform together](https://fusionauth.io/docs/operate/deploy/terraform).

## Example Usage - Provider Configuration via `tfvars variables`

```hcl
terraform {
  required_providers {
    fusionauth = {
      source  = "fusionauth/fusionauth"
      version = "~> 1.0.0"
    }
  }
}

// Provider configuration set via tfvars variables
provider "fusionauth" {
  api_key = var.fusionauth_api_key
  host    = var.fusionauth_host
}
```

## Example Usage - Provider Configuration via `environment variables`

```hcl
terraform {
  required_providers {
    fusionauth = {
      source  = "fusionauth/fusionauth"
      version = "~> 1.0.0"
    }
  }
}

// Provider configuration set via environment variables
provider "fusionauth" {}
```

## Argument Reference

* `api_key` - (Required) The API Key for the FusionAuth instance. Alternatively, can be configured using the `FA_API_KEY` environment variable.
* `host` - (Required) Host for FusionAuth instance. Alternatively, can be configured using the `FA_DOMAIN` environment variable.
