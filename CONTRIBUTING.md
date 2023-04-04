Contributing Guide

## Current Maintainers
- Caleb @cdavisgpsi
- Josiah @Jbcampbe
- Drew @drewlesueur
- Daniel @robotdan (FusionAuth)

## Guidline for Updating

1. Scroll through [FusionAuth's API Docs](https://fusionauth.io/docs/v1/tech/apis/)
2. Update each resouce and make sure to update the docs!!!
3. Submit PR to this repo
4. Upon merge, maintainer will create new git tag kicking off the build process.
5. [Terraform Registry](https://registry.terraform.io/providers/gpsinsight/fusionauth/latest) will pick up the changes

```
git tag v0.1.71
git push origin --tag
```

## To uppdate the FusionAuth go-client

In this example, we are pulling the go-client at version `1.42.1`

```
go get -u github.com/FusionAuth/go-client@1.42.1
```