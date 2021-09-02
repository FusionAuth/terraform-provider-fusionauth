package main

import (
	"github.com/gpsinsight/terraform-provider-fusionauth/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: fusionauth.Provider,
	})
}
