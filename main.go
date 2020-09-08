package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"gitlab.com/gpsi/api/fusionauth-tf-provider/fusionauth"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: fusionauth.Provider})
}
