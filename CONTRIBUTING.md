Contributing Guide

## Guidline for Updating

1. Scroll through [FusionAuth's API Docs](https://fusionauth.io/docs/v1/tech/apis/)
2. Update each resouce and make sure to update the docs!!!
3. Make sure tests work (and maybe add new ones!)
3. Submit PR to this repo
4. Upon merge, maintainer will create new git tag kicking off the build process.
5. [Terraform Registry](https://registry.terraform.io/providers/FusionAuth/fusionauth/latest) will pick up the changes

```
git tag v0.1.71
git push origin --tag
```

## Running tests

The tests require 3 variables set in order to run.
```
TF_ACC=true
FA_DOMAIN=https://YOUR.fusionauth.io
FA_API_KEY=YOUR_API_KEY
```

The API key needs to have key manager permissions (ability to create other API keys). If you add these to your computer/shell environment variables then executing the tests are as simple as:
```
go test ./...
```

Alternately you can supply on the command line when executing the tests.

```
TF_ACC=true FA_DOMAIN=https://YOUR.fusionauth.io FA_API_KEY=YOUR_API_KEY go test ./...
```

## Running lint

If you want to head off lint errors before hitting CI you can execute them locally.

First, install golangci-lint. Instructions can be found here: https://golangci-lint.run/usage/install/#local-installation

Then run it:
```
golangci-lint run
```

## To update the FusionAuth go-client

In this example, we are pulling the go-client at version `1.42.1`

```
go get -u github.com/FusionAuth/go-client@1.42.1
```

## Releasing

1. Ensure the tests pass (see above)
2. Tag the repo with v1.3.4 if 1.3.4 is the next, unreleased version. That will kick off a goreleaser workflow -  Ensure that succeeds.
3. After that, edit the release notes in GitHub Releases for the new release. Use the most recent PR (which is release specific) and past release notes as a guide.
